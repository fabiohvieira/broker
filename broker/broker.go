package broker

import "context"

type MessageBroker interface {
	SendMessage(ctx context.Context, message *Message, topic string) error
	ReceiveMessages(ctx context.Context, maxMessages int64, topic string) (*Message, error)
	Ack(ctx context.Context, message *Message, topic string) error
}

type NullBroker struct{}

func (n *NullBroker) SendMessage(_ context.Context, _ *Message, _ string) error {
	return nil
}

func (n *NullBroker) ReceiveMessages(_ context.Context, _ int64, _ string) (*Message, error) {
	return nil, nil
}

func (n *NullBroker) Ack(_ context.Context, _ *Message, _ string) error {
	return nil
}
