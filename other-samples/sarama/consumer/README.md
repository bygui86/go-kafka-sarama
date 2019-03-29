
# Shopify - Sarama - Consumer

This example shows you how to use the Sarama consumer group consumer. The example simply starts consuming the given Kafka topics and logs the consumed messages.

```bash
go run main.go -brokers="localhost:9092" -topics="message" -group="sample-consumers" -verbose
```

```bash
docker build . --tag sarama-kafka-consumer-sample:1.0
```

See https://github.com/Shopify/sarama/tree/master/examples/consumergroup
