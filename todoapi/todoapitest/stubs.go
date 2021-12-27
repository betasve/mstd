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
package todoapitest

import (
	api "github.com/betasve/mstd/todoapi"
)

type TodoApiMock struct {
	token string
}

var ListsIndexMockFn = func() (*[]api.ListsItem, error) {
	return &[]api.ListsItem{}, nil
}

var ListsCreateMockFn = func(n string) (*api.ListsItem, error) {
	return &api.ListsItem{}, nil
}

var ListsUpdateMockFn = func(i, n string) (*api.ListsItem, error) {
	return &api.ListsItem{}, nil
}

func (ta *TodoApiMock) ListsIndex() (*[]api.ListsItem, error) {
	return ListsIndexMockFn()
}

func (ta *TodoApiMock) ListsCreate(name string) (*api.ListsItem, error) {
	return ListsCreateMockFn(name)
}

func (ta *TodoApiMock) ListsUpdate(id, name string) (*api.ListsItem, error) {
	return ListsUpdateMockFn(id, name)
}

func (ta *TodoApiMock) SetToken(token string) {
	ta.token = token
}

func (ta *TodoApiMock) Token() string {
	return ta.token
}
