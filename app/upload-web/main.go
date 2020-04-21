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
)

func main() {
	// registry by etcd
	etcdRegistry := etcd.NewRegistry(
		registry.Addrs("127.0.0.1:2379"),
		registry.Timeout(5*time.Second),
	)
	// create new web service
	service := web.NewService(
		web.Name("go.micro.cs.web.upload"),
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
	file := engine.Group("/file", handler.Auth())
	file.POST("/upload", handler.FileUpload)
	file.GET("/detail", handler.FileDetail)
	//engine.StaticFS("/", http.Dir("./file"))
	// register html handler
	//service.Handle("/", http.FileServer(http.Dir("html")))

	// register call handler
	//service.HandleFunc("/upload/call", handler.UploadCall)

	// register by gin
	service.Handle("/", engine)

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
