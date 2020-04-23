package cache

import "sync"

var (
	once sync.Once
	err  error
)

func Init() {

}

func o() {}
