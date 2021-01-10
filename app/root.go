package app

import (
	"github.com/betasve/mstd/conf"
	"github.com/betasve/mstd/ext/log"
)

var config *conf.Config
var CfgFilePath string

func InitAppConfig() {
	config = &conf.Config{}
	if err := config.InitConfig(CfgFilePath); err != nil {
		log.Client.Fatal(err)
	}
}
