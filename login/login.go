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
package login

import (
	"encoding/json"
	"fmt"
	httpService "github.com/betasve/mstd/http"
	t "github.com/betasve/mstd/time"
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

// TODO: Add logout command to remove attributes from conf file
func (c *Creds) PerformLogin() error {
	if c.alreadyLoggedIn() {
		return c.refreshTokenIfNeeded()
	} else {
		return c.performLogin()
	}
}

func (c *Creds) performLogin() error {
	err := c.loginUrlHandlerFn(c.prepareLoginUrl())
	if err != nil {
		return err
	}

	return CallbackListen(c.authCallbackPath, c.getAccessToken)
}

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

func (c *Creds) alreadyLoggedIn() bool {
	return c.isAccessTokenValid() || c.isRefreshTokenValid()
}

func (c *Creds) isAccessTokenValid() bool {
	return len(c.accessToken) != 0 &&
		t.Client.Now().Before(c.accessTokenExpiresAt)
}

func (c *Creds) isRefreshTokenValid() bool {
	return len(c.refreshToken) != 0 &&
		t.Client.Now().Before(c.refreshTokenExpiresAt)
}

func (c *Creds) refreshTokenIfNeeded() error {
	if !c.isAccessTokenValid() {
		return c.getRefreshToken()
	}

	return nil
}

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

func (c *Creds) buildRequestBodyForRefreshToken() url.Values {
	data := url.Values{}
	data.Set("client_id", c.clientId)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", c.refreshToken)
	data.Set("client_secret", c.clientSecret)

	return data
}

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

func CallbackListen(callbackUrl string, cb func(string) error) error {
	callbackFn = cb

	httpClient.HandleFunc(callbackUrl, responder)
	err := httpClient.ListenAndServe(":8080", nil)

	if err != nil {
		return err
	}

	return nil
}

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
