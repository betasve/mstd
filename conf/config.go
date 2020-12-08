package conf

import (
	"fmt"
	"github.com/betasve/mstd/homedir"
	l "github.com/betasve/mstd/logger"
	"github.com/betasve/mstd/viper"
	"log"
)

const defaultConfigFileName string = ".mstd"
const defaultClientIdConfig string = "client_id"
const defaultClientSecretConfig string = "client_secret"

// const defaultPermissionsConfig string = "permissions"
// const defaultAccessTokenConfig string = "access_token"
// const defaultRefreshTokenConfig string = "refresh_token"

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
		l.Client.Log(fmt.Sprintf("Using config file: %s", viper.Client.ConfigFileUsed()))
	} else {
		l.Client.Fatal(err)
	}
}

func validateConfigFileAttributes() {
	validateClientIdConfigPresence()
	validateClientSecretConfigPresence()
}

func validateClientIdConfigPresence() {
	if len(GetClientId()) == 0 {
		log.Fatalf("Missing %s in config file", defaultClientIdConfig)
	}
}

func validateClientSecretConfigPresence() {
	if len(GetClientSecret()) == 0 {
		l.Client.Fatalf("Missing %s in config file", defaultClientSecretConfig)
	}
}

func homeDir() string {
	home, err := homedir.Client.Dir()
	if err != nil {
		l.Client.Fatal(err)
	}

	return home
}
