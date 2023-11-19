package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/fabiohvieira/broker/broker"
)

type Broker struct {
	client Client
}

func New(client Client) *Broker {
	return &Broker{
		client: client,
	}
}

type Client interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	Ack(ctx context.Context) error
}

func (b *Broker) SendMessage(ctx context.Context, message broker.Message, topic string) error {
	params := &sqs.SendMessageInput{
		MessageBody: aws.String(string(message.Body)),
		QueueUrl:    aws.String(topic),
	}

	_, err := b.client.SendMessage(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func (b *Broker) ReceiveMessages(ctx context.Context, topic string) ([]broker.Message, error) {
	params := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(topic),
	}

	output, err := b.client.ReceiveMessage(ctx, params)
	if err != nil {
		return nil, err
	}

	var messages []broker.Message
	for _, message := range output.Messages {
		messages = append(messages, broker.Message{
			Header: map[string]string{},
			Body:   []byte(*message.Body),
		})
	}

	return messages, nil
}

func (b *Broker) Ack(ctx context.Context) error {
	return b.client.Ack(ctx)
}

func NewSQSClient(ctx context.Context) (*sqs.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return sqs.NewFromConfig(cfg), nil
}
