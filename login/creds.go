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
	loginDataCallbackFn   func(*AuthData) error
	loginUrlHandlerFn     func(string) error
}

// A setter method for authCallbackPath.
func (c *Creds) SetAuthCallbackPath(path string) {
	c.authCallbackPath = path
}

// A setter method for authCallbackHost.
func (c *Creds) SetAuthCallbackHost(host string) {
	c.authCallbackHost = host
}

// A setter method for clientId.
func (c *Creds) SetId(id string) {
	c.clientId = id
}

// A setter method for clientSecret.
func (c *Creds) SetSecret(secret string) {
	c.clientSecret = secret
}

// A setter method for permissions.
func (c *Creds) SetPermissions(permissions string) {
	c.permissions = permissions
}

// A setter method for accessToken.
func (c *Creds) SetAccessToken(accessToken string) {
	c.accessToken = accessToken
}

// A setter method for refreshToken.
func (c *Creds) SetRefreshToken(refreshToken string) {
	c.refreshToken = refreshToken
}

// A setter method for accessTokenExpiresAt.
func (c *Creds) SetAccessTokenExpiresAt(accessTokenExpiresAt time.Time) {
	c.accessTokenExpiresAt = accessTokenExpiresAt
}

// A setter method for refreshTokenExpiresAt.
func (c *Creds) SetRefreshTokenExpiresAt(refreshTokenExpiresAt time.Time) {
	c.refreshTokenExpiresAt = refreshTokenExpiresAt
}

// A setter method for loginDataCallbackFn.
func (c *Creds) SetLoginDataCallbackFn(fn func(*AuthData) error) {
	c.loginDataCallbackFn = fn
}

// A setter method for loginUrlHandlerFn.
func (c *Creds) SetLoginUrlHandlerFn(fn func(string) error) {
	c.loginUrlHandlerFn = fn
}
