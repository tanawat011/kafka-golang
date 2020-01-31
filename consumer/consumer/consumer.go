package consumer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var producer *kafka.Producer

// InitKafkaConsumer is Initialize kafka consumer
func InitKafkaConsumer(bootstrapServer string, groupID string, topic string) {
	fmt.Println("Started consume")

	// Initialize message consumer
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  bootstrapServer,
		"group.id":           groupID,
		"enable.auto.commit": false,
		"auto.offset.reset":  "earliest",
		"log.queue":          true,
	})

	if err != nil {
		panic(err)
	}

	topics := []string{topic}
	err = c.SubscribeTopics(topics, nil)
	fmt.Printf("Subscribed Topics : %s", topics)

	timeStart1 := time.Now()
	for {
		timeStart := time.Now()
		msg, err := c.ReadMessage(-1)
		c.Commit()
		if err != nil {
			// Log error and let container restart
			elapsed := time.Since(timeStart)
			fmt.Printf("Read error %s", err)
			_ = fmt.Sprintf("Read timeout happen in %f hour(s)", elapsed.Hours())
			// os.Exit(1)
		}
		// consumed := fmt.Sprintf("Consumed from => topic [%s] [%d] at offset [%v]:",
		// 	*msg.TopicPartition.Topic,
		// 	msg.TopicPartition.Partition,
		// 	msg.TopicPartition.Offset,
		// )
		// debuggerDebug(consumed, string(msg.Value))

		getOrderCaptureNameAndRequestCondition(msg.Value)

		elapsed := time.Since(timeStart)
		elapsedConsumed := fmt.Sprintf("Consumed end in %s", elapsed)
		debuggerDebug(elapsedConsumed, fmt.Sprintf("at offset [%v] in partition [%d] and in time => (%s)", msg.TopicPartition.Offset, msg.TopicPartition.Partition, timeStart))
		fmt.Printf("\n***************************************************\n")
		fmt.Printf("All Consumed end in %s \n", time.Since(timeStart1))
		fmt.Printf("***************************************************\n\n")
	}

	// fmt.Println("Consumed!")
}

// getOrderCaptureNameAndRequestCondition is get order capture name and request condition
func getOrderCaptureNameAndRequestCondition(bodyMsg []byte) {
	raw := make(map[string]interface{})
	json.Unmarshal(bodyMsg, &raw)
	orderReq := raw["submitOrderRequest"].(map[string]interface{})
	order := orderReq["Order"].(map[string]interface{})
	str := (order["orderCaptureName"]).(string)
	lastStr := lastString(strings.Split(str, ""))
	fmt.Printf("Last string => %s \n", lastStr)
	lastStrInt, _ := strconv.Atoi(lastStr)
	switch lastStrInt {
	case 5:
		go func() {
			_, err := http.Get("http://5e3263b5b92d240014ea51e6.mockapi.io/api/true/example")
			if err != nil {
				fmt.Printf("The HTTP request failed with error %s \n", err)
			}
			// else {
			// 	data, _ := ioutil.ReadAll(response.Body)
			// 	fmt.Printf(" Response => %s \n", data)
			// }
		}()
	}
}

func lastString(ss []string) string {
	return ss[len(ss)-1]
}

func debuggerInfo(str string) {
	fmt.Printf("\033[1;34m%s \033[0m \n", str)
}

func debuggerNotice(str string) {
	fmt.Printf("\033[1;36m%s \033[0m \n", str)
}

func debuggerWarning(str string) {
	fmt.Printf("\033[1;33m%s \033[0m \n", str)
}

func debuggerError(str string) {
	fmt.Printf("\033[1;31m%s \033[0m \n", str)
}

func debuggerDebug(str string, str2 string) {
	fmt.Printf("\033[0;36m%s \033[0m %s \n", str, str2)
}
