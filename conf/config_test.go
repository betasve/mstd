package conf

import (
	"errors"
	"fmt"
	"github.com/betasve/mstd/homedir"
	log "github.com/betasve/mstd/log"
	logtest "github.com/betasve/mstd/log/logtest"
	t "github.com/betasve/mstd/time"
	tt "github.com/betasve/mstd/time/timetest"
	"github.com/betasve/mstd/viper"
	vt "github.com/betasve/mstd/viper/vipertest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestInitConfig(test *testing.T) {
	CfgFilePath = "conf.yml"
	viper.Client = vt.ViperServiceMock{}
	CfgFilePath = "file/path"
	vt.ConfigFileUsed = "file/path"
	vt.AutomaticEnvFunc = func() {}
	vt.SetCfgFilePathFunc = func(in string) {}
	vt.GetString = "clientId"

	log.Client = logtest.LoggerServiceMock{}
	var logResult interface{}
	logtest.PrintlnMock = func(in ...interface{}) { logResult = in[0] }

	InitConfig()

	if !strings.Contains(fmt.Sprintf("%v", logResult), "Using config file: file/path") {
		test.Errorf(
			"expected log to contain \n%s \n but got\n%s",
			"Using config file: file/path",
			logResult,
		)
	}
}

func TestGetClientId(test *testing.T) {
	testAccessorMethodFor(GetClientId, defaultClientIdConfig, test)
}

func TestGetClientSecret(test *testing.T) {
	testAccessorMethodFor(GetClientSecret, defaultClientSecretConfig, test)
}

func TestGetClientPermissions(test *testing.T) {
	testAccessorMethodFor(GetClientPermissions, defaultPermissionsConfig, test)
}

func TestGetClientAccessToken(test *testing.T) {
	testAccessorMethodFor(GetClientAccessToken, defaultAccessTokenConfig, test)
}

func TestGetClientRefreshToken(test *testing.T) {
	testAccessorMethodFor(GetClientRefreshToken, defaultRefreshTokenConfig, test)
}

func TestGetClientAccessTokenExpirySuccess(test *testing.T) {
	vt.GetString = "1607534638"
	result := GetClientAccessTokenExpiry()
	resultType := reflect.TypeOf(result).Name()

	if resultType != "Time" {
		test.Errorf("expected Time\nbut got\n%s", resultType)
	}
}

func TestGetClientAccessTokenExpiryFailureWithEmptyExpiry(test *testing.T) {
	vt.GetString = ""
	expiry := GetClientAccessTokenExpiry()

	if expiry.Unix() != 0 {
		test.Errorf("\nexpected Unix timestamp of 0\nbut got\n%v", expiry.Unix())
	}
}

func TestGetClientAccessTokenExpiryFailureWithInvalidExpiry(test *testing.T) {
	vt.GetString = "abc"
	expiry := GetClientAccessTokenExpiry()

	if expiry.Unix() != 0 {
		test.Errorf("\nexpected Unix timestamp of 0\nbut got\n%v", expiry.Unix())
	}
}

func TestGetClientRefreshTokenExpirySuccess(test *testing.T) {
	vt.GetString = "1607519079"
	result := GetClientRefreshTokenExpiry()
	resultType := reflect.TypeOf(result).Name()

	if resultType != "Time" {
		test.Errorf("expected Time\nbut got\n%s", resultType)
	}
}

func TestGetClientRefreshTokenExpiryFailureWithEmptyExpiry(test *testing.T) {
	vt.GetString = ""
	expiry := GetClientRefreshTokenExpiry()

	if expiry.Unix() != 0 {
		test.Errorf("\nexpected Unix timestamp of 0\nbut got\n%v", expiry.Unix())
	}
}

func TestGetClientRefreshTokenExpiryFailureWithInvalidExpiry(test *testing.T) {
	expiry := GetClientRefreshTokenExpiry()

	if expiry.Unix() != 0 {
		test.Errorf("expected 0 Unix timestamp but got %v", expiry.Unix())
	}
}

