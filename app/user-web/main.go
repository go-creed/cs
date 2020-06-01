package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/cli/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"

	"cs/app/user-web/handler"
	_const "cs/public/const"
	middleware "cs/public/gin-middleware"
)

func main() {
	//registry by etcd
	etcdRegistry := etcd.NewRegistry(
		registry.Addrs("127.0.0.1:2379"),
		registry.Timeout(5*time.Second),
	)
	// create new web service
	service := web.NewService(
		web.Name(_const.UserWeb),
		web.Version("latest"),
		web.Registry(etcdRegistry),
		web.Address("127.0.0.1:12003"),
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
	engine := gin.New()
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
