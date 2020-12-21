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
	conf "github.com/betasve/mstd/conf"
	exec "github.com/betasve/mstd/exec"
	l "github.com/betasve/mstd/log"
	"github.com/betasve/mstd/runtime"
	t "github.com/betasve/mstd/time"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type auth struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var baseRequestUrl = "https://login.microsoftonline.com/common/oauth2/v2.0"
var authRequestPath = "/authorize"
var httpClient = &http.Client{}
var tokenRequestPath = "/token"
var callbackFn func(string)

// TODO: Add logout command to remove attributes from conf file
func Perform() {
	if alreadyLoggedIn() {
		refreshTokenIfNeeded()

		l.Client.Println("You are already logged in. If you want to log in anew, please use the logut command first.")
	} else {
		performLogin()
	}
}

func performLogin() {
	openLoginUrl(prepareLoginUrl())

	CallbackListen(conf.CurrentState.AuthCallbackPath, getAccessToken)
}

func prepareLoginUrl() string {
	urlParams := url.Values{}

	urlParams.Add("client_id", conf.CurrentState.ClientId)
	urlParams.Add("response_type", "code")
	urlParams.Add(
		"redirect_uri",
		conf.CurrentState.AuthCallbackHost+
			conf.CurrentState.AuthCallbackPath,
	)
	urlParams.Add("response_mode", "query")
	urlParams.Add("scope", conf.CurrentState.Permissions)

	return fmt.Sprintf(
		"%s%s?%s",
		baseRequestUrl,
		authRequestPath,
		urlParams.Encode(),
	)
}

func openLoginUrl(url string) {
	var err error

	switch runtime.Client.GetOS() {
	case "linux":
		err = exec.CmdClient.Command("xdg-open", url).Run()
	case "windows":
		err = exec.CmdClient.Command("rundll32", "url.dll,FileProtocolHandler", url).Run()
	case "darwin":
		err = exec.CmdClient.Command("open", url).Run()
	default:
		l.Client.Printf("Please visit \n\r %s \n\r and login.", url)
	}

	if err != nil {
		l.Client.Fatal(err)
	}
}

// TODO: To replace with a hook function
func writeDataToConfigFile(a *auth) {
	// TODO: Set those via externally with a callback func
	conf.SetClientAccessToken(a.AccessToken)
	conf.SetClientAccessTokenExpirySeconds(a.ExpiresIn)

	conf.SetClientRefreshToken(a.RefreshToken)
	expiryDaysToHours := 200 * 24
	day200, err := t.Client.ParseDuration(fmt.Sprintf("%dh", expiryDaysToHours))
	if err != nil {
		l.Client.Fatal(err)
	}
	conf.SetClientRefreshTokenExpirySeconds(int(day200.Seconds()))
}

