package log

import (
	"errors"
	"fmt"
	"os"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/sirupsen/logrus"
)

type EsHook struct {
	client *elasticsearch.Client
	Index  string
	Decode Reader
}


func (hook *EsHook) Fire(entry *logrus.Entry) error {
	var (
		err   error
		index *esapi.Response
	)
	if index, err = hook.client.Index(hook.Index, hook.Decode.Marshal(entry)); err != nil {
		fmt.Fprintf(os.Stderr, "Create index failure, %v", err)
		return err
	} else if index.IsError() {
		fmt.Fprintf(os.Stderr, "Create index failure, %v", index.Body)
		_ = index.Body.Close()
		return errors.New(index.String())
	}
	return nil
}

func (hook *EsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
