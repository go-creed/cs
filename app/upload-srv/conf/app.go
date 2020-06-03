package conf

import (
	log "github.com/micro/go-micro/v2/logger"

	"cs/public/config"
	_const "cs/public/const"
)

type appConfig struct {
	Name    string `json:"name"`
	Address string `json:"address"` // include ip and port
	Version string `json:"version"`
}

var (
	app appConfig
	err error
	c   = config.C()
)

func Init() {
	if err = c.Get(_const.UploadSrv, &app); err != nil {
		log.Fatal(err)
	}
	log.Infof("APP【%s】 configuration of current service is %+v", _const.AuthSrv, app)
}

func App() appConfig {
	return app
}
