package log

import (
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	once sync.Once
)

func Init(appName string) {
	once.Do(func() {
		initLogger(appName)
	})
}
