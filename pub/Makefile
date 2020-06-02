#export NATS_URL="nats://172.22.0.3:4222" #,nats://172.22.0.2:4222,nats://172.22.0.4:4222"
export NATS_ADDR=nats://nats1:4222 #,nats://172.22.0.2:4222,nats://172.22.0.4:4222"
export NATS_CLUSTER_ID=cluster-one
export NATS_CLIENT_ID=pub-1
export SERVER_PORT=9090

init-project:
	go mod init github.com/renegmed/learn-nats-stream-pattern/pub
up:


run:
	go run . 
.PHONY: run

post:
	curl -X POST http://localhost:9090/publish -d {"topic": "test.topic", "message":"this is a test."}
.PHONY: post
