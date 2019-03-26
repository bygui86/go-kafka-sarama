FROM golang:1.9

MAINTAINER Eranga Bandara (erangaeb@gmail.com)

# install dependencies
RUN go get github.com/Shopify/sarama

# env
ENV KAFKA_HOST dev.localhost
ENV KAFKA_PORT 9092
ENV KAFKA_TOPIC senz

# copy app
ADD . /app
WORKDIR /app

# build
RUN go build -o build/produzer src/*.go

ENTRYPOINT ["/app/docker-entrypoint.sh"]
