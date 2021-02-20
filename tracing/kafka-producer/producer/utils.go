package producer

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/opentracing/opentracing-go"

	"github.com/bygui86/go-kafka/tracing/kafka-producer/commons"
	"github.com/bygui86/go-kafka/tracing/kafka-producer/logging"
	"github.com/bygui86/go-kafka/tracing/kafka-producer/tracing"
)

// Delivery report handler for produced messages
func (p *KafkaProducer) startEventListener() {
	logging.Log.Info("Start kafka event listener")

	// WARN: headers are not yet supported by kafka.Producer.Events channel
	for e := range p.producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				topicInfo := fmt.Sprintf("name[%s], partition[%d], offset[%d]",
					*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				msgInfo := fmt.Sprintf("key[%s], value[%s], timestamp[%v]",
					string(ev.Key), string(ev.Value), ev.Timestamp)
				headersInfo := fmt.Sprintf("%+v", ev.Headers)
				logging.SugaredLog.Errorf("Message delivery FAILED: topicInfo[%s], msgInfo[%s], headersInfo[%s]",
					topicInfo, msgInfo, headersInfo)
			} else {
				topicInfo := fmt.Sprintf("name[%s], partition[%d], offset[%d]",
					*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				msgInfo := fmt.Sprintf("key[%s], value[%s], timestamp[%v]",
					string(ev.Key), string(ev.Value), ev.Timestamp)
				headersInfo := fmt.Sprintf("%+v", ev.Headers)
				logging.SugaredLog.Infof("Message delivery: topicInfo[%s], msgInfo[%s], headersInfo[%s]",
					topicInfo, msgInfo, headersInfo)
			}
		}
	}
}

// Produce messages to topic (asynchronously)
func (p *KafkaProducer) startProducer() {
	p.ticker = time.NewTicker(1 * time.Second)
	topicPartition := kafka.TopicPartition{Topic: &p.config.kafkaTopic, Partition: kafka.PartitionAny}
	counter := 0
	for {
		select {
		case <-p.stop:
			return
		case <-p.ticker.C:
			span := opentracing.StartSpan(p.name)

			span.SetTag("app", commons.ServiceName)

			msg := p.messages[counter]

			carrier := tracing.KafkaHeadersCarrier([]kafka.Header{
				{"example", []byte("example-value")},
				{"app", []byte(commons.ServiceName)},
			})
			traceErr := tracing.Inject(span, &carrier)
			if traceErr != nil {
				logging.SugaredLog.Errorf("Producer failed to inject tracing span: %s", traceErr.Error())
				// continue
			}

			kafkaMsg := &kafka.Message{
				TopicPartition: topicPartition,
				Value:          []byte(msg),
				Headers:        carrier,
			}

			err := p.producer.Produce(kafkaMsg, nil)
			if err != nil {
				logging.SugaredLog.Errorf("Producer failed to publish message %s: %s", msg, err.Error())
				continue
			}

			if counter == len(p.messages)-1 {
				counter = 0
			} else {
				counter++
			}

			span.Finish()
		}
	}
}
