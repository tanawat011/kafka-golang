package consumer

import (
	"fmt"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var producer *kafka.Producer

// InitKafkaConsumer is Initialize kafka consumer
func InitKafkaConsumer(bootstrapServer string, groupID string, topic string) {
	fmt.Println("Started consume")

	// Initialize message consumer
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	topics := []string{topic}
	c.SubscribeTopics(topics, nil)
	fmt.Printf("Subscribed Topics : [%s]", topics)

	for {
		startRead := time.Now()
		msg, err := c.ReadMessage(-1)
		if err != nil {
			// Log error and let container restart
			duration := time.Since(startRead)
			fmt.Printf("Read error %s", err)
			_ = fmt.Sprintf("Read timeout happen in %f hour(s)", duration.Hours())
			// os.Exit(1)
		}
		fmt.Printf("consumed from topic %s [%d] at offset %v: "+
			string(msg.Value), *msg.TopicPartition.Topic,
			msg.TopicPartition.Partition, msg.TopicPartition.Offset)
	}

	// c.Close()
}
