package conf

import (
	log "github.com/micro/go-micro/v2/logger"

	cLog "cs/plugin/log"
	"cs/public/config"
	_const "cs/public/const"
)

type appConfig struct {
	Name    string      `json:"name"`
	Address string      `json:"address"` // include ip and port
	Version string      `json:"version"`
	Log     cLog.Options `json:"log"`
}

var (
	app appConfig
	err error
	c   = config.C()
)

func Init() {
	if err = c.Get(_const.UserWeb, &app); err != nil {
		log.Fatal(err)
	}
	log.Infof("APP【%s】 configuration of current service is %+v", _const.UserWeb, app)
}

func App() appConfig {
	return app
}
