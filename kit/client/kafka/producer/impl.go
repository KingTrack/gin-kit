package client

import (
	"context"

	"github.com/KingTrack/gin-kit/kit/runtime"
	"github.com/KingTrack/gin-kit/kit/types/kafka/producer/unknown"
	"github.com/Shopify/sarama"
)

type Producer struct {
	name string
}

func New(name string) IProducer {
	return &Producer{name: name}
}

func (p *Producer) SyncProducer(ctx context.Context) sarama.SyncProducer {
	if runtime.Get().KafkaProducer() == nil {
		return unknown.New()
	}
	return runtime.Get().KafkaProducer().GetProducer(ctx, p.name)
}
