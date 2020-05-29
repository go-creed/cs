package main

import (
	"os"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"

	"cs/app/auth-srv/handler"
	"cs/app/auth-srv/model"
	auth "cs/app/auth-srv/proto/auth"
	"cs/plugin/rd"
	_const "cs/public/const"
)

func main() {
	// Registry by etcd
	// 用读取环境变量的方式来做
	os.Getenv("")
	etcdRegistry := etcd.NewRegistry(
		registry.Addrs("127.0.0.1:2379"),
		registry.Timeout(time.Second*5),
	)
	// New Service
	service := micro.NewService(
		micro.Name(_const.AuthSrv),
		micro.Version("latest"),
		micro.Registry(etcdRegistry),
		micro.Address("127.0.0.1:12004"),
	)

	// Initialise service
	service.Init(
		micro.Action(func(context *cli.Context) error {
			// Init Redis
			rd.Init()
			// Init Model
			model.Init()
			// Init Handler
			handler.Init()
			return nil
		}),
	)

	// Register Handler
	auth.RegisterAuthHandler(service.Server(), new(handler.Auth))

	// Register Struct as Subscriber
	//micro.RegisterSubscriber(AuthSrv, service.Server(), new(subscriber.Auth))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
