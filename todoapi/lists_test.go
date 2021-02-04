package todoapi

import (
	"bytes"
	"errors"
	"fmt"
	httpService "github.com/betasve/mstd/ext/http/httptest"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

const listResponse1 string = `{
  "@odata.type": "#microsoft.graph.todoTaskList",
  "id": "1",
  "displayName": "List Title 1",
  "isOwner": true,
  "isShared": false,
  "wellknownListName": "none"
}`

const listResponse2 string = `{
				"@odata.type": "#microsoft.graph.todoTaskList",
				"id": "2",
        "displayName": "List Title 2",
        "isOwner": false,
        "isShared": true,
        "wellknownListName": "system"
        }`

func init() {
	httpClient = &httpService.ClientMock{}
}

func TestListsIndex(test *testing.T) {
	api := TodoApi{}
	api.SetToken("token")

	stubHttp(
		200,
		fmt.Sprintf(
			`{ "@odata.context": "some", "value": [%s, %s] }`,
			listResponse1,
			listResponse2,
		),
	)

	lists, err := api.ListsIndex()

	checkListsIndexExpectations(test, lists, err)
}

func TestListsCreate(test *testing.T) {
	api := TodoApi{}
	api.SetToken("token")

	stubHttp(201, listResponse1)

	listItem, err := api.ListsCreate("name")

	checkCreatedListExpectations(test, listItem, err)
}

func TestRetrieveListsSuccess(test *testing.T) {
	stubHttp(
		200,
		fmt.Sprintf(
			`{ "@odata.context": "some", "value": [%s, %s] }`,
			listResponse1,
			listResponse2,
		),
	)

	lists, err := retrieveLists("token")

	checkListsIndexExpectations(test, lists, err)
}

func TestRetrieveListsFailureUnmarshalling(test *testing.T) {
	stubHttp(
		200,
		fmt.Sprintf(
			`{ "@odata.context": "some", "value" [%s, %s] }`,
			listResponse1,
			listResponse2,
		),
	)

	_, err := retrieveLists("token")

	if err == nil {
		test.Errorf("\nExpected error\nbut got\nnil%s", err)
	}
}

func TestCreateAListSuccess(test *testing.T) {
	stubHttp(201, listResponse1)

	listItem, err := createAList("token", "name")

	checkCreatedListExpectations(test, listItem, err)
}

func TestCreateAListFailureWithWrongCode(test *testing.T) {
	stubHttp(304, listResponse1)

	_, err := createAList("token", "name")

	if err == nil {
		test.Error("\nExpected to return error\nbut it was\nnil")
	}
}

func TestConstructRequestSuccess(test *testing.T) {
	method := "GET"
	path := "some/path"
	token := "token"
	body := "body"

	req, _ := constructRequest(
		method,
		path,
		token,
		bytes.NewBuffer([]byte(body)),
		formCT,
	)

	if req.Method != method {
		test.Errorf("\nExpected method to be:\n%s\nbut was\n%s", method, req.Method)
	}

	if req.URL.Path != path {
		test.Errorf("\nExpected path to be:\n%s\nbut was\n%s", path, req.URL.Path)
	}

	if !strings.Contains(req.Header["Authorization"][0], fmt.Sprintf("Bearer %s", token)) {
		test.Errorf("\nExpected header to containt:\n%s\nbut\nit did not", token)
	}

	if req.Header["Content-Type"][0] != string(formCT) {
		test.Errorf("\nExpected header content type to be:\n%s\nbut was\n%s", string(formCT), req.Header["Content-Type"][0])
	}

	if reqBody, _ := ioutil.ReadAll(req.Body); string(reqBody) != body {
		test.Errorf("\nExpected body to be:\n%s\nbut was\n%s", body, reqBody)
	}
}

func TestConstructRequestFailure(test *testing.T) {
	httpClient = &httpService.ClientMock{}
	expectedErr := errors.New("error")

	httpService.NewRequestStubFn = func(method, url string, body io.Reader) (*http.Request, error) {
		return nil, expectedErr
	}

	method := "GET"
	path := "some/path"
	token := "token"
	body := "body"

	_, err := constructRequest(
		method,
		path,
		token,
		bytes.NewBuffer([]byte(body)),
		formCT,
	)

	if err == nil {
		test.Errorf("\nExpected error to be:\n%s\nbut got\nnil", expectedErr)
	}
}

func stubHttp(status int, body string) {
	httpService.MockFn = func(req *http.Request) (*http.Response, error) {
		res := &http.Response{}
		res.StatusCode = status
		res.Body = ioutil.NopCloser(
			strings.NewReader(body),
		)
		return res, nil
	}
}

func checkListsIndexExpectations(test *testing.T, lists *[]ListsItem, err error) {
	if err != nil {
		test.Errorf("\nExpected error to be:\nnil\nbut was\n%s", err)
	}

	if len(*lists) != 2 {
		test.Errorf("\nExpected a list of 2:\nbut got\n%d", len(*lists))
	}

	if (*lists)[0].Name != "List Title 1" {
		test.Errorf("\nExpected first list to have title:\nList Title 1\nbut got\n%s", (*lists)[0].Name)
	}
}

func checkCreatedListExpectations(test *testing.T, listItem *ListsItem, err error) {
	if err != nil {
		test.Errorf("\nExpected error to be:\nnil\nbut was\n%s", err)
	}

	if listItem.Id != "1" {
		test.Errorf("\nExpected id to be:\n1\nbut was\n%s", listItem.Id)
	}

	if listItem.Name != "List Title 1" {
		test.Errorf("\nExpected name to be:\nList Title 1\nbut was\n%s", listItem.Name)
	}

	if !listItem.Owner {
		test.Errorf("\nExpected owner to be:\ntrue\nbut was\n%v", listItem.Owner)
	}

	if listItem.Shared {
		test.Errorf("\nExpected shared to be:\nfalse\nbut was\n%v", listItem.Shared)
	}

	if listItem.System != "none" {
		test.Errorf("\nExpected system to be:\nnone\nbut was\n%v", listItem.System)
	}
}
