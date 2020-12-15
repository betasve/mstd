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
	"encoding/json"
	"fmt"
	l "github.com/betasve/mstd/logger"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	uri "net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var url string

type auth struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// TODO: export host and port to settings
var callbackHost = "http://localhost:8080"
var callbackPath = "/login/authorized"
var CallbackUrl = callbackHost + callbackPath

var baseRequestUrl = "https://login.microsoftonline.com/common/oauth2/v2.0"
var authRequestPath = "/authorize"
var tokenRequestPath = "/token"

var appId = "b1a43d92-35c5-4654-ab80-1380211060a1"
var permissions = "Tasks.ReadWrite,offline_access"

// TODO: Add logout command to remove attributes from conf file
func Login() {
	if alreadyLogedIn() {
		l.Client.Log("You are already logged in. If you want to log in anew, please use the logut command first.")
		os.Exit(0)
	}

	// TODO: Interpolate string with variables
	redirectUri := uri.Values{}
	redirectUri.Add("redirect_uri", CallbackUrl)

	url =
		baseRequestUrl + authRequestPath +
			"?client_id=" + appId +
			"&response_type=code" +
			"&" + redirectUri.Encode() +
			"&response_mode=query" +
			"&scope=" + permissions +
			"&state=12345"

	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Run()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Run()
	case "darwin":
		err = exec.Command("open", url).Run()
	default:
		log.Printf("Please visit \n\r %s \n\r and login.", url)
	}

	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	CallbackListen(callbackPath, authCallbackFn)
}

