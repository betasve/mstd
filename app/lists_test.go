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
package app

import (
	"github.com/betasve/mstd/conf"
	api "github.com/betasve/mstd/todoapi"
	apiTest "github.com/betasve/mstd/todoapi/todoapitest"
	"testing"
)

func TestListIndex(test *testing.T) {
	config = &conf.Config{}
	_ = config.SetClientAccessToken("some")

	apiTest.ListsIndexMockFn = func() (*[]api.ListsItem, error) {
		return &[]api.ListsItem{
			api.ListsItem{
				Id:     "some-id",
				Name:   "example name",
				Owner:  true,
				Shared: false,
				System: "",
			},
		}, nil
	}
	apiClient = &apiTest.TodoApiMock{}
	err := ListsIndex([]string{"display name"})
	test.Log(err)
	if err != nil {
		test.Error(err)
	}
}
