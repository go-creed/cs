package main

import (
	"flag"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"

	"cs/app/upload-srv/conf"
	"cs/app/upload-srv/handler"
	"cs/app/upload-srv/model"
	upload "cs/app/upload-srv/proto/upload"
	"cs/plugin/db"
	cLog "cs/plugin/log"
	"cs/plugin/rd"
	"cs/plugin/trace"
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
	cLog.Init(
		cLog.SetEsIndex(conf.App().Log.EsIndex),
	)
}

func main() {
	// Init Config
	initCfg()
	// Registry etcd
	etcdRegistry := etcd.NewRegistry(
		registry.Timeout(5*time.Second),
		registry.Addrs("127.0.0.1:2379"),
	)
	// New Service
	service := micro.NewService(
		micro.Name(conf.App().Name),
		micro.Version(conf.App().Version),
		micro.Registry(etcdRegistry),
		micro.Address(conf.App().Address),
		micro.WrapHandler(trace.Wrapper()),
	)

	// Initialise service
	service.Init(
		micro.Action(func(c *cli.Context) error {
			// init db
			db.Init()
			// init redis
			rd.Init()
			// init model
			model.Init()
			// init gin-middleware
			handler.Init()
			return nil
		}),
	)

	// Register Handler
	upload.RegisterUploadHandler(service.Server(), new(handler.Upload))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
