package consumer

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/bygui86/go-kafka/tracing/kafka-consumer/commons"
	"github.com/bygui86/go-kafka/tracing/kafka-consumer/logging"
	"github.com/bygui86/go-kafka/tracing/kafka-consumer/tracing"
)

func (c *KafkaConsumer) subscribeToTopics() error {
	subErr := c.consumer.SubscribeTopics(c.config.kafkaTopics, nil)
	if subErr != nil {
		return subErr
	}
	return nil
}

func (c *KafkaConsumer) startConsumer() {
	for {
		msg, err := c.consumer.ReadMessage(-1)
		if err == nil {
			carrier := tracing.KafkaHeadersCarrier(msg.Headers)
			spanCtx, extErr := tracing.Extract(&carrier)
			if extErr != nil {
				logging.SugaredLog.Errorf("Error extracting span context from message: %s", extErr.Error())
			}

			var span opentracing.Span
			if spanCtx != nil {
				span = opentracing.StartSpan(c.name, ext.RPCServerOption(spanCtx))
			}

			topicInfo, msgInfo, headersInfo := c.getMessageInfo(msg)
			logging.SugaredLog.Infof("Message received: topicInfo[%s], msgInfo[%s], headersInfo[%s]",
				topicInfo, msgInfo, headersInfo)

			if span != nil {
				span.SetTag("app", commons.ServiceName)
				span.Finish()
			}

		} else {
			// INFO: The client will automatically try to recover from all errors.
			if msg != nil {
				topicInfo, msgInfo, headersInfo := c.getMessageInfo(msg)
				logging.SugaredLog.Errorf("Consumer error on message: topicInfo[%s], msgInfo[%s], headersInfo[%s], error[%s]",
					topicInfo, msgInfo, headersInfo, err.Error())
			} else {
				logging.SugaredLog.Errorf("Consumer error: %s", err.Error())
			}
		}
	}
}

func (c *KafkaConsumer) getMessageInfo(msg *kafka.Message) (string, string, string) {
	topicInfo := fmt.Sprintf("name[%s], partition[%d], offset[%d]",
		*msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
	msgInfo := fmt.Sprintf("key[%s], value[%s], timestamp[%v]",
		string(msg.Key), string(msg.Value), msg.Timestamp)
	headersInfo := fmt.Sprintf("%+v", msg.Headers)
	return topicInfo, msgInfo, headersInfo
}
