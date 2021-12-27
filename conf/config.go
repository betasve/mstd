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

// The `conf` package is in charge of handling the configuration of the app.
// It 'knows' only about what string values it needs from a config file and
// provides an API (for the rest of the packages) on how to read the valus (tokens)
// they need.
package conf

import (
	"errors"
	"fmt"
	"github.com/betasve/mstd/ext/homedir"
	tm "github.com/betasve/mstd/ext/time"
	"github.com/betasve/mstd/ext/viper"
	"strconv"
	"strings"
	"sync"
	t "time"
)

const defaultConfigFileName string = ".mstd"
const defaultClientIdConfig string = "client_id"
const defaultClientSecretConfig string = "client_secret"
const defaultPermissionsConfig string = "permissions"
const defaultAccessTokenConfig string = "access_token"
const defaultRefreshTokenConfig string = "refresh_token"
const defaultAuthCallbackHost string = "auth_callback_host_and_port"
const defaultAuthCallbackPath string = "auth_callback_path"
const defaultAccessTokenExpiryConfig string = "ate"
const defaultRefreshTokenExpiryConfig string = "rte"
const nanosecondsInASecond int64 = 1_000_000_000

type Config struct {
	mu                    sync.Mutex
	clientId              string
	clientSecret          string
	permissions           string
	accessToken           string
	refreshToken          string
	accessTokenExpiresAt  t.Time
	refreshTokenExpiresAt t.Time
	authCallbackHost      string
	authCallbackPath      string
}

// Initializes the Config struct, holding most of the configuration related
// information that the app is currently needing. For doing so it relies only
// on reading the information (and formatting it) from the config file we've
// specified upon running the app (or using the default one).
func (c *Config) InitConfig(cfgFilePath string) error {
	setEnvVariables()

	if err := setViperConfig(cfgFilePath); err != nil {
		return err
	}
	if err := validateConfigFileAttributes(); err != nil {
		return err
	}

	c.populateConfigValues()

	return nil
}

// A getter function for the clientId.
func (c *Config) ClientId() string {
	return c.clientId
}

// A getter function for the clientSecret.
func (c *Config) ClientSecret() string {
	return c.clientSecret
}

// A getter function for the permissions.
func (c *Config) ClientPermissions() string {
	return c.permissions
}

// A getter function for the accessToken.
func (c *Config) ClientAccessToken() string {
	return c.accessToken
}

// A getter function for the refreshToken.
func (c *Config) ClientRefreshToken() string {
	return c.refreshToken
}

// A getter function for the accessTokenExpiresAt.
func (c *Config) ClientAccessTokenExpiresAt() t.Time {
	return c.accessTokenExpiresAt
}

// A getter function for the refreshTokenExpiresAt.
func (c *Config) ClientRefreshTokenExpiresAt() t.Time {
	return c.refreshTokenExpiresAt
}

// A getter function for the authCallbackHost.
func (c *Config) AuthCallbackHost() string {
	return c.authCallbackHost
}

// A getter function for the authCallbackPath.
func (c *Config) AuthCallbackPath() string {
	return c.authCallbackPath
}

// A getter function for the clientId key string.
func clientId() string {
	return viper.Client.GetString(defaultClientIdConfig)
}

// A getter function for the clientSecret key string.
func clientSecret() string {
	return viper.Client.GetString(defaultClientSecretConfig)
}

// A getter function for the permissionsConfig key string.
func clientPermissions() string {
	return viper.Client.GetString(defaultPermissionsConfig)
}

// A getter function for the accessTokent key string.
func clientAccessToken() string {
	return viper.Client.GetString(defaultAccessTokenConfig)
}

// A getter function for the refreshTokent key string.
func clientRefreshToken() string {
	return viper.Client.GetString(defaultRefreshTokenConfig)
}

// A getter function for the refreshTokent key string.
func clientAccessTokenExpiry() t.Time {
	expires := viper.Client.GetInt64(defaultAccessTokenExpiryConfig)

	return t.Unix(expires, 0)
}

// A getter function to provide the token expiratoin Unix timestamp.
func clientRefreshTokenExpiry() t.Time {
	expires := viper.Client.GetInt64(defaultRefreshTokenExpiryConfig)

	return t.Unix(expires, 0)
}

// A getter function to provide the authCallbackHost, needed complete
// the MS authentication process (we are using a local host as we are
// authenticating for the machine the command-line is run on).
func authCallbackHost() string {
	return viper.Client.GetString(defaultAuthCallbackHost)
}

// A getter function to provide the authCallbackPath, needed complete
// the MS authentication process (we are using a local host as we are
// authenticating for the machine the command-line is run on).
func authCallbackPath() string {
	return viper.Client.GetString(defaultAuthCallbackPath)
}

// A setter method for the accessToken.
func (c *Config) SetClientAccessToken(in string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	viper.Client.Set(defaultAccessTokenConfig, in)

	if err := viper.Client.WriteConfig(); err != nil {
		return err
	}

	c.accessToken = in

	return nil
}

