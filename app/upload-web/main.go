package main

import (
	"flag"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/cli/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"

	"cs/app/upload-web/conf"
	"cs/app/upload-web/handler"
	"cs/public/config"
	"cs/public/gin-middleware"
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
	// registry by etcd
	etcdRegistry := etcd.NewRegistry(
		registry.Addrs(etcdCfg.Addrs),
		registry.Timeout(time.Duration(etcdCfg.Timeout)),
	)
	// create new web service
	service := web.NewService(
		web.Name(conf.Config().Name),
		web.Version(conf.Config().Version),
		web.Registry(etcdRegistry),
		web.Address(conf.Config().Address),
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

	file := engine.Group("/file").
		Use(middleware.AuthWrapper(handler.Auth()))
	{
		file.POST("/upload", middleware.C(handler.FileUpload))
		file.GET("/detail", middleware.C(handler.FileDetail))
		file.GET("/chunk", middleware.C(handler.FileChunk))
		file.POST("/merge", middleware.C(handler.FileMerge))
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
