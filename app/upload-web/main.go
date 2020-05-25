package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/cli/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"

	"cs/app/upload-web/handler"
	_const "cs/public/const"
	"cs/public/gin-middleware"
)

func main() {
	// registry by etcd
	etcdRegistry := etcd.NewRegistry(
		registry.Addrs("127.0.0.1:2379"),
		registry.Timeout(5*time.Second),
	)
	// create new web service
	service := web.NewService(
		web.Name(_const.UploadWeb),
		web.Version("latest"),
		web.Registry(etcdRegistry),
		web.Address("127.0.0.1:12001"),
	)

	// initialise service
	if err := service.Init(
		web.Action(func(c *cli.Context) { //执行某一些初始化动作
			handler.Init()
		}),
	); err != nil {
		log.Fatal(err)
	}
	engine := gin.New()

	file := engine.Group("/file").Use(middleware.AuthWrapper(handler.Auth()))
	{
		file.POST("/upload", middleware.C(handler.FileUpload))
		file.GET("/detail", middleware.C(handler.FileDetail))
		file.GET("/chunk", middleware.C(handler.FileChunk))
	}

	//engine.StaticFS("/", http.Dir("./file"))
	// register html gin-middleware
	//service.Handle("/", http.FileServer(http.Dir("html")))

	// register call gin-middleware
	//service.HandleFunc("/upload/call", gin-middleware.UploadCall)

	// register by gin
	service.Handle("/", engine)

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
