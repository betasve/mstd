package conf

import (
	"errors"
	"github.com/betasve/mstd/homedir"
	log "github.com/betasve/mstd/logger"
	"github.com/betasve/mstd/viper"
	"strings"
	"testing"
)

func TestGetClientSecret(t *testing.T) {
	getString = defaultClientSecretConfig

	viper.Client = ViperServiceMock{}

	result := GetClientSecret()
	if result != getString {
		t.Errorf("expected \n%s \n but got\n%s", getString, result)
	}
}

func TestSetViperConfigWithConfigFilePath(t *testing.T) {
	var logResult string
	var result string
	CfgFilePath = "file/path"
	configFileUsed = "conf.yml"

	viper.Client = ViperServiceMock{}
	log.Client = LoggerServiceMock{}

	logMock = func(in string) { logResult = in }
	setCfgFilePathFunc = func(in string) { result = CfgFilePath }

	setViperConfig()
	if result != CfgFilePath {
		t.Errorf("expected \n%s \n but got\n%s", CfgFilePath, result)
	}

	if !strings.Contains(logResult, "Using config file: "+configFileUsed) {
		t.Errorf("expected \n%s \n to contain\n%s", logResult, configFileUsed)
	}
}

func TestSetViperConfigWithoutConfigFilePath(t *testing.T) {
	var logResult string
	var addConfigPathResult string
	var setConfigNameResult string
	CfgFilePath = ""
	configFileUsed = "conf.yml"

	viper.Client = ViperServiceMock{}
	log.Client = LoggerServiceMock{}

	logMock = func(in string) { logResult = in }
	addConfigPathFunc = func(in string) { addConfigPathResult = "homedir" }
	setConfigNameFunc = func(in string) { setConfigNameResult = configFileUsed }

	setViperConfig()
	if addConfigPathResult != "homedir" && setConfigNameResult != configFileUsed {
		t.Errorf(
			"expected config path \n%s \n but was\n%s"+
				"\nexpected config name \n%s but was\n%s",
			addConfigPathResult,
			"homedir",
			setConfigNameResult,
			configFileUsed,
		)
	}

	if !strings.Contains(logResult, "Using config file: "+configFileUsed) {
		t.Errorf("expected \n%s \n to contain\n%s", configFileUsed, logResult)
	}
}

func TestSetEnvironmentVariables(t *testing.T) {
	viper.Client = ViperServiceMock{}
	var funcInvoked bool
	automaticEnvFunc = func() { funcInvoked = true }

	setEnvVariables()
	if !funcInvoked {
		t.Error("expected to invoke function but did not")
	}
}

func TestReadConfigFileSuccess(t *testing.T) {
	var result string
	configFileUsed = "Using config file: .mstd.yml"
	logMock = func(in string) { result = in }
	log.Client = LoggerServiceMock{}
	viper.Client = ViperServiceMock{}

	readConfigFile()

	if !strings.Contains(result, configFileUsed) {
		t.Errorf("expected \n%s \n to contain\n%s", configFileUsed, result)
	}
}

func TestReadConfigFileFailure(t *testing.T) {
	var result error
	err := errors.New("Cannot read in config")
	configErr = err
	fatalMock = func(in ...interface{}) { result = err }
	log.Client = LoggerServiceMock{}

	viper.Client = ViperServiceMock{}
	readConfigFile()

	if result != err {
		t.Errorf("expected error\n%s \n but got\n%s", result.Error(), configErr.Error())
	}
}

func TestValidateConfigFileAttributesSuccess(t *testing.T) {
	viper.Client = ViperServiceMock{}
	log.Client = LoggerServiceMock{}
	var result error = nil
	fatalMock = func(in ...interface{}) { result = errors.New("invalid call") }
	getString = "clientId"

	validateConfigFileAttributes()
	if result != nil {
		t.Errorf("expected no errors\n \n but got\n%s", result.Error())
	}
}

func TestValidateConfigFileAttributesClientIdFailure(t *testing.T) {
	viper.Client = ViperServiceMock{}
	log.Client = LoggerServiceMock{}
	var result error = nil
	fatalfMock = func(f string, in ...interface{}) { result = errors.New("missing client id") }
	getStringFunc = func(key string) string {
		if key == defaultClientIdConfig {
			return ""
		} else {
			return "clientSecret"
		}
	}

	validateConfigFileAttributes()
	if result == nil {
		t.Errorf("expected %s\n \n but got\nnil", result.Error())
	}
}

func TestValidateConfigFileAttributesClientSecretFailure(t *testing.T) {
	viper.Client = ViperServiceMock{}
	log.Client = LoggerServiceMock{}
	var result error = nil
	fatalfMock = func(f string, in ...interface{}) { result = errors.New("missing client secret") }
	getStringFunc = func(key string) string {
		if key == defaultClientSecretConfig {
			return "clientId"
		} else {
			return ""
		}
	}

	validateConfigFileAttributes()
	if result == nil {
		t.Errorf("expected %s\n \n but got\nnil", result.Error())
	}
}

func TestValidateClientIdConfigPresenceSuccess(t *testing.T) {
	viper.Client = ViperServiceMock{}
	log.Client = LoggerServiceMock{}
	var result error = nil
	fatalMock = func(in ...interface{}) { result = errors.New("invalid call") }
	getString = "clientId"

	validateClientIdConfigPresence()
	if result != nil {
		t.Errorf("expected no errors\n \n but got\n%s", result.Error())
	}
}

func TestValidateClientIdConfigPresenceFailure(t *testing.T) {
	viper.Client = ViperServiceMock{}
	log.Client = LoggerServiceMock{}
	var result error = nil

	fatalfMock = func(f string, in ...interface{}) {
		result = errors.New("invalid call")
	}

	getString = ""

	validateClientIdConfigPresence()
	if result == nil {
		t.Errorf("expected \n%s \n but got\nnil", result.Error())
	}
}

func TestValidateClientSecretConfigPresenceSuccess(t *testing.T) {
	viper.Client = ViperServiceMock{}
	log.Client = LoggerServiceMock{}
	var result error = nil
	fatalMock = func(in ...interface{}) { result = errors.New("invalid call") }
	getString = "clientSecret"

	validateClientSecretConfigPresence()
	if result != nil {
		t.Errorf("expected no errors\n \n but got\n%s", result.Error())
	}
}

func TestValidateClientSecretConfigPresenceFailure(t *testing.T) {
	viper.Client = ViperServiceMock{}
	log.Client = LoggerServiceMock{}
	var result error = nil

	fatalfMock = func(format string, in ...interface{}) {
		result = errors.New("invalid call")
	}

	getString = ""

	validateClientSecretConfigPresence()
	if result != nil {
		t.Errorf("expected \n%s \n but got\nnil", result.Error())
	}
}

func TestHomedirSuccess(t *testing.T) {
	homedir.Client = HomedirServiceMock{}

	result := homeDir()
	if result != homedirPath {
		t.Errorf("expected \n%s \n but got\n%s", result, homedirPath)
	}
}

func TestHomedirFailure(t *testing.T) {
	homedir.Client = HomedirServiceMock{}

	var resultErr error
	homedirErr = errors.New("No homedir")
	fatalMock = func(in ...interface{}) { resultErr = homedirErr }
	log.Client = LoggerServiceMock{}

	homeDir()
	if resultErr != homedirErr {
		t.Errorf("expected \n%s \n but got\n%s", resultErr.Error(), homedirErr.Error())
	}
}
