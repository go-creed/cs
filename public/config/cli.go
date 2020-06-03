package config

import (
	"github.com/coreos/etcd/clientv3"
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

type UpdateInit func()

func Init(registryAddress string, update ...UpdateInit) {
	var err error
	if registryAddress == "" {
		registryAddress = "127.0.0.1:2379"
	}
	log.Infof("当前的 etcd %s", registryAddress)
	etcdSource := etcd.NewSource(
		etcd.WithAddress(registryAddress),
		etcd.WithPrefix(prefix),
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
	for _, i := range update {
		i()
	}
	go func() {
		for {
			next, err := watch.Next()
			if err != nil {
				log.Error(err)
				break
			}
			log.Infof("Watch changes : %v", string(next.Bytes()))
			for _, i := range update {
				i()
			}
		}
	}()
}

func RegistryAddress() string {
	return registryAddress
}
