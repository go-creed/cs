package upload

import (
	"sync"

	log "github.com/micro/go-micro/v2/logger"
)

var (
	once sync.Once
	s    *service
)

// Init Service Model Like Redis, Mysql ....
func Init() {
	once.Do(func() {
		log.Info("[Upload][Model] init service model like redis,mysql...")
		s = &service{} //init service
	})
}
