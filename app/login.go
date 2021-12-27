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
	"fmt"
	"github.com/betasve/mstd/ext/exec"
	"github.com/betasve/mstd/ext/log"
	"github.com/betasve/mstd/ext/runtime"
	"github.com/betasve/mstd/login"
)

var creds login.Creds

// Facilitates the login procedure for the app. Mainly - reads credentials from
// the config file and saves them in a place they can be easily accessed.
func Login() {
	creds = login.Creds{}

	creds.SetAuthCallbackHost(config.AuthCallbackHost())
	creds.SetAuthCallbackPath(config.AuthCallbackPath())
	creds.SetId(config.ClientId())
	creds.SetSecret(config.ClientSecret())
	creds.SetPermissions(config.ClientPermissions())
	creds.SetAccessToken(config.ClientAccessToken())
	creds.SetAccessTokenExpiresAt(config.ClientAccessTokenExpiresAt())
	creds.SetRefreshToken(config.ClientRefreshToken())
	creds.SetRefreshTokenExpiresAt(config.ClientRefreshTokenExpiresAt())
	creds.SetLoginDataCallbackFn(writeDataToConfigFile)
	creds.SetLoginUrlHandlerFn(openLoginUrl)

	if err := creds.PerformLogin(); err != nil {
		log.Client.Fatal(err)
	}
}

// Checks if the user needs to be logged in (again) or his current session is
// still active.
func LoginNeeded() bool {
	return creds.LoginNeeded()
}

// Writes data to the config file for the app.
// TODO: Update tests by covering the api token setup
func writeDataToConfigFile(a *login.AuthData) error {
	err := config.SetClientAccessToken(a.AccessToken)
	if err != nil {
		return err
	}

	apiClient.SetToken(a.AccessToken)

	err = config.SetClientAccessTokenExpirySeconds(a.ExpiresIn)
	if err != nil {
		return err
	}

	err = config.SetClientRefreshToken(a.RefreshToken)
	if err != nil {
		return err
	}
	err = config.SetClientRefreshTokenExpirySeconds(a.ExtExpiresIn)
	if err != nil {
		return err
	}

	log.Client.Println("Logged in successfully.\nPlease Ctr+C to exit.")
	return nil
}

// Opens the login url (provided as an argument) with the default browser for
// the operating system the app is being run on.
func openLoginUrl(url string) error {
	var err error

	switch runtime.Client.GetOS() {
	case "linux":
		err = exec.CmdClient.Command("xdg-open", url).Run()
	case "windows":
		err = exec.CmdClient.Command("rundll32", "url.dll,FileProtocolHandler", url).Run()
	case "darwin":
		err = exec.CmdClient.Command("open", url).Run()
	default:
		err = fmt.Errorf("OS not recognized. Please visit \n\r%s\n\r and login.", url)
	}

	if err != nil {
		return err
	}

	return nil
}
