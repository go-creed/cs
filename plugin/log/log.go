package log

import (
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	once sync.Once
	err  error
)

func Init() {
	once.Do(func() {
		initLogger()
	})
}
