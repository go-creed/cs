package conf

type Config struct {
	Mysql MysqlConfig `json:"mysql"`
	Etcd  EtcdConfig  `json:"etcd"`
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
