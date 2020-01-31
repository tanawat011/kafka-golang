package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kafka-golang-example/producer/kafka_action"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	p "kafka-golang-example/producer/kafka_action/producer"
	mw "kafka-golang-example/producer/middleware"

	"github.com/gorilla/mux"
)

var topic = os.Getenv("KAFKA_TOPIC_NAME")
var bootstrapServer = os.Getenv("KAFKA_BOOTSTRAP_SERVER")
var numberOfPartition = os.Getenv("KAFKA_NUMBER_OF_PARTITION")

func main() {

	// Create initial topic
	numberOfPartition, _ := strconv.Atoi(numberOfPartition)
	kafka_action.CreateTopic(topic, bootstrapServer, numberOfPartition)

	// Initial kafka producer
	err := p.InitKafkaProducer(bootstrapServer)
	if err != nil {
		log.Fatal("Kafka producer ERROR: ", err)
	}

	// Initial router
	r := mux.NewRouter()

	// Initial Custom middleware
	r.Use(mw.CommonMiddleware)

	// Router
	r.HandleFunc("/", welcome).Methods("GET")
	r.HandleFunc("/producer", produceMessage).Methods("POST")

	//create http server
	port := os.Getenv("APP_PORT")
	fmt.Println("App listen localhost on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// Welcome is welcome
func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Welcome to golang api")

	response := map[string]string{
		"status":  "ok",
		"message": "Welcome to golang api",
	}
	json.NewEncoder(w).Encode(response)
}

func produceMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Started produce")

	// Parsing body request to json string
	bodyMsg, _ := ioutil.ReadAll(r.Body)
	// removeAllWhiteSpaceNotWithinQuote := regexp.MustCompile(`\s+(^([^"]*"[^"]*")*[^"]*$)`)
	// newMsg := removeAllWhiteSpaceNotWithinQuote.ReplaceAllString(string(newBodyMsg), "")
	// removeAllNewLineAndMore := regexp.MustCompile(`\r?\n?\t`)
	// newMsg = removeAllNewLineAndMore.ReplaceAllString(newBodyMsg, "")
	// fmt.Printf("Body message %s", newMsg)

	// Produce bodyMsg
	timeStart := time.Now()

	newBodyMsg, _ := setOrderCaptureName(bodyMsg)

	producerErr := p.Produce(topic, string(newBodyMsg)) // 1
	for i := 0; i < 1000000; i++ {

		newBodyMsgLoop, _ := setOrderCaptureName(newBodyMsg)

		producerErr = p.Produce(topic, string(newBodyMsgLoop)) // any
	}
	if producerErr != nil {
		log.Print(producerErr)
		os.Exit(1)
	}
	elapsed := time.Since(timeStart)
	fmt.Printf("End produce in %s", elapsed)

	json.NewEncoder(w).Encode("Produced success!")
}

// randomStringByte is random string byte and return string
func randomStringByte(n int) string {
	// rule := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	rule := "1234567890"
	b := make([]byte, n)
	for i := range b {
		b[i] = rule[rand.Intn(len(rule))]
	}
	return string(b)
}

// setOrderCaptureName is set order capture name
func setOrderCaptureName(bodyMsg []byte) ([]byte, error) {
	raw := make(map[string]interface{})
	json.Unmarshal(bodyMsg, &raw)
	orderReq := raw["submitOrderRequest"].(map[string]interface{})
	order := orderReq["Order"].(map[string]interface{})
	order["orderCaptureName"] = randomStringByte(20)
	fmt.Printf("Order capture name => %s", order["orderCaptureName"])
	orderReq["Order"] = order
	raw["submitOrderRequest"] = orderReq
	return json.Marshal(raw)
}
