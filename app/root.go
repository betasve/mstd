package app

import (
	"github.com/betasve/mstd/conf"
	"github.com/betasve/mstd/ext/log"
	api "github.com/betasve/mstd/todoapi"
)

var config *conf.Config
var CfgFilePath string
var apiClient api.TodoApiClient

func InitAppConfig() {
	config = &conf.Config{}
	if err := config.InitConfig(CfgFilePath); err != nil {
		log.Client.Fatal(err)
	}

	apiClient = &api.TodoApi{}
}
