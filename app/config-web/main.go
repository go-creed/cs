package main

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/cli/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/web"

	"cs/app/config-web/handler"
	"cs/public/conf"
	_const "cs/public/const"
)

func main() {

	// create new web service
	service := web.NewService(
		web.Name(_const.ConfigWeb),
		web.Version(_const.VersionLatest),
		web.Address("localhost:12888"),
	)

	// initialise service
	if err := service.Init(
		web.Action(func(c *cli.Context) {
			conf.Init(c)
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
	}
	service.Handle("/", engine)
	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
