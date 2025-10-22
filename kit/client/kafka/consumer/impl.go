package consumer

import (
	"context"

	"github.com/KingTrack/gin-kit/kit/runtime"

	"github.com/Shopify/sarama"
)

type Consumer struct {
	name string
}

func New(name string) IConsumer {
	return &Consumer{name: name}
}

func (c *Consumer) ConsumerGroup(ctx context.Context) sarama.ConsumerGroup {
	return runtime.Get().KafkaConsumer().GetConsumer(ctx, c.name)
}
