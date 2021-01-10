package app

import (
	"fmt"
	"github.com/betasve/mstd/ext/exec"
	"github.com/betasve/mstd/ext/log"
	"github.com/betasve/mstd/ext/runtime"
	"github.com/betasve/mstd/login"
)

func Login() {
	creds := login.Creds{}

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

func writeDataToConfigFile(a *login.AuthData) error {
	err := config.SetClientAccessToken(a.AccessToken)
	if err != nil {
		return err
	}

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
