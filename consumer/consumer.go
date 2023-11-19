package consumer

import (
	"context"
	"errors"
	"github.com/fabiohvieira/broker/broker"
	"sync"
)

type Config struct {
	Broker       broker.MessageBroker
	WorkersCount uint8
	Topic        string
	MaxMessages  int64
}

type Consumer interface {
	Start(ctx context.Context)
}

type WorkerGetter = func(config Config) Worker

func New(config Config) (Consumer, error) {
	if config.WorkersCount < 1 {
		return nil, errors.New("workers count must be greater than 0")
	}
	if config.MaxMessages < 1 {
		return nil, errors.New("max messages must be greater than 0")
	}

	return &DefaultConsumer{config: config, getWorker: func(config Config) Worker {
		return NewWorker(config.Broker, workerConfig{
			Topic:       config.Topic,
			MaxMessages: config.MaxMessages,
		})
	}}, nil
}

type DefaultConsumer struct {
	config    Config
	getWorker WorkerGetter
}

func (c *DefaultConsumer) Start(ctx context.Context) {
	var wg sync.WaitGroup
	for i := uint8(0); i < c.config.WorkersCount; i++ {
		wg.Add(1)
		go c.run(ctx, &wg)
	}
}

func (c *DefaultConsumer) run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	w := c.getWorker(c.config)
	w.Start(ctx)
}
