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

// The `todoapi` is the package responsible for holding the layer of abstraction
// that communicates with MS' API regarding the management of Lists and ToDos.
// It strives to remaing as isolated and unaware of the rest of the environment
// as possible. It aims to expose as little as possible. Namely only the
// structures that are used to hold the data, as well as the methods who
// wrap its retrieving from the API.
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

// Retrieves the collection of `ListItem`s.
func (ta *TodoApi) ListsIndex() (*[]ListsItem, error) {
	return retrieveLists(ta.token)
}

// Creates a ListItem setting its name.
func (ta *TodoApi) ListsCreate(name string) (*ListsItem, error) {
	return createAList(ta.token, name)
}

// Updates a ListItem finding it by its id and changing its name.
func (ta *TodoApi) ListsUpdate(id, name string) (*ListsItem, error) {
	return updateList(ta.token, id, name)
}

// Sets a token to be used for the API communication.
func (ta *TodoApi) SetToken(token string) {
	ta.token = token
}

// Gets the token that is used for the API communication.
func (ta *TodoApi) Token() string {
	return ta.token
}

// The function that is responsible for building the HTTP request and handling
// the response of the Lists API endpoint.
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

// The function that is responsible for building the HTTP request and handling
// the response of the 'Create a list' API endpoint.
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

// The function that is responsible for building the HTTP request and handling
// the response of the 'Update a list' API endpoint.
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

// A 'helper' function to construct requests with the needed (valid) Auth
// headers for successfully communicating with the API.
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
