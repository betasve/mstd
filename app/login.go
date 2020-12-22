package app

import (
	"fmt"
	"github.com/betasve/mstd/conf"
	"github.com/betasve/mstd/log"
	"github.com/betasve/mstd/login"
	t "github.com/betasve/mstd/time"
)

func Login() {
	creds := login.Creds{}

	creds.SetAuthCallbackHost(conf.CurrentState.AuthCallbackHost)
	creds.SetAuthCallbackPath(conf.CurrentState.AuthCallbackPath)
	creds.SetId(conf.CurrentState.ClientId)
	creds.SetSecret(conf.CurrentState.ClientSecret)
	creds.SetPermissions(conf.CurrentState.Permissions)
	creds.SetAccessToken(conf.CurrentState.AccessToken)
	creds.SetAccessTokenExpiresAt(conf.CurrentState.AccessTokenExpiresAt)
	creds.SetRefreshToken(conf.CurrentState.RefreshToken)
	creds.SetRefreshTokenExpiresAt(conf.CurrentState.RefreshTokenExpiresAt)
	creds.SetLoginDataCallbackFn(writeDataToConfigFile)

	if err := creds.Perform(); err != nil {
		log.Client.Fatal(err)
	} else {
		log.Client.Println("Logged in successfully.\nPlease Ctr+C to exit.")
	}
}

func writeDataToConfigFile(a *login.AuthData) {
	conf.SetClientAccessToken(a.AccessToken)
	conf.SetClientAccessTokenExpirySeconds(a.ExpiresIn)

	conf.SetClientRefreshToken(a.RefreshToken)
	expiryDaysToHours := 200 * 24
	// TODO: Move token duration to a more approriate place
	day200, err := t.Client.ParseDuration(fmt.Sprintf("%dh", expiryDaysToHours))
	if err != nil {
		log.Client.Fatal(err)
	}
	conf.SetClientRefreshTokenExpirySeconds(int(day200.Seconds()))
}
