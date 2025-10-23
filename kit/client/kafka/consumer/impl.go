package consumer

import (
	"context"

	"github.com/KingTrack/gin-kit/kit/runtime"
	"github.com/KingTrack/gin-kit/kit/types/kafka/consumer/unknown"
	"github.com/Shopify/sarama"
)

type Consumer struct {
	name string
}

func New(name string) IConsumer {
	return &Consumer{name: name}
}

func (c *Consumer) ConsumerGroup(ctx context.Context) sarama.ConsumerGroup {
	if runtime.Get().KafkaConsumer() == nil {
		return unknown.New()
	}
	return runtime.Get().KafkaConsumer().GetConsumer(ctx, c.name)
}
