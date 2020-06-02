FROM golang:1.14.0-alpine

WORKDIR src/subscriber

COPY . .
 
# RUN pwd
# RUN ls -als  

RUN go mod vendor 

RUN go build -o subscriber .

CMD ["./subscriber"]
