package log

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/micro/go-micro/v2/logger"
	"github.com/sirupsen/logrus"
)

const (
	defaultEsAddr  = "http://127.0.0.1:9200"
	defaultEsIndex = "go.test"
)

type defaultLogger struct {
	opts Options
}

func (d *defaultLogger) Init(opts ...Option) error {
	for _, o := range opts {
		o(&d.opts)
	}
	return nil
}

func (d *defaultLogger) Hook() logrus.LevelHooks {
	var err error
	esHook := &EsHook{}
	if d.opts.EsAddr == "" {
		d.opts.EsAddr = defaultEsAddr
	}
	if d.opts.EsIndex == "" {
		d.opts.EsIndex = defaultEsIndex
	}
	esHook.client, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{d.opts.EsAddr},
	})
	if err != nil {
		logger.Fatal(err)
	}
	esHook.Decode = &general{}
	esHook.Index = d.opts.EsIndex
	hooks := logrus.LevelHooks{}
	hooks.Add(esHook)
	return hooks
}

func (d *defaultLogger) Formatter() logrus.Formatter {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	return customFormatter
}

type general struct {
	Message string    `json:"message"`
	Time    time.Time `json:"timestamp"`
	File    string    `json:"file"`
	Line    int       `json:"line"`
}

func (g general) Marshal(entry *logrus.Entry) io.Reader {
	var (
		file string
		line int
	)

	for i := 5; i < 15; i++ {
		_, file, line, _ = runtime.Caller(i)
		if !strings.Contains(file, "github.com") {
			break
		}
	}
	g.Line = line
	g.File = file
	g.Message = entry.Message
	g.Time = entry.Time
	marshal, _ := json.Marshal(g)
	entry.Message = fmt.Sprintf("[%s:%d] \n\t%s\n", file, line, entry.Message)
	return strings.NewReader(string(marshal))
}
