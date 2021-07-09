package server

import (
	"examresult/config"
	log "github.com/sirupsen/logrus"
)

func Init() {
	initEmail()
	initHttp()
}

func Run() {
	if err := GlobalEmailServer.DialTest(); err != nil {
		panic(err)
	}
	log.Info("GlobalEmailServer run")

	log.WithField("HttpAddress", config.HttpServer.HttpAddress).Info("HttpRouter run")
	err := HttpRouter.Run(config.HttpServer.HttpAddress)
	if err != nil {
		panic(err)
	}
}
