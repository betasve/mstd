# login
--
    import "."

The `login` package is in charge of holding the login logic needed for working
with MS' ToDo API. It's holding the knowledge of what requests to build and
where to send them in order to retrieve the tokens we need so we can
successfully retrieve and create lists and todo items onward in the app.

## Usage

#### func  CallbackListen

```go
func CallbackListen(callbackUrl string, cb func(string) error) error
```
Spins up a tiny HTTP serer on :8008 to listen for a callback and handle the
passed params.

#### type AuthData

```go
type AuthData struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
```


#### type Creds

```go
type Creds struct {
}
```


#### func (*Creds) LoginNeeded

```go
func (c *Creds) LoginNeeded() bool
```
Checks if the user needs to log in or she is already logged in.

#### func (*Creds) PerformLogin

```go
func (c *Creds) PerformLogin() error
```
Logs in a user. TODO: Add logout command to remove attributes from conf file

#### func (*Creds) SetAccessToken

```go
func (c *Creds) SetAccessToken(accessToken string)
```
A setter method for accessToken.

#### func (*Creds) SetAccessTokenExpiresAt

```go
func (c *Creds) SetAccessTokenExpiresAt(accessTokenExpiresAt time.Time)
```
A setter method for accessTokenExpiresAt.

#### func (*Creds) SetAuthCallbackHost

```go
func (c *Creds) SetAuthCallbackHost(host string)
```
A setter method for authCallbackHost.

#### func (*Creds) SetAuthCallbackPath

```go
func (c *Creds) SetAuthCallbackPath(path string)
```
A setter method for authCallbackPath.

#### func (*Creds) SetId

```go
func (c *Creds) SetId(id string)
```
A setter method for clientId.

#### func (*Creds) SetLoginDataCallbackFn

```go
func (c *Creds) SetLoginDataCallbackFn(fn func(*AuthData) error)
```
A setter method for loginDataCallbackFn.

#### func (*Creds) SetLoginUrlHandlerFn

```go
func (c *Creds) SetLoginUrlHandlerFn(fn func(string) error)
```
A setter method for loginUrlHandlerFn.

#### func (*Creds) SetPermissions

```go
func (c *Creds) SetPermissions(permissions string)
```
A setter method for permissions.

#### func (*Creds) SetRefreshToken

```go
func (c *Creds) SetRefreshToken(refreshToken string)
```
A setter method for refreshToken.

#### func (*Creds) SetRefreshTokenExpiresAt

```go
func (c *Creds) SetRefreshTokenExpiresAt(refreshTokenExpiresAt time.Time)
```
A setter method for refreshTokenExpiresAt.

#### func (*Creds) SetSecret

```go
func (c *Creds) SetSecret(secret string)
```
A setter method for clientSecret.
