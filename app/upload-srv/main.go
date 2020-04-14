package main

import (
	"cs/app/upload-srv/handler"
	"cs/app/upload-srv/model"
	upload "cs/app/upload-srv/proto/upload"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"time"
)

func main() {
	// Registry etcd
	etctRegistry := etcd.NewRegistry(
		registry.Timeout(5*time.Second),
		registry.Addrs("127.0.0.1:2379"),
	)
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.cs.service.upload"),
		micro.Version("latest"),
		micro.Registry(etctRegistry),
		micro.Address("127.0.0.1:12000"),
	)

	// Initialise service
	service.Init(
		micro.Action(func(c *cli.Context) error {
			// init model
			model.Init()
			// init handler
			handler.Init()
			return nil
		}),
	)

	// Register Handler
	upload.RegisterUploadHandler(service.Server(), new(handler.Upload))

	// Register Struct as Subscriber
	//micro.RegisterSubscriber("go.micro.cs.service.upload", service.Server(), new(subscriber.Upload))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
