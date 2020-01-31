package main

import (
	"os"

	"kafka-golang-example/consumer/consumer"
)

var topic = os.Getenv("KAFKA_TOPIC_NAME")
var groupID = os.Getenv("KAFKA_CONSUME_GROUP_ID")
var bootstrapServer = os.Getenv("KAFKA_BOOTSTRAP_SERVER")

func main() {
	// Initial kafka consumer
	consumer.InitKafkaConsumer(bootstrapServer, groupID, topic)
}
