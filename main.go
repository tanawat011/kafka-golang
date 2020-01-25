package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kafka-golang-multitreading/kafka_action"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	c "kafka-golang-multitreading/kafka_action/consumer"
	p "kafka-golang-multitreading/kafka_action/producer"
	mw "kafka-golang-multitreading/middleware"

	"github.com/gorilla/mux"
)

var topic = os.Getenv("KAFKA_TOPIC_NAME")
var groupID = os.Getenv("KAFKA_CONSUME_GROUP_ID")
var bootstrapServer = os.Getenv("KAFKA_BOOTSTRAP_SERVER")

func main() {

	// Create initial topic
	kafka_action.CreateTopic(topic, bootstrapServer)

	// Initialize kafka consumer
	go c.InitKafkaConsumer(bootstrapServer, groupID, topic)

	// Initialize kafka producer
	err := p.InitKafkaProducer(bootstrapServer)
	if err != nil {
		log.Fatal("Kafka producer ERROR: ", err)
	}

	// Initialize router
	r := mux.NewRouter()

	// Initialize Custom middleware
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

	// Produce bodyMsg
	timeStart := time.Now()
	producerErr := p.Produce(topic, string(bodyMsg)) // 1
	for i := 0; i < 30; i++ {
		newMsg := string(bodyMsg) + strconv.Itoa(i)
		producerErr = p.Produce(topic, newMsg) // any
	}
	if producerErr != nil {
		log.Print(producerErr)
		os.Exit(1)
	}
	elapsed := time.Since(timeStart)
	fmt.Printf("End in %s", elapsed)

	json.NewEncoder(w).Encode("Produced success!")
}
