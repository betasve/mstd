package conf

import (
	"errors"
	"github.com/betasve/mstd/homedir"
	log "github.com/betasve/mstd/logger"
	tm "github.com/betasve/mstd/time"
	"github.com/betasve/mstd/viper"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestInitConfig(t *testing.T) {
	CfgFilePath = "conf.yml"
	viper.Client = ViperServiceMock{}
	CfgFilePath = "file/path"
	configFileUsed = "file/path"
	automaticEnvFunc = func() {}
	setCfgFilePathFunc = func(in string) {}
	getString = "clientId"

	log.Client = LoggerServiceMock{}
	var logResult string
	logMock = func(in string) { logResult = in }

	InitConfig()

	if !strings.Contains(logResult, "Using config file: file/path") {
		t.Errorf(
			"expected log to contain \n%s \n but got\n%s",
			"Using config file: file/path",
			logResult,
		)
	}
}

func TestGetClientId(t *testing.T) {
	testAccessorMethodFor(GetClientId, defaultClientIdConfig, t)
}

func TestGetClientSecret(t *testing.T) {
	testAccessorMethodFor(GetClientSecret, defaultClientSecretConfig, t)
}

func TestGetClientPermissions(t *testing.T) {
	testAccessorMethodFor(GetClientPermissions, defaultPermissionsConfig, t)
}

func TestGetClientAccessToken(t *testing.T) {
	testAccessorMethodFor(GetClientAccessToken, defaultAccessTokenConfig, t)
}

func TestGetClientRefreshToken(t *testing.T) {
	testAccessorMethodFor(GetClientRefreshToken, defaultRefreshTokenConfig, t)
}

func TestGetClientAccessTokenExpirySuccess(t *testing.T) {
	getString = "1607534638"
	result := GetClientAccessTokenExpiry()
	resultType := reflect.TypeOf(result).Name()

	if resultType != "Time" {
		t.Errorf("expected Time\nbut got\n%s", resultType)
	}
}

func TestGetClientAccessTokenExpiryFailureWithEmptyExpiry(t *testing.T) {
	getString = ""
	var logResult error
	err := errors.New("empty access token expiry")
	log.Client = LoggerServiceMock{}
	fatalMock = func(in ...interface{}) { logResult = err }
	GetClientAccessTokenExpiry()

	if err != logResult {
		t.Errorf("expected error %s\nbut got\nnil", err.Error())
	}
}

func TestGetClientAccessTokenExpiryFailureWithInvalidExpiry(t *testing.T) {
	getString = "abc"
	var logResult error
	err := errors.New("invalid access token expiry")
	log.Client = LoggerServiceMock{}
	fatalMock = func(in ...interface{}) { logResult = err }
	GetClientAccessTokenExpiry()

	if err != logResult {
		t.Errorf("expected error %s\nbut got\nnil", err.Error())
	}
}

func TestGetClientRefreshTokenExpirySuccess(t *testing.T) {
	getString = "1607519079"
	result := GetClientRefreshTokenExpiry()
	resultType := reflect.TypeOf(result).Name()

	if resultType != "Time" {
		t.Errorf("expected Time\nbut got\n%s", resultType)
	}
}

func TestGetClientRefreshTokenExpiryFailureWithEmptyExpiry(t *testing.T) {
	getString = ""
	var logResult error
	err := errors.New("empty refresh token")
	log.Client = LoggerServiceMock{}
	fatalMock = func(in ...interface{}) { logResult = err }
	GetClientRefreshTokenExpiry()

	if err != logResult {
		t.Errorf("expected error %s\nbut got\nnil", err.Error())
	}
}

func TestGetClientRefreshTokenExpiryFailureWithInvalidExpiry(t *testing.T) {
	getString = "abc"
	var logResult error
	err := errors.New("invalid refresh token expiry")
	log.Client = LoggerServiceMock{}
	fatalMock = func(in ...interface{}) { logResult = err }
	GetClientRefreshTokenExpiry()

	if err != logResult {
		t.Errorf("expected error %s\nbut got\nnil", err.Error())
	}
}

func TestSetClientAccessTokenSuccess(t *testing.T) {
	var key string
	var value interface{}
	var writeCalled bool
	expectedVal := "testVal"
	setKeyValue = func(k string, v interface{}) {
		key = k
		value = v
	}
	writeConfigFunc = func() error { writeCalled = true; return nil }

	SetClientAccessToken(expectedVal)

	if key != defaultAccessTokenConfig || value != expectedVal || !writeCalled {
		t.Errorf(
			"expected\nkey: %s,\nvalue: %s,\nwrite called: true\n"+
				"but got\nkey: %s,\nvalue: %s,\nwrite called: %t",
			defaultAccessTokenConfig,
			expectedVal,
			key,
			value,
			writeCalled,
		)
	}
}

func TestSetClientAccessTokenFailure(t *testing.T) {
	err := errors.New("could not write config")
	setKeyValue = func(k string, v interface{}) {}
	writeConfigFunc = func() error { return err }
	log.Client = LoggerServiceMock{}
	var expectedErr error
	fatalMock = func(in ...interface{}) { expectedErr = err }

	SetClientAccessToken("")

	if expectedErr != err {
		t.Errorf("expected error\n %s but got %s", expectedErr.Error(), err.Error())
	}
}

func TestSetClientRefreshTokenSuccess(t *testing.T) {
	var key string
	var value interface{}
	var writeCalled bool
	expectedVal := "testVal"
	setKeyValue = func(k string, v interface{}) {
		key = k
		value = v
	}
	writeConfigFunc = func() error { writeCalled = true; return nil }

	SetClientRefreshToken(expectedVal)

	if key != defaultRefreshTokenConfig || value != expectedVal || !writeCalled {
		t.Errorf(
			"expected\nkey: %s,\nvalue: %s,\nwrite called: true\n"+
				"but got\nkey: %s,\nvalue: %s,\nwrite called: %t",
			defaultRefreshTokenConfig,
			expectedVal,
			key,
			value,
			writeCalled,
		)
	}
}

func TestSetClientRefreshTokenFailure(t *testing.T) {
	err := errors.New("could not write config")
	setKeyValue = func(k string, v interface{}) {}
	writeConfigFunc = func() error { return err }
	log.Client = LoggerServiceMock{}
	var expectedErr error
	fatalMock = func(in ...interface{}) { expectedErr = err }

	SetClientRefreshToken("")

	if expectedErr != err {
		t.Errorf("expected error\n %s but got %s", expectedErr.Error(), err.Error())
	}
}

func TestSetCleintAccessTokenExpirySuccess(t *testing.T) {
	tm.Client = TimeMock{}
	// The unix timestamp equivalent of the testTime + 5 seconds
	var expectedResult int64 = 1607591728
	testTimestamp := "2020-12-10T09:15:23Z"
	testTime, err := time.Parse(time.RFC3339, testTimestamp)

	if err != nil {
		t.Error(err)
	}

	timeNowMockFunc = func() time.Time { return testTime }
	timeParseDurationMockFunc = func(s string) (time.Duration, error) { return time.ParseDuration(s) }

	var resultKey string
	var resultValue int64
	var wroteConfigFlag bool
	setKeyValue = func(k string, v interface{}) {
		resultKey = k
		resultValue = v.(int64)
	}

	writeConfigFunc = func() error { wroteConfigFlag = true; return nil }

	SetClientAccessTokenExpirySeconds(5)
	if defaultAccessTokenExpiryConfig != resultKey {
		t.Errorf("expected \n%s \nbut got\n%s", defaultAccessTokenExpiryConfig, resultKey)
	}

	if expectedResult != resultValue {
		t.Errorf("\nexpected \n%d \nbut got\n%d", expectedResult, resultValue)
	}

	if !wroteConfigFlag {
		t.Error("expected to write to config\nbut did not")
	}
}

func TestSecondsToDurationSuccess(t *testing.T) {
	tm.Client = TimeMock{}
	dur := time.Since(time.Now())
	timeParseDurationMockFunc = func(s string) (time.Duration, error) { return dur, nil }

	result := secondsToDuration(42)
	if result != dur {
		t.Errorf("expected \n%s \n but got\n%s", dur.String(), result.String())
	}
}

func TestSecondsToDurationFailure(t *testing.T) {
	tm.Client = TimeMock{}
	err := errors.New("cannot parse")
	timeParseDurationMockFunc = func(s string) (time.Duration, error) {
		return time.Since(time.Now()), err
	}

	log.Client = LoggerServiceMock{}
	var result error
	fatalMock = func(s ...interface{}) { result = err }

	secondsToDuration(42)
	if result != err {
		t.Errorf("expected \n%s \n but got\n%s", err, result)
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

func TestSetEnvVariables(t *testing.T) {
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
	err := errors.New("missing client id")
	fatalfMock = func(f string, in ...interface{}) { result = err }
	getStringFunc = func(key string) string {
		if key == defaultClientIdConfig {
			return ""
		} else {
			return "clientSecret"
		}
	}

	validateConfigFileAttributes()
	if result == nil {
		t.Errorf("expected %s\n \n but got\nnil", err.Error())
	}
}

func TestValidateConfigFileAttributesClientSecretFailure(t *testing.T) {
	viper.Client = ViperServiceMock{}
	log.Client = LoggerServiceMock{}
	var result error = nil
	err := errors.New("missing client secret")
	fatalfMock = func(f string, in ...interface{}) { result = err }
	getStringFunc = func(key string) string {
		if key == defaultClientSecretConfig {
			return "clientId"
		} else {
			return ""
		}
	}

	validateConfigFileAttributes()
	if result == nil {
		t.Errorf("expected %s\n \n but got\nnil", err.Error())
	}
}

func TestValidateConfigFileAttributesClientPermissionsFailure(t *testing.T) {
	viper.Client = ViperServiceMock{}
	log.Client = LoggerServiceMock{}
	var result error = nil
	err := errors.New("missing client permissions")
	fatalfMock = func(f string, in ...interface{}) { result = err }
	getStringFunc = func(key string) string {
		if key == defaultPermissionsConfig {
			return ""
		} else {
			return "clientPermissions"
		}
	}

	validateConfigFileAttributes()
	if result == nil {
		t.Errorf("expected %s\n \n but got\nnil", err.Error())
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
	err := errors.New("invalid call")

	fatalfMock = func(f string, in ...interface{}) {
		result = err
	}

	getString = ""
	getStringFunc = nil

	validateClientIdConfigPresence()
	if result == nil {
		t.Errorf("expected \n%s \n but got\nnil", err.Error())
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
	err := errors.New("invalid call")

	fatalfMock = func(format string, in ...interface{}) {
		result = err
	}

	getString = ""
	getStringFunc = nil

	validateClientSecretConfigPresence()
	if result == nil {
		t.Errorf("expected \n%s \n but got\nnil", err.Error())
	}
}

func TestValidateClientConfigPermissionsPresenceSuccess(t *testing.T) {
	viper.Client = ViperServiceMock{}
	log.Client = LoggerServiceMock{}
	var result error = nil
	fatalMock = func(in ...interface{}) { result = errors.New("invalid call") }
	getString = "clientPermissions"

	validateClientPermissionsConfigPresence()
	if result != nil {
		t.Errorf("expected no errors\n \n but got\n%s", result.Error())
	}
}

func TestValidateClientPermissionsConfigPresenceFailure(t *testing.T) {
	viper.Client = ViperServiceMock{}
	log.Client = LoggerServiceMock{}
	var result error = nil
	err := errors.New("invalid call")

	fatalfMock = func(format string, in ...interface{}) {
		result = err
	}

	getString = ""
	getStringFunc = nil

	validateClientPermissionsConfigPresence()
	if result == nil {
		t.Errorf("expected \n%s \n but got\nnil", err.Error())
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

func testAccessorMethodFor(method func() string, stubValue string, t *testing.T) {
	getString = stubValue

	viper.Client = ViperServiceMock{}

	result := method()
	if result != getString {
		t.Errorf("expected \n%s \n but got\n%s", getString, result)
	}
}

func TestUnixTimeAfter(t *testing.T) {
	testTime, err := time.Parse(time.RFC3339, "2020-12-10T09:15:23Z")

	if err != nil {
		t.Error(err)
	}

	// The unix timestamp equivalent of the testTime + 5 seconds
	var expectedResult int64 = 1607591728

	fiveSecodsDuration, err := time.ParseDuration("5s")
	if err != nil {
		t.Error(err)
	}

	tm.Client = TimeMock{}
	timeNowMockFunc = func() time.Time { return testTime }

	result := unixTimeAfter(fiveSecodsDuration)

	if result != expectedResult {
		t.Errorf("expected \n%d \n but got\n%d", expectedResult, result)
	}
}
