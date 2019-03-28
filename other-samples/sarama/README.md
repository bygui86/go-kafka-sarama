# Consumergroup example

This example shows you how to use the Sarama consumer group consumer. The example simply starts consuming the given Kafka topics and logs the consumed messages.

```bash
go run main.go -brokers="localhost:9092" -topics="typed-events" -group="die-consumers" -verbose
```

```bash
docker build . --tag sarama-kafka-consumer-sample:1.0
```
