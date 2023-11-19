package consumer

import (
	"context"
	"github.com/fabiohvieira/broker/broker"
	"sync"
)

type workerConfig struct {
	Topic       string
	MaxMessages int64
}

type Worker interface {
	Start(ctx context.Context)
}

func NewWorker(broker broker.MessageBroker, config workerConfig) Worker {
	return &worker{
		broker: broker,
		config: &config,
	}
}

type worker struct {
	config *workerConfig
	broker broker.MessageBroker
}

func (w *worker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			messages, err := w.broker.ReceiveMessages(ctx, w.config.MaxMessages, w.config.Topic)
			if err != nil {
				continue
			}

			if len(messages) > 0 {
				w.processMessages(ctx, messages)
			}
		}
	}
}

func (w *worker) processMessages(ctx context.Context, messages []*broker.Message) {
	var wg sync.WaitGroup

	numMessages := len(messages)

	wg.Add(numMessages)

	for i := range messages {
		go func(ctx context.Context, m *broker.Message) {
			defer wg.Done()
			w.handleMessages(ctx, m)
		}(ctx, messages[i])
	}

	wg.Wait()
}

func (w *worker) handleMessages(ctx context.Context, message *broker.Message) {
	if message == nil {
		return
	}

	err := w.processMessage(ctx, message)
	if err != nil {
		return
	}
}

func (w *worker) processMessage(ctx context.Context, message *broker.Message) error {
	return w.broker.Ack(ctx, message, w.config.Topic)
}
