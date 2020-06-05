package log

import (
	"sync"
)

var (
	once sync.Once
)

func Init(option ...Option) {
	once.Do(func() {
		initLogger(option...)
	})
}
