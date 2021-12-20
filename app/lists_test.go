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
