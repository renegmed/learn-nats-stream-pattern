up:
	docker-compose up --build -d 
.PHONY: up

cluster-info:
	curl http://127.0.0.1:8222/varz
	curl http://127.0.0.1:8222/routez
.PHONY: cluster-info


post:
	curl -XPOST http://localhost:9000/publish -d '{"topic":"topic.test","message":"this is a test."}'
	curl -XPOST http://localhost:9000/publish -d '{"topic":"topic.test","message":"this is another test."}'
.PHONY: post 

tail-pub:
	docker logs pub -f 
tail-sub1:
	docker logs sub_1 -f 
tail-sub2:
	docker logs sub_2 -f 
tail-nats1:
	docker logs nats_1 -f 
tail-nats2:
	docker logs nats_2 -f
tail-nats3:
	docker logs nats_3 -f
.PHONY: tail-pub tail-sub1 tail-sub2 tail-nats1 tail-nats2 tail-nats3