package log

import (
	"io"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/micro/go-micro/v2/logger"
	logM "github.com/micro/go-plugins/logger/logrus/v2"
	"github.com/sirupsen/logrus"

	"cs/public/config"
)

type Reader interface {
	Marshal(entry *logrus.Entry) io.Reader
}

type LogConfig struct {
	EsAddr string `json:"es_addr"`
}

var (
	cfg = &LogConfig{}
)

func initLogger(appName string) {
	var (
		c   = config.C()
		err error
	)
	if err = c.Get("log", cfg); err != nil {
		logger.Fatal(err)
	}

	formatter := myTextFormatter()
	hook := initEsHook(appName)

	logger.DefaultLogger = logM.NewLogger(
		logM.WithTextTextFormatter(formatter),
		logM.WithLevelHooks(hook),
	)
}

func initEsHook(appName string) logrus.LevelHooks {
	var err error
	esHook := &EsHook{}
	esHook.client, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.EsAddr},
	})
	if err != nil {
		logger.Fatal(err)
	}
	esHook.Decode = &general{}
	esHook.Index = appName
	hooks := logrus.LevelHooks{}
	hooks.Add(esHook)
	return hooks
}

func myTextFormatter() *logrus.TextFormatter {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	return customFormatter
}
