package log

import (
	"testing"
	"time"

	"github.com/micro/go-micro/v2/logger"
	logM "github.com/micro/go-plugins/logger/logrus/v2"
	"github.com/sirupsen/logrus"
)

func TestName(t *testing.T) {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	l := logrus.New()
	l.SetFormatter(customFormatter)
	l.WithTime(time.Now()).Info("Hello Walrus")


	l2 := logM.NewLogger(logM.WithTextTextFormatter(customFormatter))


	logger.DefaultLogger = l2
	logger.Warn(123)
	logger.Fields(map[string]interface{}{})
	logger.Log(logger.InfoLevel, "testing: Info")
	logger.Logf(logger.InfoLevel, "testing: %s", "Infof")
}
