package handler

import (
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/gin-gonic/gin"
	log "github.com/micro/go-micro/v2/logger"

	"cs/plugin/db"
	"cs/public/conf"
)

var (
	client *clientv3.Client
)

func Init() {
	var err error
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{conf.RegistryAddress(), "" + "2380"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
}
func Get(ctx *gin.Context) {
	cfg := &conf.MysqlConfig{}
	err := conf.C().Get("mysql", cfg)
	if err != nil {
		log.Error(err)
	}
	fmt.Println(cfg)
}

func C(ctx *gin.Context) {
	config := db.MysqlConfig{
		User:     "root",
		Password: "root",
		Addr:     "localhost:3306",
		DbName:   "cs",
		LogMode:  true,
	}
	conf.C().Set("mysql", config)

}
