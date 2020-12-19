package conf

import (
	"fmt"
	"github.com/betasve/mstd/homedir"
	log "github.com/betasve/mstd/log"
	tm "github.com/betasve/mstd/time"
	"github.com/betasve/mstd/viper"
	"strconv"
	t "time"
)

const defaultConfigFileName string = ".mstd"
const defaultClientIdConfig string = "client_id"
const defaultClientSecretConfig string = "client_secret"
const defaultPermissionsConfig string = "permissions"
const defaultAccessTokenConfig string = "access_token"
const defaultRefreshTokenConfig string = "refresh_token"
const defaultAuthCallbackHost string = "auth_callback_host_and_port"
const defaultAccessTokenExpiryConfig string = "ate"
const defaultRefreshTokenExpiryConfig string = "rte"
const nanosecondsInASecond int64 = 1_000_000_000

var CfgFilePath string

type State struct {
	ClientId              string
	ClientSecret          string
	Permissions           string
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  t.Time
	RefreshTokenExpiresAt t.Time
	AuthCallbackHost      string
}

var CurrentState = State{}

func InitConfig() {
	setEnvVariables()
	setViperConfig()

	validateConfigFileAttributes()
	populateCurrentState()
}

// TODO: make these public methods private if they turn out not needed
func GetClientId() string {
	return viper.Client.GetString(defaultClientIdConfig)
}

func GetClientSecret() string {
	return viper.Client.GetString(defaultClientSecretConfig)
}

func GetClientPermissions() string {
	return viper.Client.GetString(defaultPermissionsConfig)
}

func GetClientAccessToken() string {
	return viper.Client.GetString(defaultAccessTokenConfig)
}

func GetClientRefreshToken() string {
	return viper.Client.GetString(defaultRefreshTokenConfig)
}

func GetClientAccessTokenExpiry() t.Time {
	expires := viper.Client.GetInt64(defaultAccessTokenExpiryConfig)

	return t.Unix(expires, 0)
}

func GetClientRefreshTokenExpiry() t.Time {
	expires := viper.Client.GetInt64(defaultRefreshTokenExpiryConfig)

	return t.Unix(expires, 0)
}

func GetAuthCallbackHost() string {
	return viper.Client.GetString(defaultAuthCallbackHost)
}

func SetClientAccessToken(in string) {
	viper.Client.Set(defaultAccessTokenConfig, in)
	err := viper.Client.WriteConfig()
	if err != nil {
		log.Client.Fatal(err.Error())
	}
}

func SetClientRefreshToken(in string) {
	viper.Client.Set(defaultRefreshTokenConfig, in)
	err := viper.Client.WriteConfig()
	if err != nil {
		log.Client.Fatal(err.Error())
	}
}

func SetClientAccessTokenExpirySeconds(seconds int) {
	tokenDuration := secondsToDuration(seconds)

	viper.Client.Set(
		defaultAccessTokenExpiryConfig,
		unixTimeAfter(tokenDuration),
	)

	err := viper.Client.WriteConfig()
	if err != nil {
		log.Client.Fatal(err.Error())
	}
}

func SetClientRefreshTokenExpirySeconds(seconds int) {
	tokenDuration := secondsToDuration(seconds)

	viper.Client.Set(
		defaultRefreshTokenExpiryConfig,
		unixTimeAfter(tokenDuration),
	)

	err := viper.Client.WriteConfig()
	if err != nil {
		log.Client.Fatal(err.Error())
	}
}

func populateCurrentState() {
	CurrentState.ClientId = GetClientId()
	CurrentState.ClientSecret = GetClientSecret()
	CurrentState.Permissions = GetClientPermissions()
	CurrentState.AccessToken = GetClientAccessToken()
	CurrentState.RefreshToken = GetClientRefreshToken()
	CurrentState.AccessTokenExpiresAt = GetClientAccessTokenExpiry()
	CurrentState.RefreshTokenExpiresAt = GetClientRefreshTokenExpiry()
	CurrentState.AuthCallbackHost = GetAuthCallbackHost()
}

func secondsToDuration(s int) t.Duration {
	secsStr := strconv.Itoa(s)
	durSecs, err := tm.Client.ParseDuration(secsStr + "s")
	if err != nil {
		log.Client.Fatal(err.Error())
	}

	return durSecs
}

func setViperConfig() {
	if CfgFilePath != "" {
		viper.Client.SetConfigFile(CfgFilePath)
	} else {
		viper.Client.AddConfigPath(homeDir())
		viper.Client.SetConfigName(defaultConfigFileName)
	}

	readConfigFile()
}

func setEnvVariables() {
	viper.Client.AutomaticEnv()
}

func readConfigFile() {
	if err := viper.Client.ReadInConfig(); err == nil {
		log.Client.Println(fmt.Sprintf("Using config file: %s", viper.Client.ConfigFileUsed()))
	} else {
		log.Client.Fatal(err)
	}
}

// TODO: Improve method to return accumulated error
// by improving each validation to return error on its own
func validateConfigFileAttributes() {
	validateClientIdConfigPresence()
	validateClientSecretConfigPresence()
	validateClientPermissionsConfigPresence()
	validateAuthCallbackHostConfigPresence()
}

func validateClientIdConfigPresence() {
	if len(GetClientId()) == 0 {
		log.Client.Fatalf("Missing %s in config file", defaultClientIdConfig)
	}
}

func validateClientSecretConfigPresence() {
	if len(GetClientSecret()) == 0 {
		log.Client.Fatalf("Missing %s in config file", defaultClientSecretConfig)
	}
}

func validateClientPermissionsConfigPresence() {
	if len(GetClientPermissions()) == 0 {
		log.Client.Fatalf("Missing %s in config file", defaultPermissionsConfig)
	}
}

func validateAuthCallbackHostConfigPresence() {
	if len(GetAuthCallbackHost()) == 0 {
		log.Client.Fatalf("Missing %s in config file", defaultAuthCallbackHost)
	}
}

func homeDir() string {
	home, err := homedir.Client.Dir()
	if err != nil {
		log.Client.Fatal(err)
	}

	return home
}

func unixTimeAfter(d t.Duration) int64 {
	return tm.Client.Now().Add(d).UnixNano() / nanosecondsInASecond
}
