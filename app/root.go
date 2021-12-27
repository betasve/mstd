/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// The `app` package is holding logic for coordinating the features of the app
// itself. Holding its configuration too. It's faily isolated from the logic
// of doing the actual communication (part of the `todoapi` package). The `app`
// package aims to be the link between the commands in `cmd` and the communication
// results provided from the `todoapi`.
package app

import (
	"github.com/betasve/mstd/conf"
	"github.com/betasve/mstd/ext/log"
	api "github.com/betasve/mstd/todoapi"
)

var config *conf.Config
var CfgFilePath string
var apiClient api.TodoApiClient

// This is app's entry point. It's being invoked by the command-line tool
// that is being used. Here we read the config file from the path that's being
// set for it and initializing the configuration for the app.
func InitAppConfig() {
	config = &conf.Config{}
	if err := config.InitConfig(CfgFilePath); err != nil {
		log.Client.Fatal(err)
	}

	apiClient = &api.TodoApi{}
}
