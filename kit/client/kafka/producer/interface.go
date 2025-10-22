package client

import (
	"context"

	"github.com/Shopify/sarama"
)

type IProducer interface {
	SyncProducer(ctx context.Context) sarama.SyncProducer
}
