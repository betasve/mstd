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

func (c *Config) ClientId() string {
	return c.clientId
}

func (c *Config) ClientSecret() string {
	return c.clientSecret
}

func (c *Config) ClientPermissions() string {
	return c.permissions
}

func (c *Config) ClientAccessToken() string {
	return c.accessToken
}

func (c *Config) ClientRefreshToken() string {
	return c.refreshToken
}

func (c *Config) ClientAccessTokenExpiresAt() t.Time {
	return c.accessTokenExpiresAt
}

func (c *Config) ClientRefreshTokenExpiresAt() t.Time {
	return c.refreshTokenExpiresAt
}

func (c *Config) AuthCallbackHost() string {
	return c.authCallbackHost
}

func (c *Config) AuthCallbackPath() string {
	return c.authCallbackPath
}

func clientId() string {
	return viper.Client.GetString(defaultClientIdConfig)
}

func clientSecret() string {
	return viper.Client.GetString(defaultClientSecretConfig)
}

func clientPermissions() string {
	return viper.Client.GetString(defaultPermissionsConfig)
}

func clientAccessToken() string {
	return viper.Client.GetString(defaultAccessTokenConfig)
}

func clientRefreshToken() string {
	return viper.Client.GetString(defaultRefreshTokenConfig)
}

func clientAccessTokenExpiry() t.Time {
	expires := viper.Client.GetInt64(defaultAccessTokenExpiryConfig)

	return t.Unix(expires, 0)
}

func clientRefreshTokenExpiry() t.Time {
	expires := viper.Client.GetInt64(defaultRefreshTokenExpiryConfig)

	return t.Unix(expires, 0)
}

func authCallbackHost() string {
	return viper.Client.GetString(defaultAuthCallbackHost)
}

func authCallbackPath() string {
	return viper.Client.GetString(defaultAuthCallbackPath)
}

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

func secondsToDuration(s int) (t.Duration, error) {
	secsStr := strconv.Itoa(s)
	durSecs, err := tm.Client.ParseDuration(secsStr + "s")

	return durSecs, err
}

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

func setEnvVariables() {
	viper.Client.AutomaticEnv()
}

func readConfigFile() error {
	if err := viper.Client.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

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

func validateClientIdConfigPresence() error {
	if len(clientId()) == 0 {
		return fmt.Errorf("Missing %s in config file", defaultClientIdConfig)
	}

	return nil
}

func validateClientSecretConfigPresence() error {
	if len(clientSecret()) == 0 {
		return fmt.Errorf("Missing %s in config file", defaultClientSecretConfig)
	}

	return nil
}

func validateClientPermissionsConfigPresence() error {
	if len(clientPermissions()) == 0 {
		return fmt.Errorf("Missing %s in config file", defaultPermissionsConfig)
	}

	return nil
}

func validateAuthCallbackHostConfigPresence() error {
	if len(authCallbackHost()) == 0 {
		return fmt.Errorf("Missing %s in config file", defaultAuthCallbackHost)
	}

	return nil
}

func validateAuthCallbackPathConfigPresence() error {
	if len(authCallbackPath()) == 0 {
		return fmt.Errorf("Missing %s in config file", defaultAuthCallbackPath)
	}

	return nil
}

func homeDir() (string, error) {
	home, err := homedir.Client.Dir()
	if err != nil {
		return "", err
	}

	return home, nil
}

func unixTimeAfter(d t.Duration) int64 {
	return tm.Client.Now().Add(d).UnixNano() / nanosecondsInASecond
}
