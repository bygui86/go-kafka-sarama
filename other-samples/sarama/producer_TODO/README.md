
# Shopify - Sarama - Producer

This example shows you how to use the Sarama producer. The example simply starts pushing to the given Kafka topics.

```bash
go run main.go -brokers="localhost:9092" -topics="message" -group="sample-consumers" -verbose
```

```bash
docker build . --tag sarama-kafka-producer-sample:1.0
```

See https://github.com/Shopify/sarama/tree/master/examples/http_server
