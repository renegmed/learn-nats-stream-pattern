package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	stan "github.com/nats-io/stan.go"
)

type Message struct {
	MessageID string `json:"message-id"`
	Topic     string `json:"topic"`
	Message   string `json:"message"`
}

func processMessage(m *stan.Msg) {
	data := m.Data
	msg := Message{}
	if err := json.Unmarshal([]byte(data), &msg); err != nil {
		log.Fatalf("Error, could not unmarshal message, %v", err)
		return
	}

	log.Printf("Message has been processed by sub2:\n\t %s\n", msg)
}

func main() {

	natsURL := os.Getenv("NATS_ADDR")
	natsClusterID := os.Getenv("NATS_CLUSTER_ID")
	natsClientID := os.Getenv("NATS_CLIENT_ID")
	qtopic := os.Getenv("QUEUE_TOPIC")
	qgroup := os.Getenv("QUEUE_GROUP")
	serverPort := os.Getenv("SERVER_PORT")

	// Connect to NATS
	natsConn, err := stan.Connect(natsClusterID, natsClientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, natsURL)
	}
	defer natsConn.Close()

	log.Printf("Connected to %s clusterID: [%s] clientID: [%s]\n", natsURL, natsClusterID, natsClientID)

	mcb := func(msg *stan.Msg) {
		processMessage(msg)
	}

	_, err = natsConn.Subscribe(qtopic, mcb)
	if err != nil {
		natsConn.Close()
		log.Fatalf("Error on queue subscribe; %v", err)
	}

	log.Printf("Listening on [%s], clientID=[%s], qgroup=[%s]\n", qtopic, natsClientID, qgroup)

	// Serve HTTP
	r := mux.NewRouter()
	log.Printf("Starting HTTP server on '%s'", serverPort)

	if err := http.ListenAndServe(":"+serverPort, r); err != nil {
		log.Fatal(err)
	}
}
