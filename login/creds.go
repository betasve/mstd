package login

import (
	"time"
)

type Creds struct {
	authCallbackPath      string
	authCallbackHost      string
	clientId              string
	clientSecret          string
	permissions           string
	accessToken           string
	refreshToken          string
	accessTokenExpiresAt  time.Time
	refreshTokenExpiresAt time.Time
	loginDataCallbackFn   func(*AuthData)
}

func (c *Creds) SetAuthCallbackPath(path string) {
	c.authCallbackHost = path
}

func (c *Creds) SetAuthCallbackHost(host string) {
	c.authCallbackHost = host
}

func (c *Creds) SetId(id string) {
	c.clientId = id
}

func (c *Creds) SetSecret(secret string) {
	c.clientSecret = secret
}

func (c *Creds) SetPermissions(permissions string) {
	c.permissions = permissions
}

func (c *Creds) SetAccessToken(accessToken string) {
	c.accessToken = accessToken
}
func (c *Creds) SetRefreshToken(refreshToken string) {
	c.refreshToken = refreshToken
}

func (c *Creds) SetAccessTokenExpiresAt(accessTokenExpiresAt time.Time) {
	c.accessTokenExpiresAt = accessTokenExpiresAt
}

func (c *Creds) SetRefreshTokenExpiresAt(refreshTokenExpiresAt time.Time) {
	c.refreshTokenExpiresAt = refreshTokenExpiresAt
}

func (c *Creds) SetLoginDataCallbackFn(fn func(*AuthData)) {
	c.loginDataCallbackFn = fn
}