func authCallbackFn(authKey string) string {
	endpoint := baseRequestUrl + tokenRequestPath
	data := uri.Values{}
	data.Set("client_id", appId)
	data.Set("scope", "Tasks.ReadWrite.Shared,offline_access")
	data.Set("code", authKey)
	data.Set("redirect_uri", callbackHost+callbackPath)
	data.Set("grant_type", "authorization_code")
	data.Set("client_secret", "0v8Ag1_FPYO70~l.Ect_G69v-qHmTDV~cN")

	client := &http.Client{}
	r, err := http.NewRequest(
		"POST",
		endpoint,
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		log.Fatal(err)
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)

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
	a := auth{}

	stubBody := []byte(`{"token_type":"Bearer","scope":"Tasks.ReadWrite.Shared Tasks.ReadWrite User.Read Mail.Read","expires_in":3600,"ext_expires_in":3600,"access_token":"EwBgA8l6BAAU6k7+XVQzkGyMv7VHB/h4cHbJYRAAAdP4XTitRFcSaCEkgaktzueLC4mJdOBqwzWA6AQ4BlMofDsqwJfswAoD8eXnuoP80RMgW5ZM9h6Qg7gFlzSnMKaGMf9wDa51GMGK6o4Gf/Miyik8MiDvCjIQU0mDIad8dEsYFfNv9Mq6h/aOCZVgLeMgA2c6Mqnyd3ym8UMRU/0+Olk3pJa1amRhXXtJCzCtPru5bEwJfjNcsM8pLux4tqP4WRxzWuwaCFdKqAHQc7PkBnpC1qgkr8jJh7v2cuLGAA+EDbqYi66EO+KxlaMe+wfbNkS36YBr3wwAOEdeu7K5zfK7pzge1SqSlQmxaSzZzRar0QhjkzoFbcjxFu6GMxsDZgAACG+mPO5XTOPCMAJA4uZFr7NTXI9IthKkUb+Dy31lUgT0V50sG/t8cRbI6fKOXpzzVgKLoNH+gcTaoRLAISq8mjwLBuBLU7eC5VoTInIDCNdQMYDzjPhj8SRVa8saBH/r4fuHMJfGAp0NEPyv3vPEH/ackLswawg9EUUxxjSgejawTmNP/H1UtGhPukfg6MVTpwA33N6E0urBEzwqANgtIXGMjDtfWKHGGUGFtBYSauftE7UAmukETjQD928Gyvm2Rq045AxmSJeRRlz1AUC6ffEf56jmEv/uOd4Eth+MGvwtsAA5wrCUXDjixFuyeghXmX8VduN+Q63+WoVjEl5f6XtrftSUeydSR40x6s9xFbvROiGssGxDK1hOCp8uF2fWyGl4x+x+d6vjt6t1Ha/L4duQkw4SxlQE58C22WTR0m4wtaN9zB1ovIrPtOf01+9kdkoowI017SoOBUQmUuSfGyjsLQnPvJi6Edp4KzFVHhJ/FfA5jvraVh4SgA7LyxiVCO1WIHYCXDamLqelIJ0W7mwcZB7snHH+gYXMACAzZcD1qf37O1QMuT6jUTUNe2jX5uh90QRFzYzOw819YX7rqb3Z8NytLSj+qkI9eetgpXplGnfipRloiXHBePVGAGrcigMLNc4Ny5DXIvqW0SJoQdLUqZUvzYKobf4BzFabr2rGLxU98zX/KKimw+IbBwQwzePdvdjPSWeZyUUhaX+TE0RiALAX82NnFz8dI6NXw/uj712ZIwtM1JJ9Dn4C","refresh_token":"M.R3_BAY.CVyjHHFi2Rrsv!vpgLLBhfs7SbGRhMu7TLdKA0wTBxsM1rX9Tggx8bzNizGx*vp5QdvZd8eP2hL5csx7BHhdZLwsHQ3CVfK9llk30wU1NKOiKoRuJThwudUNVsCkEZs2Xz53*Kb1RpErlHT44sVpwmh9ZFta3NXD70lJ4i2Jom1G7Ma8Ia4Ha149B0GtPpmdnlb7ENbHQAVEpkwBpZrJDDMG7PRrtLn3cG*C4QqtENtYUJbI!28JS378OQB1mMeEONEmVyrFz8nnwchGpNxY9JBo00uzh*12S3CwiDsiy2J3lYi*oQFNJsPhGbRmDhJTXo4ixtC!RULY1L8a33IVf7vmifKh!iaskVdDxGDJorcuW*Qxvt4ZC7gdl*18LHQBkcx7Rc3DLHxLLx!POTzI26FF5UV78B6LQnOOXYNRnSsd"}`)
	if err := json.Unmarshal(stubBody, &a); err != nil {
		log.Fatal(err)
	}

	log.Println(a.Scope)

	viper.Set("access_token", a.AccessToken)
	viper.Set("refresh_token", a.RefreshToken)
	if err := viper.WriteConfig(); err != nil {
		log.Println(err)
	}

	log.Println(viper.GetString("access_token"))
	if false {
		return refreshToken()
	} else {
		return "s"
	}
}

