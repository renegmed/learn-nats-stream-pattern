FROM golang:1.14.0-alpine

WORKDIR src/publisher

COPY . .
 
RUN pwd
RUN ls -als  

RUN go mod vendor 

RUN go build -o publisher .

CMD ["./publisher"]
