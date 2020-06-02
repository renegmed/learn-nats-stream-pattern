package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	stan "github.com/nats-io/stan.go"
)

func main() {

	serverPort := os.Getenv("SERVER_PORT")
	natsURL := os.Getenv("NATS_ADDR")
	natsClusterID := os.Getenv("NATS_CLUSTER_ID")
	natsClientID := os.Getenv("NATS_CLIENT_ID")

	fmt.Printf("ClusterID: %s  ClientID: %s URL: %s\n", natsClusterID, natsClientID, natsURL)

	natsConn, err := stan.Connect(natsClusterID, natsClientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, natsURL)
	}
	defer natsConn.Close()

	srv := server{
		natsConn: natsConn,
	}

	// Serve HTTP
	r := mux.NewRouter()
	r.HandleFunc("/publish", srv.HandlePublishMessage)

	log.Printf("Starting HTTP server on '%s'", serverPort)

	if err := http.ListenAndServe(":"+serverPort, r); err != nil {
		log.Fatal(err)
	}
}
