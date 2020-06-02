package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	stan "github.com/nats-io/stan.go"
)

type Message struct {
	MessageID string `json:"message-id"`
	Topic     string `json:"topic"`
	Message   string `json:"message"`
}

type processed struct {
	messages map[string]Message
}

func newProcessed() processed {
	return processed{map[string]Message{}}
}

func (p *processed) ifexists(id string) bool {
	_, ok := p.messages[id]
	return ok
}

func (p *processed) save(m Message) {
	if !p.ifexists(m.MessageID) {
		p.messages[m.MessageID] = m
	} else {
		log.Printf("+++++ Message %s have been received before.\n", m.MessageID)
	}
}

func (p *processed) processMessage(m *stan.Msg) bool {
	msg, err := getMessage(m)
	if err != nil {
		log.Println(err)
		return false
	}

	if p.ifexists(msg.MessageID) {
		log.Println("++++ Message already received before", msg.MessageID)
		return false
	} else {
		p.save(msg)
	}

	if strings.Contains(msg.Message, "fail") {
		log.Printf("+++++ Message failed to process by sub1:\n\t %v\n", msg)
		return false
	}

	log.Printf("+++++ Message has been processed by sub1:\n\t %v\n", msg)
	return true
}

func getMessage(m *stan.Msg) (Message, error) {
	data := m.Data
	msg := Message{}
	if err := json.Unmarshal([]byte(data), &msg); err != nil {
		log.Fatalf("Error, could not unmarshal message, %v", err)
		return msg, err
	}
	return msg, nil
}
func main() {

	natsURL := os.Getenv("NATS_ADDR")
	natsClusterID := os.Getenv("NATS_CLUSTER_ID")
	natsClientID := os.Getenv("NATS_CLIENT_ID")
	qtopic := os.Getenv("QUEUE_TOPIC")
	qgroup := os.Getenv("QUEUE_GROUP")
	serverPort := os.Getenv("SERVER_PORT")

	natsConn, err := stan.Connect(natsClusterID, natsClientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, natsURL)
	}
	defer natsConn.Close()

	log.Printf("Connected to %s clusterID: [%s] clientID: [%s]\n", natsURL, natsClusterID, natsClientID)

	aw, _ := time.ParseDuration("1s")
	processedMsg := newProcessed()

	_, err = natsConn.Subscribe(qtopic, func(m *stan.Msg) {

		log.Printf("Received Message seq %d. Ready to process.", m.Sequence)
		success := processedMsg.processMessage(m)
		if success {
			log.Print("+++++ Sending acknowledgement.")
			m.Ack() // ack message after performing I/O intensive operation
		} else {
			if m.RedeliveryCount >= 2 {
				log.Printf("+++++ Failed to process after %d attempts. We give up. Ask NATS to stop redelivery.", m.RedeliveryCount+1)
				m.Ack()
			}
		}
	}, stan.SetManualAckMode(), stan.AckWait(aw))

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
