package handler

import (
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/gin-gonic/gin"
	log "github.com/micro/go-micro/v2/logger"

	"cs/plugin/db"
	"cs/public/config"
	_const "cs/public/const"
)

var (
	c *clientv3.Client
)

func Init() {
	var err error
	c, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{config.RegistryAddress(), "" + "2380"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
}
func Get(ctx *gin.Context) {
	cfg := &config.MysqlConfig{}
	err := config.C().Get("mysql", cfg)
	if err != nil {
		log.Error(err)
	}
	fmt.Println(cfg)
}

func C(ctx *gin.Context) {

	type appConfig struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		Port    int    `json:"port"`
		Version string `json:"version"`
	}
	config.C().Set(_const.UserSrv, appConfig{
		Name:    _const.UserSrv,
		Address: "localhost",
		Port:    12002,
		Version: "v1.0",
	})
	config.C().Set(_const.Etcd, config.EtcdConfig{
		Addrs:   "localhost:2379",
		Timeout: 5000,
	})
	config.C().Set("mysql", db.MysqlConfig{
		User:     "root",
		Password: "root",
		Addr:     "localhost:3306",
		DbName:   "cs",
		LogMode:  true,
	})

}
