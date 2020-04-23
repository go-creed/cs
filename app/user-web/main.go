package main

import (
	"time"

	"cs/app/user-web/handler"
	"github.com/gin-gonic/gin"
	"github.com/micro/cli/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"
)

func main() {
	//registry by etcd
	etcdRegistry := etcd.NewRegistry(
		registry.Addrs("127.0.0.1:2379"),
		registry.Timeout(5*time.Second),
	)
	// create new web service
	service := web.NewService(
		web.Name("go.micro.cs.web.user"),
		web.Version("latest"),
		web.Registry(etcdRegistry),
		web.Address("127.0.0.1:12003"),
	)

	// initialise service
	if err := service.Init(
		web.Action(func(context *cli.Context) {
			//Init handler
			handler.Init()
		}),
	); err != nil {
		log.Fatal(err)
	}
	engine := gin.New()
	// register gin handler
	engine.POST("/login", handler.Login)
	engine.POST("/register", handler.Register)
	service.Handle("/", engine)
	// register html handler
	//service.Handle("/", http.FileServer(http.Dir("html")))

	// register call handler
	//service.HandleFunc("/user/call", handler.UserCall)

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
