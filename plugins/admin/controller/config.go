package controller

import "github.com/compty001/go-admin/modules/config"

var Config config.Config

func SetConfig(cfg config.Config) {
	Config = cfg
}
