
FROM golang:1.9

RUN go get github.com/Shopify/sarama
RUN go get github.com/wvanbergen/kafka/consumergroup

ENV ZOOKEEPER_HOST localhost
ENV ZOOKEEPER_PORT 2181
ENV TOPIC gomessages

ADD consumer.go /app/consumer.go
WORKDIR /app

RUN go build -o build/consumer consumer.go

ENTRYPOINT ["/app/build/consumer"]
