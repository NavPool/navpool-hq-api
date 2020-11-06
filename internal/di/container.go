package di

import (
	"github.com/NavPool/navpool-hq-api/generated/dic"
	"github.com/sarulabs/dingo/v3"
	log "github.com/sirupsen/logrus"
)

var container *dic.Container

func Init() {
	container, _ = dic.NewContainer(dingo.App)
}

func Get() *dic.Container {
	if container == nil {
		log.Fatal("Container not initialised")
	}

	return container
}
