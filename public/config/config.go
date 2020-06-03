package config

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	configV2 "github.com/micro/go-micro/v2/config"
	log "github.com/micro/go-micro/v2/logger"

	_const "cs/public/const"
)

type Configurator interface {
	Get(name string, config interface{}) (err error)
	Set(name string, config interface{}) (err error)
}

type configurator struct {
	conf   configV2.Config
	client *clientv3.Client
}

func nameToName(name string) string {
	return prefix + "/" + name
}
func nameToPath(name string) []string {
	var rst []string
	vals := strings.Split(nameToName(name), "/")
	for _, val := range vals {
		if val != "" {
			rst = append(rst, val)
		}
	}
	log.Info(rst)
	return rst
}

func (c *configurator) Get(name string, config interface{}) (err error) {
	v := c.conf.Get(nameToPath(name)...)
	if v != nil {
		err = v.Scan(config)
	}
	return
}

func (c *configurator) Set(name string, config interface{}) error {
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*3)
	marshal, _ := json.Marshal(config)
	_, err := c.client.Put(ctx, nameToName(name), string(marshal))
	return err
}

func C() Configurator {
	return config
}

func ETCD() (etcd EtcdConfig, err error) {
	err = config.Get(_const.Etcd, &etcd)
	return etcd, err
}

