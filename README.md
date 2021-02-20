
# Go Kafka

## Instructions

1. Prepare environment

   ```bash
   mkdir -p $GOPATH/src/github.com/byui86
   
   cd $GOPATH/src/github.com/byui86
   git clone git@github.com:bygui86/go-kafka.git
   
   cd go-kafka
   go get ./...
   ```

2. Setup environment variables
	`TODO: replace with docker network`

   ```bash
   export DOCKER_HOST_NAME="dockerhost"
   export DOCKER_HOST_IP="$(ifconfig en0 | grep inet | grep -v inet6 | awk '{print $2}')"
   ```

3. Run Zookeeper on Docker

   ```bash
   docker run -d --name zookeeper --restart always --network bridge --add-host=$DOCKER_HOST_NAME:$DOCKER_HOST_IP -p 2181:2181 -p 2888:2888 -p 3888:3888 zookeeper
   ```

3. Run Kafka on Docker

   ```bash
   docker run -d --name kafka --restart always -e KAFKA_ZOOKEEPER_CONNECT=dockerhost:2181 -e KAFKA_ADVERTISED_HOST_NAME=$DOCKER_HOST_IP \
       --network bridge --add-host=$DOCKER_HOST_NAME:$DOCKER_HOST_IP -p 9092:9092 -p 7203:7203 wurstmeister/kafka
   ```

4. Create Kafka topic

   ```bash
   docker exec -ti kafka opt/kafka_2.12-2.1.1/bin/kafka-topics.sh --create --topic gomessages --replication-factor 1 --partitions 1 --zookeeper $DOCKER_HOST_NAME:2181
   ```

5. List all Kafka topics

   ```bash
   docker exec -ti kafka opt/kafka_2.12-2.1.1/bin/kafka-topics.sh --list --zookeeper $DOCKER_HOST_NAME:2181
   ```

6. Run producer on Docker
	
	```bash
	cd $GOPATH/src/github.com/bygui86/go-kafka/producer
	docker build . --tag go-kafka-producer:1.0
	docker run -ti --rm --name producer go-kafka-producer:1.0
	```

7. Run consumer on Docker

	```bash
	cd $GOPATH/src/github.com/bygui86/go-kafka/consumer
	docker build . --tag go-kafka-consumer:1.0
	docker run -ti --rm --name consumer go-kafka-consumer:1.0
	```

---

## Running applications without Docker

1. Run producer
	```
	cd $GOPATH/src/github.com/bygui86/go-kafka/producer
	go get ./...
	go run producer.go
	```

2. Run consumer
	```
	cd $GOPATH/src/github.com/bygui86/go-kafka/consumer
	go get ./...
	go run consumer.go
	```

---

## Links

- https://medium.com/@yusufs/getting-started-with-kafka-in-golang-14ccab5fa26
- https://medium.com/rahasak/kafka-and-zookeeper-with-docker-65cff2c2c34f
- https://medium.com/@itseranga/kafka-producer-with-golang-fab7348a5f9a
- https://medium.com/@itseranga/kafka-consumer-with-golang-a93db6131ac2
