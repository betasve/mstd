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

// The `login` package is in charge of holding the login logic needed for
// working with MS' ToDo API. It's holding the knowledge of what requests to
// build and where to send them in order to retrieve the tokens we need so we
// can successfully retrieve and create lists and todo items onward in the app.
package login

import (
	"encoding/json"
	"fmt"
	httpService "github.com/betasve/mstd/ext/http"
	t "github.com/betasve/mstd/ext/time"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type AuthData struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var baseRequestUrl = "https://login.microsoftonline.com/common/oauth2/v2.0"
var authRequestPath = "/authorize"
var httpClient httpService.HttpClient = &httpService.Client{}
var tokenRequestPath = "/token"
var callbackFn func(string) error

const refreshTokenValidityInHours = 200 * 24

// Logs in a user.
// TODO: Add logout command to remove attributes from conf file
func (c *Creds) PerformLogin() error {
	if c.alreadyLoggedIn() {
		return c.refreshTokenIfNeeded()
	} else {
		return c.performLogin()
	}
}

// Checks if the user needs to log in or she is already logged in.
func (c *Creds) LoginNeeded() bool {
	return !c.isAccessTokenValid()
}

// Performs the login operation procedure.
func (c *Creds) performLogin() error {
	err := c.loginUrlHandlerFn(c.prepareLoginUrl())
	if err != nil {
		return err
	}

	return CallbackListen(c.authCallbackPath, c.getAccessToken)
}

// Retrieves the access token using MS' login API.
func (c *Creds) getAccessToken(authKey string) error {
	request, err := buildRequestObjectWithEncodedParams(
		baseRequestUrl+tokenRequestPath,
		c.buildRequestBodyForAuthToken(authKey).Encode(),
	)

	if err != nil {
		return err
	}

	return c.processTokenRequest(request)
}

// Retrieves the refresh token using MS' login API.
func (c *Creds) getRefreshToken() error {
	request, err := buildRequestObjectWithEncodedParams(
		baseRequestUrl+tokenRequestPath,
		c.buildRequestBodyForRefreshToken().Encode(),
	)

	if err != nil {
		return err
	}

	return c.processTokenRequest(request)
}

// Processes the received data to use it for our Config sructure.
func (c *Creds) processTokenRequest(request *http.Request) error {
	body, err := sendRequest(request)
	if err != nil {
		return err
	}

	a := AuthData{}

	if err = json.Unmarshal(body, &a); err != nil {
		return err
	}

	a.ExtExpiresIn = a.ExtExpiresIn * refreshTokenValidityInHours
	return c.loginDataCallbackFn(&a)
}

// Checks if the user is already logged in.
func (c *Creds) alreadyLoggedIn() bool {
	return c.isAccessTokenValid() || c.isRefreshTokenValid()
}

// Checks if the access token is still valid.
func (c *Creds) isAccessTokenValid() bool {
	return len(c.accessToken) != 0 &&
		t.Client.Now().Before(c.accessTokenExpiresAt)
}

// Checks if the refresh token is still valid.
func (c *Creds) isRefreshTokenValid() bool {
	return len(c.refreshToken) != 0 &&
		t.Client.Now().Before(c.refreshTokenExpiresAt)
}

// Refreshes the token if it's needed.
func (c *Creds) refreshTokenIfNeeded() error {
	if !c.isAccessTokenValid() {
		return c.getRefreshToken()
	}

	return nil
}

// Builds a request body(for receiving the auth token), holding the needed
// (documented in the API) url values.
func (c *Creds) buildRequestBodyForAuthToken(authKey string) url.Values {
	data := url.Values{}
	data.Set("client_id", c.clientId)
	data.Set("scope", c.permissions)
	data.Set("code", authKey)
	data.Set("redirect_uri", c.authCallbackHost+c.authCallbackPath)
	data.Set("grant_type", "authorization_code")
	data.Set("client_secret", c.clientSecret)

	return data
}

// Builds a request body(for receiving the refresh token), holding the needed
// (documented in the API) url values.
func (c *Creds) buildRequestBodyForRefreshToken() url.Values {
	data := url.Values{}
	data.Set("client_id", c.clientId)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", c.refreshToken)
	data.Set("client_secret", c.clientSecret)

	return data
}

// Builds a request object with encoded params.
func buildRequestObjectWithEncodedParams(requestUrl, urlEncodedParams string) (*http.Request, error) {
	req, err := http.NewRequest(
		"POST",
		requestUrl,
		strings.NewReader(urlEncodedParams),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(urlEncodedParams)))

	return req, nil
}

// Sends a (preliminarily prepared) request and returns its body bytes.
func sendRequest(req *http.Request) ([]byte, error) {
	res, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

// Spins up a tiny HTTP serer on :8008 to listen for a callback and handle the
// passed params.
func CallbackListen(callbackUrl string, cb func(string) error) error {
	callbackFn = cb

	httpClient.HandleFunc(callbackUrl, responder)
	err := httpClient.ListenAndServe(":8080", nil)

	if err != nil {
		return err
	}

	return nil
}

// The handler function for the tiny HTTP server. It's doing the actual
// handling of the request values.
func responder(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	code := strings.Join(values["code"], "")

	if code == "" {
		return
	}

	if err := callbackFn(code); err != nil {
		fmt.Fprint(w, err)
	} else {
		fmt.Fprint(
			w,
			"Successfully retrieved an authorization "+
				"code \nGo back to your console and check if login succeeded.",
		)
	}
}

// Assembles the url (and its params) we need to use in order to log
// with MS' API and receive the tokens we need.
func (c *Creds) prepareLoginUrl() string {
	urlParams := url.Values{}

	urlParams.Add("client_id", c.clientId)
	urlParams.Add("response_type", "code")
	urlParams.Add(
		"redirect_uri",
		c.authCallbackHost+
			c.authCallbackPath,
	)
	urlParams.Add("response_mode", "query")
	urlParams.Add("scope", c.permissions)

	return fmt.Sprintf(
		"%s%s?%s",
		baseRequestUrl,
		authRequestPath,
		urlParams.Encode(),
	)
}