// A setter method for refreshToken.
func (c *Config) SetClientRefreshToken(in string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	viper.Client.Set(defaultRefreshTokenConfig, in)

	if err := viper.Client.WriteConfig(); err != nil {
		return err
	}

	c.refreshToken = in

	return nil
}

// A setter method for accessTokenExpiresAt (in seconds).
func (c *Config) SetClientAccessTokenExpirySeconds(seconds int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	tokenDuration, durErr := secondsToDuration(seconds)

	if durErr != nil {
		return durErr
	}

	expiresAt := unixTimeAfter(tokenDuration)

	viper.Client.Set(
		defaultAccessTokenExpiryConfig,
		expiresAt,
	)

	if err := viper.Client.WriteConfig(); err != nil {
		return err
	}

	c.accessTokenExpiresAt = t.Unix(expiresAt, 0)

	return nil
}

// A setter method for refreshTokenExpiresAt (in seconds).
func (c *Config) SetClientRefreshTokenExpirySeconds(seconds int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	tokenDuration, durErr := secondsToDuration(seconds)

	if durErr != nil {
		return durErr
	}

	expiresAt := unixTimeAfter(tokenDuration)

	viper.Client.Set(
		defaultRefreshTokenExpiryConfig,
		expiresAt,
	)

	if err := viper.Client.WriteConfig(); err != nil {
		return err
	}

	c.refreshTokenExpiresAt = t.Unix(expiresAt, 0)

	return nil
}

// A method to populate values in the Config object by reading them from the
// config file.
func (c *Config) populateConfigValues() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.clientId = clientId()
	c.clientSecret = clientSecret()
	c.permissions = clientPermissions()
	c.accessToken = clientAccessToken()
	c.refreshToken = clientRefreshToken()
	c.accessTokenExpiresAt = clientAccessTokenExpiry()
	c.refreshTokenExpiresAt = clientRefreshTokenExpiry()
	c.authCallbackHost = authCallbackHost()
	c.authCallbackPath = authCallbackPath()
}

// A function to concert seconds into a time.Duration object
func secondsToDuration(s int) (t.Duration, error) {
	secsStr := strconv.Itoa(s)
	durSecs, err := tm.Client.ParseDuration(secsStr + "s")

	return durSecs, err
}

// A funciton to set the config file path for the Viper tool.
func setViperConfig(cfgFilePath string) error {
	if cfgFilePath != "" {
		viper.Client.SetConfigFile(cfgFilePath)
	} else {
		home, err := homeDir()
		if err != nil {
			return err
		}
		viper.Client.AddConfigPath(home)
		viper.Client.SetConfigName(defaultConfigFileName)
	}

	return readConfigFile()
}

// A function that sets the environments' variables to the Viper client.
func setEnvVariables() {
	viper.Client.AutomaticEnv()
}

// Read the config file (it's a side-effect function that does not return
// anything but populates Viper's internal state.
func readConfigFile() error {
	if err := viper.Client.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

// Validates the presence of the necessary values in our config file.
func validateConfigFileAttributes() error {
	var err [5]error
	err[0] = validateClientIdConfigPresence()
	err[1] = validateClientSecretConfigPresence()
	err[2] = validateClientPermissionsConfigPresence()
	err[3] = validateAuthCallbackHostConfigPresence()
	err[4] = validateAuthCallbackPathConfigPresence()

	str := []string{"Errors in config file:"}
	for _, e := range err {
		if e != nil {
			str = append(str, e.Error())
		}
	}
	if len(str) != 1 {
		return errors.New(strings.Join(str, "\n"))
	}

	return nil
}

// Validates the presence of a Client ID in our config.
func validateClientIdConfigPresence() error {
	if len(clientId()) == 0 {
		return fmt.Errorf("Missing %s in config file", defaultClientIdConfig)
	}

	return nil
}

// Validates the presence of a Client Secret in our config.
func validateClientSecretConfigPresence() error {
	if len(clientSecret()) == 0 {
		return fmt.Errorf("Missing %s in config file", defaultClientSecretConfig)
	}

	return nil
}

// Validates the presence of a Client Permissions in our config.
func validateClientPermissionsConfigPresence() error {
	if len(clientPermissions()) == 0 {
		return fmt.Errorf("Missing %s in config file", defaultPermissionsConfig)
	}

	return nil
}

// Validates the presence of an Auth Callback Host in our config.
func validateAuthCallbackHostConfigPresence() error {
	if len(authCallbackHost()) == 0 {
		return fmt.Errorf("Missing %s in config file", defaultAuthCallbackHost)
	}

	return nil
}

// Validates the presence of an Auth Callback Path in our config.
func validateAuthCallbackPathConfigPresence() error {
	if len(authCallbackPath()) == 0 {
		return fmt.Errorf("Missing %s in config file", defaultAuthCallbackPath)
	}

	return nil
}

// Retrieves the path to user's home directory
func homeDir() (string, error) {
	home, err := homedir.Client.Dir()
	if err != nil {
		return "", err
	}

	return home, nil
}

// Adds a time.Duration and returns the new timestamp in Unix timestamp.
func unixTimeAfter(d t.Duration) int64 {
	return tm.Client.Now().Add(d).UnixNano() / nanosecondsInASecond
}
