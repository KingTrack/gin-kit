package unknown

import "github.com/redis/go-redis/v9"

func NewClient() *redis.Client {
	opts := &redis.Options{
		Addr: "localhost:0",
		DB:   0,
	}
	return redis.NewClient(opts)
}
