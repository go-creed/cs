package log

import (
	"github.com/micro/go-micro/v2/logger"
	logM "github.com/micro/go-plugins/logger/logrus/v2"
	"github.com/sirupsen/logrus"
)

type MysqlConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Addr     string `json:"addr"`     // including ip and port
	LogMode  bool   `json:"log_mode"` // 是否开启日志
	DbName   string `json:"db_name"`
}

func initLogger() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	//l := logrus.New()
	//l.SetFormatter(customFormatter)
	//l.WithTime(time.Now()).Info("Hello Walrus")
	logger.DefaultLogger = logM.NewLogger(logM.WithTextTextFormatter(customFormatter))
}
