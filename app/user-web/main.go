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

	"cs/app/user-web/conf"
	"cs/app/user-web/handler"
	"cs/public/config"
	middleware "cs/public/gin-middleware"
)

var (
	configCenter = *flag.String("cc", "127.0.0.1:2379", "")
	etcdCfg      config.EtcdConfig
)

func initCfg() {
	flag.Parse()
	var err error
	config.Init(configCenter, conf.Init)
	etcdCfg, err = config.ETCD()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initCfg()
	//registry by etcd
	etcdRegistry := etcd.NewRegistry(
		registry.Addrs(etcdCfg.Addrs),
		registry.Timeout(time.Duration(etcdCfg.Timeout)),
	)
	// create new web service
	service := web.NewService(
		web.Name(conf.App().Name),
		web.Version(conf.App().Version),
		web.Registry(etcdRegistry),
		web.Address(conf.App().Address),
	)

	// initialise service
	if err := service.Init(
		web.Action(func(context *cli.Context) {
			//Init gin-middleware
			handler.Init()
		}),
	); err != nil {
		log.Fatal(err)
	}
	engine := gin.Default()
	gin.Default()
	// register gin gin-middleware
	engine.POST("/login", handler.Login)
	engine.POST("/register", handler.Register)
	auth := engine.Use(middleware.AuthWrapper(handler.Auth()))
	{
		auth.GET("/detail", handler.Detail)
	}
	service.Handle("/", engine)

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
