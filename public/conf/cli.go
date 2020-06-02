package conf

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/micro/cli/v2"
	configV2 "github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/etcd"
	log "github.com/micro/go-micro/v2/logger"
)

const (
	prefix = "/cs/app/config"
)

var (
	registryAddress string
	config          = &configurator{}
)

func Init(c *cli.Context) {
	var err error
	registryAddress = c.String("registry_address")
	if registryAddress == "" {
		registryAddress = "127.0.0.1:2379"
	}
	log.Infof("当前的 etcd %s", registryAddress)
	etcdSource := etcd.NewSource(
		etcd.WithAddress(registryAddress),
		etcd.WithPrefix("/cs/app"),
	)
	config.conf, err = configV2.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	config.client, err = clientv3.New(clientv3.Config{
		Endpoints: []string{
			registryAddress,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	if err = config.conf.Load(etcdSource); err != nil {
		log.Fatal(err)
	}
	watch, err := config.conf.Watch()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			next, err := watch.Next()
			if err != nil {
				log.Error(err)
				break
			}
			log.Infof("Watch changes : %v", string(next.Bytes()))
		}
	}()
}

func RegistryAddress() string {
	return registryAddress
}