func getAccessToken(authKey string) {
	request := buildRequestObjectWithEncodedParams(
		baseRequestUrl+tokenRequestPath,
		buildRequestBodyForAuthToken(authKey).Encode(),
	)

	// stubBody := []byte(`{"token_type":"Bearer","scope":"Tasks.ReadWrite.Shared Tasks.ReadWrite User.Read Mail.Read","expires_in":3600,"ext_expires_in":3600,"access_token":"EwBgA8l6BAAU6k7+XVQzkGyMv7VHB/h4cHbJYRAAAdP4XTitRFcSaCEkgaktzueLC4mJdOBqwzWA6AQ4BlMofDsqwJfswAoD8eXnuoP80RMgW5ZM9h6Qg7gFlzSnMKaGMf9wDa51GMGK6o4Gf/Miyik8MiDvCjIQU0mDIad8dEsYFfNv9Mq6h/aOCZVgLeMgA2c6Mqnyd3ym8UMRU/0+Olk3pJa1amRhXXtJCzCtPru5bEwJfjNcsM8pLux4tqP4WRxzWuwaCFdKqAHQc7PkBnpC1qgkr8jJh7v2cuLGAA+EDbqYi66EO+KxlaMe+wfbNkS36YBr3wwAOEdeu7K5zfK7pzge1SqSlQmxaSzZzRar0QhjkzoFbcjxFu6GMxsDZgAACG+mPO5XTOPCMAJA4uZFr7NTXI9IthKkUb+Dy31lUgT0V50sG/t8cRbI6fKOXpzzVgKLoNH+gcTaoRLAISq8mjwLBuBLU7eC5VoTInIDCNdQMYDzjPhj8SRVa8saBH/r4fuHMJfGAp0NEPyv3vPEH/ackLswawg9EUUxxjSgejawTmNP/H1UtGhPukfg6MVTpwA33N6E0urBEzwqANgtIXGMjDtfWKHGGUGFtBYSauftE7UAmukETjQD928Gyvm2Rq045AxmSJeRRlz1AUC6ffEf56jmEv/uOd4Eth+MGvwtsAA5wrCUXDjixFuyeghXmX8VduN+Q63+WoVjEl5f6XtrftSUeydSR40x6s9xFbvROiGssGxDK1hOCp8uF2fWyGl4x+x+d6vjt6t1Ha/L4duQkw4SxlQE58C22WTR0m4wtaN9zB1ovIrPtOf01+9kdkoowI017SoOBUQmUuSfGyjsLQnPvJi6Edp4KzFVHhJ/FfA5jvraVh4SgA7LyxiVCO1WIHYCXDamLqelIJ0W7mwcZB7snHH+gYXMACAzZcD1qf37O1QMuT6jUTUNe2jX5uh90QRFzYzOw819YX7rqb3Z8NytLSj+qkI9eetgpXplGnfipRloiXHBePVGAGrcigMLNc4Ny5DXIvqW0SJoQdLUqZUvzYKobf4BzFabr2rGLxU98zX/KKimw+IbBwQwzePdvdjPSWeZyUUhaX+TE0RiALAX82NnFz8dI6NXw/uj712ZIwtM1JJ9Dn4C","refresh_token":"M.R3_BAY.CVyjHHFi2Rrsv!vpgLLBhfs7SbGRhMu7TLdKA0wTBxsM1rX9Tggx8bzNizGx*vp5QdvZd8eP2hL5csx7BHhdZLwsHQ3CVfK9llk30wU1NKOiKoRuJThwudUNVsCkEZs2Xz53*Kb1RpErlHT44sVpwmh9ZFta3NXD70lJ4i2Jom1G7Ma8Ia4Ha149B0GtPpmdnlb7ENbHQAVEpkwBpZrJDDMG7PRrtLn3cG*C4QqtENtYUJbI!28JS378OQB1mMeEONEmVyrFz8nnwchGpNxY9JBo00uzh*12S3CwiDsiy2J3lYi*oQFNJsPhGbRmDhJTXo4ixtC!RULY1L8a33IVf7vmifKh!iaskVdDxGDJorcuW*Qxvt4ZC7gdl*18LHQBkcx7Rc3DLHxLLx!POTzI26FF5UV78B6LQnOOXYNRnSsd"}`)
	body := sendRequest(request)
	a := auth{}

	if err := json.Unmarshal(body, &a); err != nil {
		log.Fatal(err)
	}

	writeDataToConfigFile(&a)
	l.Client.Println("Logged in successfully.\nPlease Ctr+C to exit.")
}

