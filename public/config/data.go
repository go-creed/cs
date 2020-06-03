package config

import "encoding/json"

type Config struct {
	Apps map[string]interface{} `json:"apps"`
}

type AppConfig struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    string `json:"port"`
	Version string `json:"version"`
}

type EtcdConfig struct {
	Addrs   string `json:"addrs"`   // including ip and port
	Timeout int    `json:"timeout"` //超时时间
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
