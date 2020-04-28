package cache

import (
	"sync"

	"github.com/gomodule/redigo/redis"
	log "github.com/micro/go-micro/v2/logger"
)

var (
	once sync.Once
	err  error
	dial redis.Conn
)

func Init() {
	once.Do(func() {
		connect()
	})
}

func Cache() redis.Conn {
	return dial
}

func connect() {
	dial, err = redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}
	if do, err := dial.Do("ping"); err != nil {
		log.Fatal(err)
	} else if do.(string) != "PONG" {
		log.Fatal("redis ping failure")
	}
	log.Info("redis init the connection success")
}
