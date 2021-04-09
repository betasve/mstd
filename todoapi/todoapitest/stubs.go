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
