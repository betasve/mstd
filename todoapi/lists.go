package todoapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	httpService "github.com/betasve/mstd/ext/http"
	"io"
	"io/ioutil"
	"net/http"
)

type TodoApi struct {
	token string
}

type TodoApiClient interface {
	ListsIndex() (*[]ListsItem, error)
	ListsCreate(string) (*ListsItem, error)
	ListsUpdate(string, string) (*ListsItem, error)
	SetToken(string)
	Token() string
}

type ListsItem struct {
	Id     string `json:"id"`
	Name   string `json:"displayName"`
	Owner  bool   `json:"isOwner"`
	Shared bool   `json:"isShared"`
	System string `json:"wellKnownListName"`
}

type ListsResponse struct {
	Context string      `json:"@odata.context"`
	Lists   []ListsItem `json:"value"`
}

type ContentType string

const (
	formCT ContentType = "application/x-www-form-urlencoded"
	jsonCT ContentType = "application/json"
)

const defaultPageSize string = "?$top=100"
const listsIndexEndpoint string = "https://graph.microsoft.com/v1.0/me/todo/lists/"

var httpClient httpService.HttpClient = &httpService.Client{}

func (ta *TodoApi) ListsIndex() (*[]ListsItem, error) {
	return retrieveLists(ta.token)
}

func (ta *TodoApi) ListsCreate(name string) (*ListsItem, error) {
	return createAList(ta.token, name)
}

func (ta *TodoApi) ListsUpdate(id, name string) (*ListsItem, error) {
	return updateList(ta.token, id, name)
}

func (ta *TodoApi) SetToken(token string) {
	ta.token = token
}

func (ta *TodoApi) Token() string {
	return ta.token
}

func retrieveLists(token string) (*[]ListsItem, error) {
	req, err := constructRequest(
		"GET",
		listsIndexEndpoint+defaultPageSize,
		token,
		nil,
		formCT,
	)

	if err != nil {
		return nil, err
	}

	res, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	listsResponse := ListsResponse{}
	if err = json.Unmarshal(body, &listsResponse); err != nil {
		return nil, err
	}

	return &listsResponse.Lists, nil
}

func createAList(token, name string) (*ListsItem, error) {
	jsonObj := []byte(fmt.Sprintf("{\"displayName\": \"%s\"}", name))

	req, err := constructRequest(
		"POST",
		listsIndexEndpoint,
		token,
		bytes.NewBuffer(jsonObj),
		jsonCT,
	)

	if err != nil {
		return nil, err
	}

	res, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if res.StatusCode != 201 {
		return nil,
			fmt.Errorf(
				"Unsuccessful request to To Do API:\n%d\n%s",
				res.StatusCode,
				string(body),
			)
	}

	if err != nil {
		return nil, err
	}

	listResponse := ListsItem{}
	if err = json.Unmarshal(body, &listResponse); err != nil {
		return nil, err
	}

	return &listResponse, nil
}

func updateList(token, id, name string) (*ListsItem, error) {
	jsonObj := []byte(fmt.Sprintf("{\"displayName\": \"%s\"}", name))

	req, err := constructRequest(
		"PUT",
		listsIndexEndpoint+id,
		token,
		bytes.NewBuffer(jsonObj),
		jsonCT,
	)

	if err != nil {
		return nil, err
	}

	res, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if res.StatusCode != 200 {
		return nil,
			fmt.Errorf(
				"Unsuccessful request to To Do API:\n%d\n%s",
				res.StatusCode,
				string(body),
			)
	}

	if err != nil {
		return nil, err
	}

	listResponse := ListsItem{}
	if err = json.Unmarshal(body, &listResponse); err != nil {
		return nil, err
	}

	return &listResponse, nil
}

func constructRequest(
	method, path, token string,
	body io.Reader,
	contentType ContentType,
) (*http.Request, error) {
	req, err := httpClient.NewRequest(method, path, body)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", string(contentType))
	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Bearer %s", token),
	)

	return req, nil
}
