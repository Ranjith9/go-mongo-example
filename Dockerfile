FROM golang:alpine

RUN apk add --no-cache bash git && git clone https://github.com/Ranjith9/go-mongo-example.git /go/src/go-mongo-example/

WORKDIR /go/src/go-mongo-example

RUN go get gopkg.in/mgo.v2

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sample

ENTRYPOINT ./sample

EXPOSE 3000
