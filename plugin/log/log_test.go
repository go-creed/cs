package log

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v6"
	//"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/micro/go-micro/v2/logger"
	logM "github.com/micro/go-plugins/logger/logrus/v2"
	"github.com/sirupsen/logrus"
)

func TestDefaultLogger(t *testing.T) {
	//initDefaultLogger()
	//logrus.SetLevel(logrus.DebugLevel)
	//logrus.Debug(1)
	newLogger := logger.NewLogger(logger.WithLevel(logger.DebugLevel))
	logger.DefaultLogger = newLogger
	//logger.WithLevel(logger.DebugLevel)
	logger.Debug(123)
	//k
	//logger.Info("123")
}
func TestName(t *testing.T) {
	hooks := logrus.LevelHooks{}
	e := &EsHook{}
	var err error
	e.client, err = elasticsearch.NewDefaultClient()
	e.Decode = &general{}
	if err != nil {
		t.Fatal(err)
	}
	if _, err = e.client.Info(); err != nil {
		t.Fatal(err)
	}
	hooks.Add(e)
	logger.DefaultLogger = logM.NewLogger(
		//logM.WithTextTextFormatter(myTextFormatter()),
		logM.WithLevelHooks(hooks),
	)
	logger.Warn()
}
