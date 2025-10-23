package unknown

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

var (
	ErrUnknownConsumer = errors.New("unknown kafka consumer")
)

type Consumer struct{}

func New() *Consumer {
	return &Consumer{}
}

func (c *Consumer) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	return ErrUnknownConsumer
}

func (c *Consumer) Errors() <-chan error {
	ch := make(chan error, 1)
	ch <- ErrUnknownConsumer
	close(ch)
	return ch
}

// Close stops the ConsumerGroup and detaches any running sessions. It is required to call
// this function before the object passes out of scope, as it will otherwise leak memory.
func (c *Consumer) Close() error {
	return ErrUnknownConsumer
}

func (c *Consumer) Pause(partitions map[string][]int32) {}

func (c *Consumer) Resume(partitions map[string][]int32) {}

func (c *Consumer) PauseAll() {}

func (c *Consumer) ResumeAll() {}
