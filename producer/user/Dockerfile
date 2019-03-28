
FROM golang:1.9

RUN go get github.com/Shopify/sarama

ENV KAFKA_HOST localhost
ENV KAFKA_PORT 9092
ENV KAFKA_TOPIC gomessages

ADD . /app/
WORKDIR /app

RUN go build -o build/producer producer.go

ENTRYPOINT ["/app/build/producer"]
