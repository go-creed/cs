package trace

import "sync"

var (
	once sync.Once
)

func Init(name string) {
	once.Do(func() {
		initTrace(name)
	})
}
