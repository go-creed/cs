package rd

import (
	"sync"

	"github.com/go-redis/redis"
)

var (
	once   sync.Once
	err    error
	client *redis.Client
)

func Init() {
	once.Do(func() {
		initRedis()
	})
}

func Cache() *redis.Client {
	return client
}
