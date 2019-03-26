
# Kafka producer/concumer by erengaeb

## Instructions

1. Setup environment variables
	```
	export DOCKER_HOST_NAME="dockerhost"
	export DOCKER_HOST_IP="$(ifconfig en0 | grep inet | grep -v inet6 | awk '{print $2}')"
	```

2. Run Zookeeper on Docker
	```
	docker run -d --name zookeeper -p 2181:2181 --network bridge --add-host=$DOCKER_HOST_NAME:$DOCKER_HOST_IP jplock/zookeeper
	```

3. Run Kafka on Docker
	```
	docker run -d --name kafka -e KAFKA_ADVERTISED_HOST_NAME=$DOCKER_HOST_NAME -e ZOOKEEPER_IP=$DOCKER_HOST_NAME --network bridge --add-host=$DOCKER_HOST_NAME:$DOCKER_HOST_IP -p 7203:7203 -p 9092:9092 ches/kafka
	```

4. Create a topic
	```
	docker run --rm ches/kafka kafka-topics.sh --create --topic senz --replication-factor 1 --partitions 1 --zookeeper $DOCKER_HOST_NAME:2181
	```

5. List all topics
	```
	docker run --rm ches/kafka kafka-topics.sh --list --zookeeper $DOCKER_HOST_NAME:2181
	```

6. Run sample producer - `NOT WORKING`
	```
	docker run --rm --interactive ches/kafka kafka-console-producer.sh --topic senz --broker-list $DOCKER_HOST_NAME:2181
	```

7. Run sample consumer - `NOT WORKING`
	```
	docker run --rm ches/kafka kafka-console-consumer.sh --topic senz --from-beginning --zookeeper $DOCKER_HOST_NAME:2181
	docker run --rm ches/kafka kafka-console-consumer.sh --topic senz --from-beginning --bootstrap-server $DOCKER_HOST_NAME:2181
	```

---

## Links

* [Kafka/Zookeeper on Docker](https://medium.com/rahasak/kafka-and-zookeeper-with-docker-65cff2c2c34f)
* [Kafka producer](https://medium.com/@itseranga/kafka-producer-with-golang-fab7348a5f9a)
* [Kafka consumer](https://medium.com/@itseranga/kafka-consumer-with-golang-a93db6131ac2)
