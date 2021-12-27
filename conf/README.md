# conf
--
    import "."

The `conf` package is in charge of handling the configuration of the app. It
'knows' only about what string values it needs from a config file and provides
an API (for the rest of the packages) on how to read the valus (tokens) they
need.

## Usage

#### type Config

```go
type Config struct {
}
```


#### func (*Config) AuthCallbackHost

```go
func (c *Config) AuthCallbackHost() string
```
A getter function for the authCallbackHost.

#### func (*Config) AuthCallbackPath

```go
func (c *Config) AuthCallbackPath() string
```
A getter function for the authCallbackPath.

#### func (*Config) ClientAccessToken

```go
func (c *Config) ClientAccessToken() string
```
A getter function for the accessToken.

#### func (*Config) ClientAccessTokenExpiresAt

```go
func (c *Config) ClientAccessTokenExpiresAt() t.Time
```
A getter function for the accessTokenExpiresAt.

#### func (*Config) ClientId

```go
func (c *Config) ClientId() string
```
A getter function for the clientId.

#### func (*Config) ClientPermissions

```go
func (c *Config) ClientPermissions() string
```
A getter function for the permissions.

#### func (*Config) ClientRefreshToken

```go
func (c *Config) ClientRefreshToken() string
```
A getter function for the refreshToken.

#### func (*Config) ClientRefreshTokenExpiresAt

```go
func (c *Config) ClientRefreshTokenExpiresAt() t.Time
```
A getter function for the refreshTokenExpiresAt.

#### func (*Config) ClientSecret

```go
func (c *Config) ClientSecret() string
```
A getter function for the clientSecret.

#### func (*Config) InitConfig

```go
func (c *Config) InitConfig(cfgFilePath string) error
```
Initializes the Config struct, holding most of the configuration related
information that the app is currently needing. For doing so it relies only on
reading the information (and formatting it) from the config file we've specified
upon running the app (or using the default one).

#### func (*Config) SetClientAccessToken

```go
func (c *Config) SetClientAccessToken(in string) error
```
A setter method for the accessToken.

#### func (*Config) SetClientAccessTokenExpirySeconds

```go
func (c *Config) SetClientAccessTokenExpirySeconds(seconds int) error
```
A setter method for accessTokenExpiresAt (in seconds).

#### func (*Config) SetClientRefreshToken

```go
func (c *Config) SetClientRefreshToken(in string) error
```
A setter method for refreshToken.

#### func (*Config) SetClientRefreshTokenExpirySeconds

```go
func (c *Config) SetClientRefreshTokenExpirySeconds(seconds int) error
```
A setter method for refreshTokenExpiresAt (in seconds).
