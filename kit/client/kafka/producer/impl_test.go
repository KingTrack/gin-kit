package client

import (
	"context"
	"testing"

	"github.com/KingTrack/gin-kit/kit/types/kafka/producer/unknown"
	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

func TestProducer_SyncProducer(t *testing.T) {
	_, _, err := New("my.producer").SyncProducer(context.Background()).SendMessage(&sarama.ProducerMessage{})
	assert.Error(t, err, unknown.ErrUnknownProducer)
}
