package main

import (
	"flag"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"

	"cs/app/user-srv/conf"
	"cs/app/user-srv/handler"
	"cs/app/user-srv/model"
	user "cs/app/user-srv/proto/user"
	"cs/plugin/db"
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
		micro.Flags(&cli.StringSliceFlag{Name: "cc"}),
		micro.Action(func(c *cli.Context) error {
			// Init Mysql
			db.Init()
			// Init Model
			model.Init()
			// Init gin-middleware
			handler.Init()
			return nil
		}),
	)

	// Register Handler
	user.RegisterUserHandler(service.Server(), new(handler.User))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
