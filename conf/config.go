package conf

import (
	"fmt"
	"github.com/betasve/mstd/homedir"
	log "github.com/betasve/mstd/logger"
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
const defaultAccessTokenExpiryConfig string = "ate"
const defaultRefreshTokenExpiryConfig string = "rte"
const nanosecondsInASecond int64 = 1_000_000_000

var CfgFilePath string

func InitConfig() {
	setEnvVariables()
	setViperConfig()

	validateConfigFileAttributes()
}

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
	expires := viper.Client.GetString(defaultAccessTokenExpiryConfig)

	if len(expires) == 0 {
		log.Client.Fatal("Token expiration time is empty.Please login again")
	}

	i, err := strconv.ParseInt(expires, 10, 64)
	if err != nil {
		log.Client.Fatal("Problem retrieving token expiration time. Please login again")
	}

	return t.Unix(i, 0)
}

func GetClientRefreshTokenExpiry() t.Time {
	expires := viper.Client.GetString(defaultRefreshTokenExpiryConfig)

	if len(expires) == 0 {
		log.Client.Fatal("Refresh token expiration time is empty.Please login again")
	}

	i, err := strconv.ParseInt(expires, 10, 64)
	if err != nil {
		log.Client.Fatal("Problem retrieving refresh token expiration time. Please login again")
	}

	return t.Unix(i, 0)
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

// TODO: Write tests for that method
func SetClientAccessTokenExpiry(seconds int) {
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
		log.Client.Log(fmt.Sprintf("Using config file: %s", viper.Client.ConfigFileUsed()))
	} else {
		log.Client.Fatal(err)
	}
}

func validateConfigFileAttributes() {
	validateClientIdConfigPresence()
	validateClientSecretConfigPresence()
	validateClientPermissionsConfigPresence()
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

func homeDir() string {
	home, err := homedir.Client.Dir()
	if err != nil {
		log.Client.Fatal(err)
	}

	return home
}

// TODO: Cover with tests
func unixTimeAfter(d t.Duration) int64 {
	return t.Now().Add(d).UnixNano() / nanosecondsInASecond
}
