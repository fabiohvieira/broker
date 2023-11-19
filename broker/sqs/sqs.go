package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/fabiohvieira/broker/broker"
)

type Broker struct {
	client Client
}

func New(client Client) *Broker {
	return &Broker{client: client}
}

type Client interface {
	SendMessage(ctx context.Context, message Message, topic string) error
	ReceiveMessages(ctx context.Context, topic string) (*Message, error)
	Ack() error
}

func (b *Broker) SendMessage(ctx context.Context, message broker.Message, topic string) error {
	return b.client.SendMessage(ctx, message, topic)
}

func (b *Broker) ReceiveMessages(ctx context.Context, topic string) (*broker.Message, error) {
	return b.client.ReceiveMessages(ctx, topic)
}

func (b *Broker) Ack() error {
	return b.client.Ack()
}

func NewSQSClient() (*sqs.Client, error) {
	return &sqsClient{}
}