func getRefreshToken() {
	request := buildRequestObjectWithEncodedParams(
		baseRequestUrl+tokenRequestPath,
		buildRequestBodyForRefreshToken().Encode(),
	)

	// stubBody := []byte(`{"token_type":"Bearer","scope":"Tasks.ReadWrite.Shared Tasks.ReadWrite User.Read Mail.Read","expires_in":3600,"ext_expires_in":3600,"access_token":"EwBgA8l6BAAU6k7+XVQzkGyMv7VHB/h4cHbJYRAAAdP4XTitRFcSaCEkgaktzueLC4mJdOBqwzWA6AQ4BlMofDsqwJfswAoD8eXnuoP80RMgW5ZM9h6Qg7gFlzSnMKaGMf9wDa51GMGK6o4Gf/Miyik8MiDvCjIQU0mDIad8dEsYFfNv9Mq6h/aOCZVgLeMgA2c6Mqnyd3ym8UMRU/0+Olk3pJa1amRhXXtJCzCtPru5bEwJfjNcsM8pLux4tqP4WRxzWuwaCFdKqAHQc7PkBnpC1qgkr8jJh7v2cuLGAA+EDbqYi66EO+KxlaMe+wfbNkS36YBr3wwAOEdeu7K5zfK7pzge1SqSlQmxaSzZzRar0QhjkzoFbcjxFu6GMxsDZgAACG+mPO5XTOPCMAJA4uZFr7NTXI9IthKkUb+Dy31lUgT0V50sG/t8cRbI6fKOXpzzVgKLoNH+gcTaoRLAISq8mjwLBuBLU7eC5VoTInIDCNdQMYDzjPhj8SRVa8saBH/r4fuHMJfGAp0NEPyv3vPEH/ackLswawg9EUUxxjSgejawTmNP/H1UtGhPukfg6MVTpwA33N6E0urBEzwqANgtIXGMjDtfWKHGGUGFtBYSauftE7UAmukETjQD928Gyvm2Rq045AxmSJeRRlz1AUC6ffEf56jmEv/uOd4Eth+MGvwtsAA5wrCUXDjixFuyeghXmX8VduN+Q63+WoVjEl5f6XtrftSUeydSR40x6s9xFbvROiGssGxDK1hOCp8uF2fWyGl4x+x+d6vjt6t1Ha/L4duQkw4SxlQE58C22WTR0m4wtaN9zB1ovIrPtOf01+9kdkoowI017SoOBUQmUuSfGyjsLQnPvJi6Edp4KzFVHhJ/FfA5jvraVh4SgA7LyxiVCO1WIHYCXDamLqelIJ0W7mwcZB7snHH+gYXMACAzZcD1qf37O1QMuT6jUTUNe2jX5uh90QRFzYzOw819YX7rqb3Z8NytLSj+qkI9eetgpXplGnfipRloiXHBePVGAGrcigMLNc4Ny5DXIvqW0SJoQdLUqZUvzYKobf4BzFabr2rGLxU98zX/KKimw+IbBwQwzePdvdjPSWeZyUUhaX+TE0RiALAX82NnFz8dI6NXw/uj712ZIwtM1JJ9Dn4C","refresh_token":"M.R3_BAY.CVyjHHFi2Rrsv!vpgLLBhfs7SbGRhMu7TLdKA0wTBxsM1rX9Tggx8bzNizGx*vp5QdvZd8eP2hL5csx7BHhdZLwsHQ3CVfK9llk30wU1NKOiKoRuJThwudUNVsCkEZs2Xz53*Kb1RpErlHT44sVpwmh9ZFta3NXD70lJ4i2Jom1G7Ma8Ia4Ha149B0GtPpmdnlb7ENbHQAVEpkwBpZrJDDMG7PRrtLn3cG*C4QqtENtYUJbI!28JS378OQB1mMeEONEmVyrFz8nnwchGpNxY9JBo00uzh*12S3CwiDsiy2J3lYi*oQFNJsPhGbRmDhJTXo4ixtC!RULY1L8a33IVf7vmifKh!iaskVdDxGDJorcuW*Qxvt4ZC7gdl*18LHQBkcx7Rc3DLHxLLx!POTzI26FF5UV78B6LQnOOXYNRnSsd"}`)
	body := sendRequest(request)
	a := auth{}

	if err := json.Unmarshal(body, &a); err != nil {
		log.Fatal(err)
	}

	log.Println(a.Scope)

	writeDataToConfigFile(&a)
}

func alreadyLoggedIn() bool {
	return isAccessTokenValid() || isRefreshTokenValid()
}

func isAccessTokenValid() bool {
	return len(conf.CurrentState.AccessToken) != 0 &&
		t.Client.Now().Before(conf.CurrentState.AccessTokenExpiresAt)
}

func isRefreshTokenValid() bool {
	return len(conf.CurrentState.RefreshToken) != 0 &&
		t.Client.Now().Before(conf.CurrentState.RefreshTokenExpiresAt)
}

func refreshTokenIfNeeded() {
	if !isAccessTokenValid() {
		getRefreshToken()
	}
}

// TODO: Cover with tests
func buildRequestBodyForAuthToken(authKey string) url.Values {
	data := url.Values{}
	data.Set("client_id", conf.CurrentState.ClientId)
	data.Set("scope", "Tasks.ReadWrite.Shared,offline_access")
	data.Set("code", authKey)
	data.Set("redirect_uri", conf.CurrentState.AuthCallbackHost+conf.CurrentState.AuthCallbackPath)
	data.Set("grant_type", "authorization_code")
	data.Set("client_secret", "0v8Ag1_FPYO70~l.Ect_G69v-qHmTDV~cN")

	return data
}

// TODO: Cover with tests
func buildRequestBodyForRefreshToken() url.Values {
	data := url.Values{}
	data.Set("client_id", conf.CurrentState.ClientId)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", conf.CurrentState.RefreshToken)
	data.Set("client_secret", conf.CurrentState.ClientSecret)

	return data
}

func buildRequestObjectWithEncodedParams(requestUrl, urlEncodedParams string) *http.Request {
	req, err := http.NewRequest(
		"POST",
		requestUrl,
		strings.NewReader(urlEncodedParams),
	)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(urlEncodedParams)))

	return req

}

func sendRequest(req *http.Request) []byte {
	res, err := httpClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(res.Status)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))

	return body
}

func CallbackListen(callbackUrl string, cb func(string)) {
	callbackFn = cb

	http.HandleFunc(callbackUrl, responder)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Panic(err)
	}
}

func responder(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	code := strings.Join(values["code"], "")

	fmt.Fprint(
		w,
		"Successfully retrieved an authorization "+
			"code \nGo back to your console and check if login succeeded.",
	)

	if code == "" {
		return
	}

	callbackFn(code)
}
