package client

import (
	"context"
	"testing"

	"github.com/KingTrack/gin-kit/kit/types/redis/unknown"
	"github.com/stretchr/testify/assert"
)

func TestRedis_Client(t *testing.T) {
	assert.Equal(t, New("my.redis").Client(context.Background()).Get(context.Background(), "my.key").Err(), unknown.ErrUnknownRedisDB)
}
