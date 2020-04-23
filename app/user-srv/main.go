package main

import (
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"

	"cs/app/user-srv/handler"
	"cs/app/user-srv/model"
	user "cs/app/user-srv/proto/user"
	"cs/plugin/db"
)

func main() {
	// Registry by etcd
	etcdRegistry := etcd.NewRegistry(
		registry.Addrs("127.0.0.1:2379"),
		registry.Timeout(5*time.Second),
	)
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.cs.service.user"),
		micro.Version("latest"),
		micro.Registry(etcdRegistry),
		micro.Address("127.0.0.1:12002"),
	)

	// Initialise service
	service.Init(
		micro.Action(func(context *cli.Context) error {
			// Init Mysql
			db.Init()
			// Init Model
			model.Init()
			// Init handler
			handler.Init()
			return nil
		}),
	)

	// Register Handler
	user.RegisterUserHandler(service.Server(), new(handler.User))

	// Register Struct as Subscriber
	//micro.RegisterSubscriber("go.micro.cs.service.user", service.Server(), new(subscriber.User))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
