package consumer

import (
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/bygui86/go-kafka/tracing/kafka-consumer/logging"
)

const (
	kafkaBootstrapServersConfigKey = "bootstrap.servers"
	kafkaConsumerGroupConfigKey    = "group.id"
	kafkaAutoOffsetResetConfigKey  = "auto.offset.reset"
)

func New(consumerName string) (*KafkaConsumer, error) {
	logging.Log.Info("Create new kafka consumer")

	cfg := loadConfig()

	kafkaConsumer, err := kafka.NewConsumer(
		// for other configurations, see https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
		&kafka.ConfigMap{
			kafkaBootstrapServersConfigKey: cfg.kafkaBootstrapServers,
			kafkaConsumerGroupConfigKey:    cfg.kafkaConsumerGroup,
			kafkaAutoOffsetResetConfigKey:  cfg.kafkaAutoOffsetReset,
		},
	)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		config:   cfg,
		name:     consumerName,
		consumer: kafkaConsumer,
		running:  false,
	}, nil
}

func (c *KafkaConsumer) Start() error {
	logging.Log.Info("Start kafka consumer")

	if c.consumer != nil {
		err := c.subscribeToTopics()
		if err != nil {
			return err
		}
		logging.SugaredLog.Infof("Kafka consumer subscribed to topics: %s",
			strings.Join(c.config.kafkaTopics, ","))

		go c.startConsumer()
		c.running = true
		logging.Log.Info("Kafka consumer started")
		return nil
	}

	return fmt.Errorf("kafka consumer start failed: consumer not initialized or already running")
}

func (c *KafkaConsumer) Shutdown(timeout int) {
	logging.SugaredLog.Warnf("Shutdown kafka consumer, timeout %d", timeout)

	if c.consumer != nil {
		err := c.consumer.Close()
		if err != nil {
			logging.SugaredLog.Errorf("Error closing kafka consumer: %s", err.Error())
		}
	}
}
