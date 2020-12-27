package app

import (
	"fmt"
	"github.com/betasve/mstd/log"
	"github.com/betasve/mstd/login"
	t "github.com/betasve/mstd/time"
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

	if err := creds.Perform(); err != nil {
		log.Client.Fatal(err)
	} else {
		log.Client.Println("Logged in successfully.\nPlease Ctr+C to exit.")
	}
}

func writeDataToConfigFile(a *login.AuthData) {
	config.SetClientAccessToken(a.AccessToken)
	config.SetClientAccessTokenExpirySeconds(a.ExpiresIn)

	config.SetClientRefreshToken(a.RefreshToken)
	expiryDaysToHours := 200 * 24
	// TODO: Move token duration to a more approriate place
	day200, err := t.Client.ParseDuration(fmt.Sprintf("%dh", expiryDaysToHours))
	if err != nil {
		log.Client.Fatal(err)
	}
	config.SetClientRefreshTokenExpirySeconds(int(day200.Seconds()))
}
