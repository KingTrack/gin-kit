package consumer

import (
	"context"

	"github.com/Shopify/sarama"
)

type IConsumer interface {
	ConsumerGroup(ctx context.Context) sarama.ConsumerGroup
}
