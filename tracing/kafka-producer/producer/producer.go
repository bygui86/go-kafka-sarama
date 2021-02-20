package producer

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/bygui86/go-kafka/tracing/kafka-producer/logging"
)

const (
	kafkaBootstrapServersConfigKey = "bootstrap.servers"
)

func New(producerName string) (*KafkaProducer, error) {
	logging.Log.Info("Create new kafka producer")

	cfg := loadConfig()

	kafkaProducer, err := kafka.NewProducer(
		// for other configurations, see https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
		&kafka.ConfigMap{
			kafkaBootstrapServersConfigKey: cfg.kafkaBootstrapServers,
		},
	)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		config:   cfg,
		name:     producerName,
		producer: kafkaProducer,
		stop:     make(chan bool, 1),
		running:  false,
		messages: []string{"Frodo Baggins", "Samvise Gamgee", "Meriadoc Brandibuck", "Peregrino Tuc", "Aragorn", "Boromir", "Legolas", "Gimli", "Gandalf"},
	}, nil
}

func (p *KafkaProducer) Start() {
	logging.Log.Info("Start kafka producer")

	if p.producer != nil {
		go p.startEventListener()
		go p.startProducer()
		p.running = true
		logging.Log.Info("Kafka producer started")
		return
	}

	logging.Log.Error("Kafka producer start failed: producer not initialized or already running")
}

func (p *KafkaProducer) Shutdown(timeout int) {
	logging.SugaredLog.Warnf("Shutdown kafka producer, timeout %d", timeout)

	time.Sleep(time.Duration(timeout) * time.Second)
	p.ticker.Stop()
	p.stop <- true
	if p.producer != nil {
		// Wait for message deliveries before shutting down
		p.producer.Flush(timeout * 1000)
		p.producer.Close()
	}
}
