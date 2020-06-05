package log

import (
	"io"

	"github.com/micro/go-micro/v2/logger"
	logM "github.com/micro/go-plugins/logger/logrus/v2"

	"github.com/sirupsen/logrus"

	"cs/public/config"
)

type Reader interface {
	Marshal(entry *logrus.Entry) io.Reader
}

type Logger interface {
	Init(opts ...Option) error
	Hook() logrus.LevelHooks
	Formatter(entry *logrus.Entry) logrus.Formatter
}

func initLogger(opts ...Option) {
	var (
		c   = config.C()
		err error
		opt Options
	)
	if err = c.Get("log", &opt); err != nil {
		logger.Error(err)
	}
	opts = append([]Option{SetOptions(opt)}, opts...)
	initDefaultLogger(opts...)
}

func initDefaultLogger(opts ...Option) {

	d := defaultLogger{}
	d.Init(opts...)

	// format
	var formatterOption logger.Option
	formatter := d.Formatter()
	switch formatter.(type) {
	case *logrus.TextFormatter:
		formatterOption = logM.WithTextTextFormatter(formatter.(*logrus.TextFormatter))
	case *logrus.JSONFormatter:
		formatterOption = logM.WithTextTextFormatter(formatter.(*logrus.TextFormatter))
	}
	// hook
	hook := logM.WithLevelHooks(d.Hook())
	// level
	level := logger.WithLevel(logger.DebugLevel)
	// new logger
	logger.DefaultLogger = logM.NewLogger(
		formatterOption,
		hook,
		level,
	)
}
