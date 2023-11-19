package broker

import "context"

type MessageBroker interface {
	SendMessage(ctx context.Context, message Message, topic string) error
	ReceiveMessages(ctx context.Context, topic string) (*Message, error)
	Ack(ctx context.Context) error
}

type NullBroker struct{}

func (n *NullBroker) SendMessage(ctx context.Context, message Message, topic string) error {
	return nil
}

func (n *NullBroker) ReceiveMessages(ctx context.Context, topic string) (*Message, error) {
	return nil, nil
}

func (n *NullBroker) Ack() error {
	return nil
}
