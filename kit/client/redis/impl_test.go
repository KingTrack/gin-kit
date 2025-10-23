package client

import (
	"context"
	"testing"

	"github.com/KingTrack/gin-kit/kit/types/redis/unknown"
	"github.com/stretchr/testify/assert"
)

func TestRedis_Client(t *testing.T) {
	r := New("my.redis")
	assert.Equal(t, r.Client(context.Background()).Get(context.Background(), "my.key").Err(), unknown.ErrUnknownRedisDB)
}