func TestSetClientAccessTokenSuccess(test *testing.T) {
	var key string
	var value interface{}
	var writeCalled bool
	expectedVal := "testVal"
	vt.SetKeyValue = func(k string, v interface{}) {
		key = k
		value = v
	}
	vt.WriteConfigFunc = func() error { writeCalled = true; return nil }

	SetClientAccessToken(expectedVal)

	if key != defaultAccessTokenConfig || value != expectedVal || !writeCalled {
		test.Errorf(
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

func TestSetClientAccessTokenFailure(test *testing.T) {
	err := errors.New("could not write config")
	vt.SetKeyValue = func(k string, v interface{}) {}
	vt.WriteConfigFunc = func() error { return err }
	log.Client = logtest.LoggerServiceMock{}
	var expectedErr error
	logtest.FatalMock = func(in ...interface{}) { expectedErr = err }

	SetClientAccessToken("")

	if expectedErr != err {
		test.Errorf("expected error\n %s but got %s", expectedErr.Error(), err.Error())
	}
}

func TestSetClientRefreshTokenSuccess(test *testing.T) {
	var key string
	var value interface{}
	var writeCalled bool
	expectedVal := "testVal"
	vt.SetKeyValue = func(k string, v interface{}) {
		key = k
		value = v
	}
	vt.WriteConfigFunc = func() error { writeCalled = true; return nil }

	SetClientRefreshToken(expectedVal)

	if key != defaultRefreshTokenConfig || value != expectedVal || !writeCalled {
		test.Errorf(
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

func TestSetClientRefreshTokenFailure(test *testing.T) {
	err := errors.New("could not write config")
	vt.SetKeyValue = func(k string, v interface{}) {}
	vt.WriteConfigFunc = func() error { return err }
	log.Client = logtest.LoggerServiceMock{}
	var expectedErr error
	logtest.FatalMock = func(in ...interface{}) { expectedErr = err }

	SetClientRefreshToken("")

	if expectedErr != err {
		test.Errorf("expected error\n %s but got %s", expectedErr.Error(), err.Error())
	}
}

func TestSetClientAccessTokenExpirySuccess(test *testing.T) {
	t.Client = tt.TimeMock{}
	// The unix timestamp equivalent of the testTime + 5 seconds
	var expectedResult int64 = 1607591728
	testTimestamp := "2020-12-10T09:15:23Z"
	testTime, err := time.Parse(time.RFC3339, testTimestamp)

	if err != nil {
		test.Error(err)
	}

	tt.TimeNowMockFunc = func() time.Time { return testTime }
	tt.TimeParseDurationMockFunc = func(s string) (time.Duration, error) { return time.ParseDuration(s) }

	var resultKey string
	var resultValue int64
	var wroteConfigFlag bool
	vt.SetKeyValue = func(k string, v interface{}) {
		resultKey = k
		resultValue = v.(int64)
	}

	vt.WriteConfigFunc = func() error { wroteConfigFlag = true; return nil }

	SetClientAccessTokenExpirySeconds(5)
	if defaultAccessTokenExpiryConfig != resultKey {
		test.Errorf("expected \n%s \nbut got\n%s", defaultAccessTokenExpiryConfig, resultKey)
	}

	if expectedResult != resultValue {
		test.Errorf("\nexpected \n%d \nbut got\n%d", expectedResult, resultValue)
	}

	if !wroteConfigFlag {
		test.Error("expected to write to config\nbut did not")
	}
}

func TestSetClientAccessTokenExpiryFailure(test *testing.T) {
	t.Client = tt.TimeMock{}
	log.Client = logtest.LoggerServiceMock{}

	expectedErr := errors.New("Could not write config")
	var err error
	vt.WriteConfigFunc = func() error { return expectedErr }
	logtest.FatalMock = func(in ...interface{}) { err = expectedErr }
	SetClientAccessTokenExpirySeconds(5)

	if err != expectedErr {
		test.Errorf("expected \n%s \nbut got\n%s", err, expectedErr)
	}
}

func TestSetClientRefreshTokenExpirySecondsSuccess(test *testing.T) {
	t.Client = tt.TimeMock{}
	// The unix timestamp equivalent of the testTime + 5 seconds
	var expectedResult int64 = 1607591728
	testTimestamp := "2020-12-10T09:15:23Z"
	testTime, err := time.Parse(time.RFC3339, testTimestamp)

	if err != nil {
		test.Error(err)
	}

	tt.TimeNowMockFunc = func() time.Time { return testTime }
	tt.TimeParseDurationMockFunc = func(s string) (time.Duration, error) { return time.ParseDuration(s) }

	var resultKey string
	var resultValue int64
	var wroteConfigFlag bool
	vt.SetKeyValue = func(k string, v interface{}) {
		resultKey = k
		resultValue = v.(int64)
	}

	vt.WriteConfigFunc = func() error { wroteConfigFlag = true; return nil }

	SetClientRefreshTokenExpirySeconds(5)
	if defaultRefreshTokenExpiryConfig != resultKey {
		test.Errorf("expected \n%s \nbut got\n%s", defaultRefreshTokenExpiryConfig, resultKey)
	}

	if expectedResult != resultValue {
		test.Errorf("\nexpected \n%d \nbut got\n%d", expectedResult, resultValue)
	}

	if !wroteConfigFlag {
		test.Error("expected to write to config\nbut did not")
	}
}

func TestSetClientRefreshTokenExpirySecondsFailure(test *testing.T) {
	t.Client = tt.TimeMock{}
	log.Client = logtest.LoggerServiceMock{}

	expectedErr := errors.New("Could not write config")
	var err error
	vt.WriteConfigFunc = func() error { return expectedErr }
	logtest.FatalMock = func(in ...interface{}) { err = expectedErr }
	SetClientRefreshTokenExpirySeconds(5)

	if err != expectedErr {
		test.Errorf("expected \n%s \nbut got\n%s", err, expectedErr)
	}
}

func TestPopulateCurrentState(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	vt.GetStringFunc = func(key string) string {
		switch key {
		case defaultClientIdConfig:
			return "client_id"
		case defaultClientSecretConfig:
			return "client_secret"
		case defaultPermissionsConfig:
			return "client_permissions"
		case defaultAccessTokenConfig:
			return "default_access_token"
		case defaultRefreshTokenConfig:
			return "default_refresh_token"
		default:
			return "default"
		}
	}

	vt.GetInt64Func = func(key string) int64 {
		switch key {
		case defaultAccessTokenExpiryConfig:
			return 1607590000
		case defaultRefreshTokenExpiryConfig:
			return 1607590011
		default:
			return 1
		}
	}

	populateCurrentState()

	if CurrentState.ClientId != "client_id" {
		test.Errorf("\nexpcted client_id\nbut got\n%s", CurrentState.ClientId)
	}

	if CurrentState.ClientSecret != "client_secret" {
		test.Errorf("\nexpcted client_secret\nbut got\n%s", CurrentState.ClientSecret)
	}

	if CurrentState.Permissions != "client_permissions" {
		test.Errorf("\nexpcted client_permissions\nbut got\n%s", CurrentState.Permissions)
	}

	if CurrentState.AccessToken != "default_access_token" {
		test.Errorf("\nexpcted default_access_token\nbut got\n%s", CurrentState.AccessToken)
	}

	if CurrentState.RefreshToken != "default_refresh_token" {
		test.Errorf("\nexpcted default_refresh_token\nbut got\n%s", CurrentState.RefreshToken)
	}

	if CurrentState.AccessTokenExpiresAt != time.Unix(1607590000, 0) {
		test.Errorf("\nexpcted %v\nbut got\n%s", time.Unix(1607590000, 0), CurrentState.AccessTokenExpiresAt)
	}

	if CurrentState.RefreshTokenExpiresAt != time.Unix(1607590011, 0) {
		test.Errorf("\nexpcted %v\nbut got\n%s", time.Unix(1607590011, 0), CurrentState.RefreshTokenExpiresAt)
	}
}

func TestSecondsToDurationSuccess(test *testing.T) {
	t.Client = tt.TimeMock{}
	dur := time.Since(time.Now())
	tt.TimeParseDurationMockFunc = func(s string) (time.Duration, error) { return dur, nil }

	result := secondsToDuration(42)
	if result != dur {
		test.Errorf("expected \n%s \n but got\n%s", dur.String(), result.String())
	}
}

func TestSecondsToDurationFailure(test *testing.T) {
	t.Client = tt.TimeMock{}
	err := errors.New("cannot parse")
	tt.TimeParseDurationMockFunc = func(s string) (time.Duration, error) {
		return time.Since(time.Now()), err
	}

	log.Client = logtest.LoggerServiceMock{}
	var result error
	logtest.FatalMock = func(s ...interface{}) { result = err }

	secondsToDuration(42)
	if result != err {
		test.Errorf("expected \n%s \n but got\n%s", err, result)
	}
}

func TestSetViperConfigWithConfigFilePath(test *testing.T) {
	var logResult interface{}
	var result string
	CfgFilePath = "file/path"
	vt.ConfigFileUsed = "conf.yml"

	viper.Client = vt.ViperServiceMock{}
	log.Client = logtest.LoggerServiceMock{}

	logtest.PrintlnMock = func(in ...interface{}) { logResult = in }
	vt.SetCfgFilePathFunc = func(in string) { result = CfgFilePath }

	setViperConfig()
	if result != CfgFilePath {
		test.Errorf("expected \n%s \n but got\n%s", CfgFilePath, result)
	}

	if !strings.Contains(fmt.Sprintf("%v", logResult), "Using config file: "+vt.ConfigFileUsed) {
		test.Errorf("expected \n%s \n to contain\n%s", logResult, vt.ConfigFileUsed)
	}
}

func TestSetViperConfigWithoutConfigFilePath(test *testing.T) {
	var logResult interface{}
	var addConfigPathResult string
	var setConfigNameResult string
	CfgFilePath = ""
	vt.ConfigFileUsed = "conf.yml"

	viper.Client = vt.ViperServiceMock{}
	log.Client = logtest.LoggerServiceMock{}

	logtest.PrintlnMock = func(in ...interface{}) { logResult = in }
	vt.AddConfigPathFunc = func(in string) { addConfigPathResult = "homedir" }
	vt.SetConfigNameFunc = func(in string) { setConfigNameResult = vt.ConfigFileUsed }

	setViperConfig()
	if addConfigPathResult != "homedir" && setConfigNameResult != vt.ConfigFileUsed {
		test.Errorf(
			"expected config path \n%s \n but was\n%s"+
				"\nexpected config name \n%s but was\n%s",
			addConfigPathResult,
			"homedir",
			setConfigNameResult,
			vt.ConfigFileUsed,
		)
	}

	if !strings.Contains(fmt.Sprintf("%v", logResult), "Using config file: "+vt.ConfigFileUsed) {
		test.Errorf("expected \n%s \n to contain\n%s", vt.ConfigFileUsed, logResult)
	}
}

func TestSetEnvVariables(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	var funcInvoked bool
	vt.AutomaticEnvFunc = func() { funcInvoked = true }

	setEnvVariables()
	if !funcInvoked {
		test.Error("expected to invoke function but did not")
	}
}

func TestReadConfigFileSuccess(test *testing.T) {
	var result interface{}
	vt.ConfigFileUsed = "Using config file: .mstd.yml"
	logtest.PrintlnMock = func(in ...interface{}) { result = in }
	log.Client = logtest.LoggerServiceMock{}
	viper.Client = vt.ViperServiceMock{}

	readConfigFile()

	if !strings.Contains(fmt.Sprintf("%v", result), vt.ConfigFileUsed) {
		test.Errorf("expected \n%s \n to contain\n%s", vt.ConfigFileUsed, result)
	}
}

func TestReadConfigFileFailure(test *testing.T) {
	var result error
	err := errors.New("Cannot read in config")
	vt.ConfigErr = err
	logtest.FatalMock = func(in ...interface{}) { result = err }
	log.Client = logtest.LoggerServiceMock{}

	viper.Client = vt.ViperServiceMock{}
	readConfigFile()

	if result != err {
		test.Errorf("expected error\n%s \n but got\n%s", result.Error(), vt.ConfigErr.Error())
	}
}

func TestValidateConfigFileAttributesSuccess(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	log.Client = logtest.LoggerServiceMock{}
	var result error = nil
	logtest.FatalMock = func(in ...interface{}) { result = errors.New("invalid call") }
	vt.GetString = "clientId"

	validateConfigFileAttributes()
	if result != nil {
		test.Errorf("expected no errors\n \n but got\n%s", result.Error())
	}
}

func TestValidateConfigFileAttributesClientIdFailure(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	log.Client = logtest.LoggerServiceMock{}
	var result error = nil
	err := errors.New("missing client id")
	logtest.FatalfMock = func(f string, in ...interface{}) { result = err }
	vt.GetStringFunc = func(key string) string {
		if key == defaultClientIdConfig {
			return ""
		} else {
			return "clientSecret"
		}
	}

	validateConfigFileAttributes()
	if result == nil {
		test.Errorf("\nexpected\n%s\nbut got\nnil", err.Error())
	}
}

func TestValidateConfigFileAttributesClientSecretFailure(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	log.Client = logtest.LoggerServiceMock{}
	var result error = nil
	err := errors.New("missing client secret")
	logtest.FatalfMock = func(f string, in ...interface{}) { result = err }
	vt.GetStringFunc = func(key string) string {
		if key == defaultClientSecretConfig {
			return "clientId"
		} else {
			return ""
		}
	}

	validateConfigFileAttributes()
	if result == nil {
		test.Errorf("expected %s\n \n but got\nnil", err.Error())
	}
}

func TestValidateConfigFileAttributesClientPermissionsFailure(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	log.Client = logtest.LoggerServiceMock{}
	var result error = nil
	err := errors.New("missing client permissions")
	logtest.FatalfMock = func(f string, in ...interface{}) { result = err }
	vt.GetStringFunc = func(key string) string {
		if key == defaultPermissionsConfig {
			return ""
		} else {
			return "clientPermissions"
		}
	}

	validateConfigFileAttributes()
	if result == nil {
		test.Errorf("expected %s\n \n but got\nnil", err.Error())
	}
}

func TestValidateClientIdConfigPresenceSuccess(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	log.Client = logtest.LoggerServiceMock{}
	var result error = nil
	logtest.FatalMock = func(in ...interface{}) { result = errors.New("invalid call") }
	vt.GetString = "clientId"

	validateClientIdConfigPresence()
	if result != nil {
		test.Errorf("expected no errors\n \n but got\n%s", result.Error())
	}
}

func TestValidateClientIdConfigPresenceFailure(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	log.Client = logtest.LoggerServiceMock{}
	var result error = nil
	err := errors.New("invalid call")

	logtest.FatalfMock = func(f string, in ...interface{}) {
		result = err
	}

	vt.GetString = ""
	vt.GetStringFunc = nil

	validateClientIdConfigPresence()
	if result == nil {
		test.Errorf("expected \n%s \n but got\nnil", err.Error())
	}
}

func TestValidateClientSecretConfigPresenceSuccess(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	log.Client = logtest.LoggerServiceMock{}
	var result error = nil
	logtest.FatalMock = func(in ...interface{}) { result = errors.New("invalid call") }
	vt.GetString = "clientSecret"

	validateClientSecretConfigPresence()
	if result != nil {
		test.Errorf("expected no errors\n \n but got\n%s", result.Error())
	}
}

func TestValidateClientSecretConfigPresenceFailure(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	log.Client = logtest.LoggerServiceMock{}
	var result error = nil
	err := errors.New("invalid call")

	logtest.FatalfMock = func(format string, in ...interface{}) {
		result = err
	}

	vt.GetString = ""
	vt.GetStringFunc = nil

	validateClientSecretConfigPresence()
	if result == nil {
		test.Errorf("expected \n%s \n but got\nnil", err.Error())
	}
}

func TestValidateClientConfigPermissionsPresenceSuccess(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	log.Client = logtest.LoggerServiceMock{}
	var result error = nil
	logtest.FatalMock = func(in ...interface{}) { result = errors.New("invalid call") }
	vt.GetString = "clientPermissions"

	validateClientPermissionsConfigPresence()
	if result != nil {
		test.Errorf("expected no errors\n \n but got\n%s", result.Error())
	}
}

func TestValidateClientPermissionsConfigPresenceFailure(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	log.Client = logtest.LoggerServiceMock{}
	var result error = nil
	err := errors.New("invalid call")

	logtest.FatalfMock = func(format string, in ...interface{}) {
		result = err
	}

	vt.GetString = ""
	vt.GetStringFunc = nil

	validateClientPermissionsConfigPresence()
	if result == nil {
		test.Errorf("expected \n%s \n but got\nnil", err.Error())
	}
}

func TestHomedirSuccess(test *testing.T) {
	homedir.Client = HomedirServiceMock{}

	result := homeDir()
	if result != homedirPath {
		test.Errorf("expected \n%s \n but got\n%s", result, homedirPath)
	}
}

func TestHomedirFailure(test *testing.T) {
	homedir.Client = HomedirServiceMock{}

	var resultErr error
	homedirErr = errors.New("No homedir")
	logtest.FatalMock = func(in ...interface{}) { resultErr = homedirErr }
	log.Client = logtest.LoggerServiceMock{}

	homeDir()
	if resultErr != homedirErr {
		test.Errorf("expected \n%s \n but got\n%s", resultErr.Error(), homedirErr.Error())
	}
}

func testAccessorMethodFor(method func() string, stubValue string, test *testing.T) {
	vt.GetString = stubValue

	viper.Client = vt.ViperServiceMock{}

	result := method()
	if result != vt.GetString {
		test.Errorf("expected \n%s \n but got\n%s", vt.GetString, result)
	}
}

func TestUnixTimeAfter(test *testing.T) {
	testTime, err := time.Parse(time.RFC3339, "2020-12-10T09:15:23Z")

	if err != nil {
		test.Error(err)
	}

	// The unix timestamp equivalent of the testTime + 5 seconds
	var expectedResult int64 = 1607591728

	fiveSecodsDuration, err := time.ParseDuration("5s")
	if err != nil {
		test.Error(err)
	}

	t.Client = tt.TimeMock{}
	tt.TimeNowMockFunc = func() time.Time { return testTime }

	result := unixTimeAfter(fiveSecodsDuration)

	if result != expectedResult {
		test.Errorf("expected \n%d \n but got\n%d", expectedResult, result)
	}
}
