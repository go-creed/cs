package main

import (
	"flag"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"

	"cs/app/auth-srv/conf"
	"cs/app/auth-srv/handler"
	"cs/app/auth-srv/model"
	auth "cs/app/auth-srv/proto/auth"
	"cs/plugin/rd"
	"cs/public/config"
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
	// Init Config
	initCfg()
	// Registry by etcd
	// 用读取环境变量的方式来做
	etcdRegistry := etcd.NewRegistry(
		registry.Addrs(etcdCfg.Addrs),
		registry.Timeout(time.Duration(etcdCfg.Timeout)),
	)
	// New Service
	service := micro.NewService(
		micro.Name(conf.Config().Name),
		micro.Version(conf.Config().Version),
		micro.Registry(etcdRegistry),
		micro.Address(conf.Config().Address),
	)

	// Initialise service
	service.Init(
		micro.Action(func(c *cli.Context) error {
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
