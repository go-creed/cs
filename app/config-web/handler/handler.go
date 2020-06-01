package handler

import (
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/gin-gonic/gin"
	log "github.com/micro/go-micro/v2/logger"

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

func C(ctx *gin.Context) {

}
