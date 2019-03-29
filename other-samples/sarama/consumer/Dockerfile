
FROM golang:1.12.1 as builder

RUN rm -rf /go/src/sarama-kafka-consumer-sample
COPY . /go/src/sarama-kafka-consumer-sample
WORKDIR /go/src/sarama-kafka-consumer-sample

RUN ["/bin/bash", "-c", "go get -v -d ./..."]
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/app .



FROM alpine

RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates

WORKDIR /bin/

COPY --from=builder /bin/app .

ENTRYPOINT [ "/bin/app", "-brokers='kafka:9092'", "-topics='typed-events'", "-group='die-consumers'", "-verbose" ]
