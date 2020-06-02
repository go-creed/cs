package conf

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	configV2 "github.com/micro/go-micro/v2/config"
	log "github.com/micro/go-micro/v2/logger"
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

type Config struct {
	Apps map[string]interface{} `json:"apps"`
}

type AppConfig struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    string `json:"port"`
	Version string `json:"version"`
}

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	Num      int    `json:"num"`
}
type EtcdConfig struct {
	Addr string `json:"addr"` // including ip and port
}

type MysqlConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Addr     string `json:"addr"`     // including ip and port
	LogMode  bool   `json:"log_mode"` // 是否开启日志
}

func (m MysqlConfig) String() string {
	marshal, _ := json.Marshal(m)
	return string(marshal)
}
func (e EtcdConfig) String() string {
	marshal, _ := json.Marshal(e)
	return string(marshal)
}

func (c Config) String() string {
	marshal, _ := json.Marshal(c)
	return string(marshal)
}
