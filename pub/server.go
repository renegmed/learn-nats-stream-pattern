package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"log"

	stan "github.com/nats-io/stan.go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type Message struct {
	MessageID string `json:"message-id"`
	Topic     string `json:"topic"`
	Message   string `json:"message"`
}

type server struct {
	natsConn stan.Conn
}

func (s *server) HandlePublishMessage(rw http.ResponseWriter, req *http.Request) {

	log.Println("Request method:", req.Method)
	switch req.Method {
	case "POST":
		s.publishMessage(rw, req)
	default:
		log.Printf("Invalid reques method: %s", req.Method)
		http.Error(rw, "Invalid request", http.StatusBadRequest)
	}
}

func (s *server) publishMessage(rw http.ResponseWriter, req *http.Request) {

	var msg Message

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, "Bad Request", http.StatusBadRequest)
		return
	}

	// log.Println("Request Body:", string(body))

	if err := json.Unmarshal(body, &msg); err != nil {
		log.Printf("Failed to read request: %v", err)
		http.Error(rw, "Invalid request", http.StatusBadRequest)
		return
	}

	messageID := uuid.NewV4().String()
	msg.MessageID = messageID

	if err := s.publish(msg); err != nil {
		log.Printf("Failed to publish message onto queue '%s': %v", msg.Topic, err)
		http.Error(rw, "", http.StatusInternalServerError)
		return
	}
}

func (s *server) publish(msg Message) error {

	bs, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "failed to marshal message")
	}

	nuid, err := s.natsConn.PublishAsync(msg.Topic, bs, ackHandler)
	if err != nil {
		return errors.Wrap(err, "failed to publish message")
	}

	log.Printf("Message has been published on topic: %v, with acknowledge %v\n", msg.Topic, nuid)
	return nil
}

func ackHandler(ackedNuid string, err error) {
	if err != nil {
		log.Printf("Warning: error publishing msg id %s: %v\n", ackedNuid, err.Error())
	} else {
		log.Printf("Received ack for msg id %s\n", ackedNuid)
	}
}
