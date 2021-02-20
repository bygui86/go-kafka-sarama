package consumer

import (
	"github.com/bygui86/go-kafka/tracing/kafka-consumer/logging"
	"github.com/bygui86/go-kafka/tracing/kafka-consumer/utils"
)

const (
	kafkaBootstrapServersEnvVar = "KAFKA_BOOTSTRAP_SERVERS" // comma-separated list host or host:port
	kafkaPortEnvVar             = "KAFKA_PORT"
	kafkaTopicsEnvVar           = "KAFKA_TOPICS" // comma-separated list
	kafkaConsumerGroupEnvVar    = "KAFKA_CONSUMER_GROUP"
	kafkaAutoOffsetResetEnvVar  = "KAFKA_AUTO_OFFSET_RESET" // earliest, latest, etc. (see https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md)

	kafkaBootstrapServersDefault = "0.0.0.0"
	kafkaPortDefault             = 9092
	kafkaConsumerGroupDefault    = "my-group"
	kafkaAutoOffsetResetDefault  = "latest"
)

var (
	kafkaTopicsDefault = []string{"my-topic"}
)

func loadConfig() *Config {
	logging.Log.Debug("Load kafka producer configurations")
	return &Config{
		kafkaBootstrapServers: utils.GetStringEnv(kafkaBootstrapServersEnvVar, kafkaBootstrapServersDefault),
		kafkaPort:             utils.GetIntEnv(kafkaPortEnvVar, kafkaPortDefault),
		kafkaTopics:           utils.GetStringListEnv(kafkaTopicsEnvVar, kafkaTopicsDefault),
		kafkaConsumerGroup:    utils.GetStringEnv(kafkaConsumerGroupEnvVar, kafkaConsumerGroupDefault),
		kafkaAutoOffsetReset:  utils.GetStringEnv(kafkaAutoOffsetResetEnvVar, kafkaAutoOffsetResetDefault),
	}
}
