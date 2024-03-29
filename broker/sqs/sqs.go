package sqs

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
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
	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}

func (b *Broker) SendMessage(ctx context.Context, message broker.Message, topic string) error {
	params := &sqs.SendMessageInput{
		MessageBody:       aws.String(string(message.Body)),
		QueueUrl:          aws.String(topic),
		MessageAttributes: map[string]types.MessageAttributeValue{},
	}

	for key, value := range message.Header {
		params.MessageAttributes[key] = types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(value),
		}
	}

	_, err := b.client.SendMessage(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func (b *Broker) ReceiveMessages(ctx context.Context, maxMessages int64, topic string) ([]*broker.Message, error) {
	params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(topic),
		MaxNumberOfMessages: int32(maxMessages),
	}

	output, err := b.client.ReceiveMessage(ctx, params)
	if err != nil {
		return nil, err
	}

	var messages []*broker.Message
	for _, message := range output.Messages {
		messages = append(messages, &broker.Message{
			Header: map[string]string{},
			Body:   []byte(*message.Body),
		})
	}

	return messages, nil
}

func (b *Broker) Ack(ctx context.Context, message *broker.Message, topic string) error {
	h, ok := message.Header["ReceiptHandle"]
	if !ok {
		return errors.New("ReceiptHandle not found")
	}

	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(topic),
		ReceiptHandle: &h,
	}

	_, err := b.client.DeleteMessage(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func NewSQSClient(ctx context.Context) (*sqs.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return sqs.NewFromConfig(cfg), nil
}
