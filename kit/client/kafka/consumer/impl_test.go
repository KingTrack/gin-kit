package consumer

import (
	"context"
	"testing"

	"github.com/KingTrack/gin-kit/kit/types/kafka/consumer/unknown"
	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

func TestConsumer_ConsumerGroup(t *testing.T) {
	assert.Error(t, New("my.consumer").ConsumerGroup(context.Background()).
		Consume(context.Background(), []string{"my.topic"}, &mockConsumerHandler{}), unknown.ErrUnknownConsumer)
}

type mockConsumerHandler struct{}

func (h *mockConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *mockConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *mockConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		session.MarkMessage(msg, "")
		session.Commit()
	}

	return nil
}
