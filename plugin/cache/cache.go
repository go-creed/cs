package cache

import (
	"sync"

	"github.com/go-redis/redis"
	//"github.com/gomodule/redigo/redis"
	log "github.com/micro/go-micro/v2/logger"
)

var (
	once   sync.Once
	err    error
	client *redis.Client
)

func Init() {
	once.Do(func() {
		connect()
	})
}

func Cache() *redis.Client {
	return client
}

func connect() {
	client = redis.NewClient(&redis.Options{
		Addr:     ":6379",
		DB:       0,
		Password: "",
	})
	if err = client.Ping().Err(); err != nil {
		log.Fatal(err)
	}
	log.Info("redis init the connection success")
}
