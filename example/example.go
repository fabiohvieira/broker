package main

import (
	"context"
	"github.com/fabiohvieira/broker/broker"
	"github.com/fabiohvieira/broker/broker/sqs"
)

func main() {
	sqsClient, err := sqs.NewSQSClient(context.Background())
	if err != nil {
		panic(err)
	}

	b := sqs.New(sqsClient)

	message := broker.Message{
		Body: []byte("test"),
	}

	err = b.SendMessage(context.Background(), message, "topic")
	if err != nil {
		panic(err)
	}
}
