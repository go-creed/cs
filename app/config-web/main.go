package main

import (
	"flag"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/cli/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"

	"cs/app/config-web/conf"
	"cs/app/config-web/handler"
	"cs/public/config"
)

var (
	configCenter = *flag.String("cc", "127.0.0.1:2379", "etcd's address")
	etcdCfg      config.EtcdConfig
)

func initCfg() {
	var err error
	config.Init(configCenter, conf.Init) //注册允许动态变化的配置
	etcdCfg, err = config.ETCD()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	initCfg()

	// registry by etcd
	etcdRegistry := etcd.NewRegistry(
		registry.Addrs(etcdCfg.Addrs),
		registry.Timeout(time.Duration(etcdCfg.Timeout)),
	)

	// create new web service
	service := web.NewService(
		web.Name(conf.Config().Name),
		web.Version(conf.Config().Version),
		web.Address(conf.Config().Address),
		web.Registry(etcdRegistry),
	)

	// initialise service
	if err := service.Init(
		web.Flags(&cli.StringSliceFlag{Name: "cc"}),
		web.Action(func(c *cli.Context) {
			// init handler
			handler.Init()
		}),
	); err != nil {
		log.Fatal(err)
	}

	// register call handler
	engine := gin.Default()
	config := engine.Group("/config")
	{
		config.GET("/all")
		config.PUT("/etcd", handler.C)
		config.GET("/etcd", handler.Get)
	}
	service.Handle("/", engine)
	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