func refreshToken() string {
	endpoint := baseRequestUrl + tokenRequestPath
	data := uri.Values{}
	data.Set("client_id", appId)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", viper.GetString("refresh_token"))
	data.Set("client_secret", "0v8Ag1_FPYO70~l.Ect_G69v-qHmTDV~cN")

	client := &http.Client{}
	r, err := http.NewRequest(
		"POST",
		endpoint,
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		log.Fatal(err)
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)

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
	a := auth{}

	stubBody := []byte(`{"token_type":"Bearer","scope":"Tasks.ReadWrite.Shared Tasks.ReadWrite User.Read Mail.Read","expires_in":3600,"ext_expires_in":3600,"access_token":"EwBgA8l6BAAU6k7+XVQzkGyMv7VHB/h4cHbJYRAAAdP4XTitRFcSaCEkgaktzueLC4mJdOBqwzWA6AQ4BlMofDsqwJfswAoD8eXnuoP80RMgW5ZM9h6Qg7gFlzSnMKaGMf9wDa51GMGK6o4Gf/Miyik8MiDvCjIQU0mDIad8dEsYFfNv9Mq6h/aOCZVgLeMgA2c6Mqnyd3ym8UMRU/0+Olk3pJa1amRhXXtJCzCtPru5bEwJfjNcsM8pLux4tqP4WRxzWuwaCFdKqAHQc7PkBnpC1qgkr8jJh7v2cuLGAA+EDbqYi66EO+KxlaMe+wfbNkS36YBr3wwAOEdeu7K5zfK7pzge1SqSlQmxaSzZzRar0QhjkzoFbcjxFu6GMxsDZgAACG+mPO5XTOPCMAJA4uZFr7NTXI9IthKkUb+Dy31lUgT0V50sG/t8cRbI6fKOXpzzVgKLoNH+gcTaoRLAISq8mjwLBuBLU7eC5VoTInIDCNdQMYDzjPhj8SRVa8saBH/r4fuHMJfGAp0NEPyv3vPEH/ackLswawg9EUUxxjSgejawTmNP/H1UtGhPukfg6MVTpwA33N6E0urBEzwqANgtIXGMjDtfWKHGGUGFtBYSauftE7UAmukETjQD928Gyvm2Rq045AxmSJeRRlz1AUC6ffEf56jmEv/uOd4Eth+MGvwtsAA5wrCUXDjixFuyeghXmX8VduN+Q63+WoVjEl5f6XtrftSUeydSR40x6s9xFbvROiGssGxDK1hOCp8uF2fWyGl4x+x+d6vjt6t1Ha/L4duQkw4SxlQE58C22WTR0m4wtaN9zB1ovIrPtOf01+9kdkoowI017SoOBUQmUuSfGyjsLQnPvJi6Edp4KzFVHhJ/FfA5jvraVh4SgA7LyxiVCO1WIHYCXDamLqelIJ0W7mwcZB7snHH+gYXMACAzZcD1qf37O1QMuT6jUTUNe2jX5uh90QRFzYzOw819YX7rqb3Z8NytLSj+qkI9eetgpXplGnfipRloiXHBePVGAGrcigMLNc4Ny5DXIvqW0SJoQdLUqZUvzYKobf4BzFabr2rGLxU98zX/KKimw+IbBwQwzePdvdjPSWeZyUUhaX+TE0RiALAX82NnFz8dI6NXw/uj712ZIwtM1JJ9Dn4C","refresh_token":"M.R3_BAY.CVyjHHFi2Rrsv!vpgLLBhfs7SbGRhMu7TLdKA0wTBxsM1rX9Tggx8bzNizGx*vp5QdvZd8eP2hL5csx7BHhdZLwsHQ3CVfK9llk30wU1NKOiKoRuJThwudUNVsCkEZs2Xz53*Kb1RpErlHT44sVpwmh9ZFta3NXD70lJ4i2Jom1G7Ma8Ia4Ha149B0GtPpmdnlb7ENbHQAVEpkwBpZrJDDMG7PRrtLn3cG*C4QqtENtYUJbI!28JS378OQB1mMeEONEmVyrFz8nnwchGpNxY9JBo00uzh*12S3CwiDsiy2J3lYi*oQFNJsPhGbRmDhJTXo4ixtC!RULY1L8a33IVf7vmifKh!iaskVdDxGDJorcuW*Qxvt4ZC7gdl*18LHQBkcx7Rc3DLHxLLx!POTzI26FF5UV78B6LQnOOXYNRnSsd"}`)
	if err := json.Unmarshal(stubBody, &a); err != nil {
		log.Fatal(err)
	}

	log.Println(a.Scope)

	viper.Set("access_token", a.AccessToken)
	viper.Set("refresh_token", a.RefreshToken)
	if err := viper.WriteConfig(); err != nil {
		log.Println(err)
	}

	log.Println(viper.GetString("access_token"))
	return "s"
}

func alreadyLogedIn() bool {
	return true
}