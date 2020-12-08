package conf

import (
	"errors"
	"github.com/betasve/mstd/homedir"
	l "github.com/betasve/mstd/logger"
	"github.com/betasve/mstd/viper"
	"strings"
	"testing"
)

func TestReadConfigFileSuccess(t *testing.T) {
	var result string
	configFileUsed = "Using config file: .mstd.yml"
	logMock = func(in string) { result = in }
	l.Client = LoggerServiceMock{}
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
	l.Client = LoggerServiceMock{}

	viper.Client = ViperServiceMock{}
	readConfigFile()

	if result != err {
		t.Errorf("expected error\n%s \n but got\n%s", result.Error(), configErr.Error())
	}
}

func TestValidateClientSecretConfigPresenceSuccess(t *testing.T) {
	viper.Client = ViperServiceMock{}
	l.Client = LoggerServiceMock{}
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
	l.Client = LoggerServiceMock{}
	var result error = nil

	fatalfMock = func(format string, in ...interface{}) {
		result = errors.New("invalid call")
	}

	getString = ""

	validateClientSecretConfigPresence()
	if result == nil {
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
	l.Client = LoggerServiceMock{}

	homeDir()
	if resultErr != homedirErr {
		t.Errorf("expected \n%s \n but got\n%s", resultErr.Error(), homedirErr.Error())
	}
}
