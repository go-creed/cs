package rd

import (
	"github.com/go-redis/redis"
	log "github.com/micro/go-micro/v2/logger"

	"cs/public/config"
)

type RedisConfig struct {
	Addr    string `json:"addr"`
	Password string `json:"password"`
	Num      int    `json:"num"`
}

func initRedis() {
	var (
		c   = config.C()
		cfg = &RedisConfig{}
		err error
	)
	if err = c.Get("redis", cfg); err != nil {
		log.Fatal(err)
	}
	client = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		DB:       cfg.Num,
		Password: cfg.Password,
	})
	if err = client.Ping().Err(); err != nil {
		log.Fatal(err)
	}
	log.Info("redis init the connection success")
}
