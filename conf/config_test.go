package conf

import (
	"errors"
	"github.com/betasve/mstd/homedir"
	homedirtest "github.com/betasve/mstd/homedir/homedirtest"
	t "github.com/betasve/mstd/time"
	tt "github.com/betasve/mstd/time/timetest"
	"github.com/betasve/mstd/viper"
	vt "github.com/betasve/mstd/viper/vipertest"
	"testing"
	"time"
)

func TestInitConfigSuccess(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	cfgFilePath := "file/path"
	vt.GetString = "testViperString"

	config := Config{}
	err := config.InitConfig(cfgFilePath)

	if config.ClientId() != vt.GetString || err != nil {
		test.Errorf(
			"expected config client id to have\n%s\nbut got\n%s",
			vt.GetString,
			config.ClientId(),
		)
	}
}

func TestInitConfigFailedToSetConfig(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	cfgFilePath := "file/path"
	vt.GetString = "testViperString"

	vt.ConfigErr = errors.New("error reading config file")

	config := Config{}
	err := config.InitConfig(cfgFilePath)

	if err != vt.ConfigErr {
		test.Errorf(
			"expected config client id to have\n%s\nbut got\n%s",
			vt.GetString,
			config.ClientId(),
		)
	}
}

func TestInitConfigFailedToValidateConfig(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	cfgFilePath := "file/path"
	vt.GetStringFunc = func(in string) string {
		if in == "client_id" {
			return "ClientId"
		} else {
			return ""
		}
	}

	config := Config{}
	err := config.InitConfig(cfgFilePath)

	if err != vt.ConfigErr {
		test.Errorf(
			"expected config client id to have\n%s\nbut got\n%s",
			vt.GetString,
			config.ClientId(),
		)
	}
}

func TestGetClientId(test *testing.T) {
	stub := "testClientId"
	cfg := Config{clientId: stub}

	testAccessorMethodFor(cfg.ClientId, stub, test)
}

func TestGetClientSecret(test *testing.T) {
	stub := "testClientSecret"
	cfg := Config{clientSecret: stub}

	testAccessorMethodFor(cfg.ClientSecret, stub, test)
}

func TestGetClientPermissions(test *testing.T) {
	stub := "testClientPermissions"
	cfg := Config{permissions: stub}

	testAccessorMethodFor(cfg.ClientPermissions, stub, test)
}

func TestGetClientAccessToken(test *testing.T) {
	stub := "testClientAccessToken"
	cfg := Config{accessToken: stub}

	testAccessorMethodFor(cfg.ClientAccessToken, stub, test)
}

func TestGetAuthCallbackHost(test *testing.T) {
	stub := "testClientAuthCallbackHost"
	cfg := Config{authCallbackHost: stub}

	testAccessorMethodFor(cfg.AuthCallbackHost, stub, test)
}

func TestGetAuthCallbackPath(test *testing.T) {
	stub := "testClientAuthCallbackPath"
	cfg := Config{authCallbackPath: stub}

	testAccessorMethodFor(cfg.AuthCallbackPath, stub, test)
}

func TestGetClientRefreshToken(test *testing.T) {
	stub := "testClientRefreshToken"
	cfg := Config{refreshToken: stub}

	testAccessorMethodFor(cfg.ClientRefreshToken, stub, test)
}

func TestGetClientAccessTokenExpirySuccess(test *testing.T) {
	now := time.Now()
	cfg := Config{accessTokenExpiresAt: now}

	result := cfg.ClientAccessTokenExpiresAt()

	if result != now {
		test.Errorf("expected\n%s\nbut got\n%s", now, result)
	}
}

func TestGetClientRefreshTokenExpirySuccess(test *testing.T) {
	now := time.Now()
	cfg := Config{refreshTokenExpiresAt: now}

	result := cfg.ClientRefreshTokenExpiresAt()

	if result != now {
		test.Errorf("expected %s\nbut got\n%s", now, result)
	}
}

func TestClientId(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	vt.GetString = "clientId"
	vt.GetStringFunc = nil

	result := clientId()
	if result != vt.GetString {
		test.Errorf("expected\n%s\nbut got\n%s", vt.GetString, result)
	}
}

func TestClientSecret(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	vt.GetString = "clientSecret"
	vt.GetStringFunc = nil

	result := clientSecret()
	if result != vt.GetString {
		test.Errorf("expected\n%s\nbut got\n%s", vt.GetString, result)
	}
}

func TestClientPermissions(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	vt.GetString = "clientPermissions"
	vt.GetStringFunc = nil

	result := clientPermissions()
	if result != vt.GetString {
		test.Errorf("expected\n%s\nbut got\n%s", vt.GetString, result)
	}
}

func TestClientAccessToken(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	vt.GetString = "clientAccessToken"
	vt.GetStringFunc = nil

	result := clientAccessToken()
	if result != vt.GetString {
		test.Errorf("expected\n%s\nbut got\n%s", vt.GetString, result)
	}
}

func TestClientRefreshToken(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	vt.GetString = "clientRefreshToken"
	vt.GetStringFunc = nil

	result := clientRefreshToken()
	if result != vt.GetString {
		test.Errorf("expected\n%s\nbut got\n%s", vt.GetString, result)
	}
}

func TestAuthCallbackHost(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	vt.GetString = "authCallbackHost"
	vt.GetStringFunc = nil

	result := authCallbackHost()
	if result != vt.GetString {
		test.Errorf("expected\n%s\nbut got\n%s", vt.GetString, result)
	}
}

func TestAuthCallbackPath(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	vt.GetString = "authCallbackPath"
	vt.GetStringFunc = nil

	result := authCallbackPath()
	if result != vt.GetString {
		test.Errorf("expected\n%s\nbut got\n%s", vt.GetString, result)
	}
}

func TestClientAccessTokenExpiry(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	vt.GetInt64 = 0
	vt.GetInt64Func = nil
	expected := time.Unix(0, 0)

	result := clientAccessTokenExpiry()
	if result != expected {
		test.Errorf("expected\n%s\nbut got\n%s", expected, result)
	}
}

func TestClientRefreshTokenExpiry(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	vt.GetInt64 = 0
	vt.GetInt64Func = nil
	expected := time.Unix(0, 0)

	result := clientRefreshTokenExpiry()
	if result != expected {
		test.Errorf("expected\n%s\nbut got\n%s", expected, result)
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

	cfg := Config{}
	err := cfg.SetClientAccessToken(expectedVal)

	if err != nil ||
		key != defaultAccessTokenConfig ||
		value != expectedVal ||
		!writeCalled {

		test.Errorf(
			"expected\nkey: %s,\nvalue: %s,\nwrite called: true\n"+
				"config accessToken set to %s"+
				"but got\nkey: %s,\nvalue: %s,\nwrite called: %t"+
				"config set to %s",
			defaultAccessTokenConfig,
			expectedVal,
			expectedVal,
			key,
			value,
			writeCalled,
			cfg.accessToken,
		)
	}
}

func TestSetClientAccessTokenFailure(test *testing.T) {
	err := errors.New("could not write config")
	viper.Client = vt.ViperServiceMock{}
	vt.SetKeyValue = func(k string, v interface{}) {}
	vt.WriteConfigFunc = func() error { return err }

	cfg := Config{}
	setErr := cfg.SetClientAccessToken("")

	if setErr != err {
		test.Errorf("expected error\n%s\nbut got\n%s", err, setErr)
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

	cfg := Config{}
	cfg.SetClientRefreshToken(expectedVal)

	if cfg.refreshToken != expectedVal {
		test.Errorf("\nexpected:\n%s\nbut got\n%s", expectedVal, cfg.refreshToken)
	}

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

	cfg := Config{}
	expectedErr := cfg.SetClientRefreshToken("")

	if expectedErr != err {
		test.Errorf("expected error\n %s but got %s", expectedErr, err)
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

	cfg := Config{}
	setErr := cfg.SetClientAccessTokenExpirySeconds(5)

	expectedTime := time.Unix(expectedResult, 0)

	if cfg.accessTokenExpiresAt != expectedTime {
		test.Errorf("\nexpected:\n%s\nbut got\n%s", expectedTime, cfg.accessTokenExpiresAt)
	}

	if setErr != nil {
		test.Errorf("\nexpected:\nno errors\nbut got\n%s", setErr)
	}

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
	viper.Client = vt.ViperServiceMock{}

	expectedErr := errors.New("Could not write config")
	vt.WriteConfigFunc = func() error { return expectedErr }

	cfg := Config{}
	err := cfg.SetClientAccessTokenExpirySeconds(5)

	if err != expectedErr {
		test.Errorf("expected\n%s\nbut got\n%s", err, expectedErr)
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
	cfg := Config{}

	setErr := cfg.SetClientRefreshTokenExpirySeconds(5)
	expectedTime := time.Unix(expectedResult, 0)

	if cfg.refreshTokenExpiresAt != expectedTime {
		test.Errorf("\nexpected:\n%s\nbut got\n%s", expectedTime, cfg.accessTokenExpiresAt)
	}

	if setErr != nil {
		test.Errorf("\nexpected:\nno errors\nbut got\n%s", setErr)
	}
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

	expectedErr := errors.New("Could not write config")
	vt.WriteConfigFunc = func() error { return expectedErr }
	cfg := Config{}
	err := cfg.SetClientRefreshTokenExpirySeconds(5)

	if err != expectedErr {
		test.Errorf("expected \n%s \nbut got\n%s", err, expectedErr)
	}
}

func TestPopulateConfigValues(test *testing.T) {
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

	cfg := Config{}
	cfg.populateConfigValues()

	if cfg.ClientId() != "client_id" {
		test.Errorf("\nexpcted client_id\nbut got\n%s", cfg.ClientId())
	}

	if cfg.ClientSecret() != "client_secret" {
		test.Errorf("\nexpcted client_secret\nbut got\n%s", cfg.ClientSecret())
	}

	if cfg.ClientPermissions() != "client_permissions" {
		test.Errorf("\nexpcted client_permissions\nbut got\n%s", cfg.ClientPermissions())
	}

	if cfg.ClientAccessToken() != "default_access_token" {
		test.Errorf("\nexpcted default_access_token\nbut got\n%s", cfg.ClientAccessToken())
	}

	if cfg.ClientRefreshToken() != "default_refresh_token" {
		test.Errorf("\nexpcted default_refresh_token\nbut got\n%s", cfg.ClientRefreshToken())
	}

	if cfg.ClientAccessTokenExpiresAt() != time.Unix(1607590000, 0) {
		test.Errorf("\nexpcted %v\nbut got\n%s", time.Unix(1607590000, 0), cfg.ClientAccessTokenExpiresAt())
	}

	if cfg.ClientRefreshTokenExpiresAt() != time.Unix(1607590011, 0) {
		test.Errorf("\nexpcted %v\nbut got\n%s", time.Unix(1607590011, 0), cfg.ClientRefreshTokenExpiresAt())
	}
}

func TestSecondsToDurationSuccess(test *testing.T) {
	t.Client = tt.TimeMock{}
	dur := time.Since(time.Now())
	tt.TimeParseDurationMockFunc = func(s string) (time.Duration, error) { return dur, nil }

	result, err := secondsToDuration(42)
	if result != dur || err != nil {
		test.Errorf("expected \n%s \n but got\n%s", dur.String(), result.String())
	}
}

func TestSecondsToDurationFailure(test *testing.T) {
	t.Client = tt.TimeMock{}
	err := errors.New("cannot parse")
	tt.TimeParseDurationMockFunc = func(s string) (time.Duration, error) {
		return time.Since(time.Now()), err
	}

	_, durErr := secondsToDuration(42)
	if durErr != err {
		test.Errorf("expected \n%s \n but got\n%s", err, durErr)
	}
}

func TestSetViperConfigWithConfigFilePath(test *testing.T) {
	cfgFilePath := "file/path"
	vt.ConfigFileUsed = "conf.yml"
	vt.ConfigErr = nil

	viper.Client = vt.ViperServiceMock{}

	err := setViperConfig(cfgFilePath)
	if err != nil {
		test.Errorf("expected\nno errors\nbut got\n%s", err)
	}
}

func TestSetViperConfigWithoutConfigFilePath(test *testing.T) {
	var addConfigPathResult string
	var setConfigNameResult string
	cfgFilePath := ""
	vt.ConfigFileUsed = "conf.yml"

	viper.Client = vt.ViperServiceMock{}
	vt.AddConfigPathFunc = func(in string) { addConfigPathResult = "homedir" }
	vt.SetConfigNameFunc = func(in string) { setConfigNameResult = vt.ConfigFileUsed }

	err := setViperConfig(cfgFilePath)
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

	if err != nil {
		test.Errorf("expected\nnil\nbut got\n%s", err)
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
	vt.ConfigErr = nil
	viper.Client = vt.ViperServiceMock{}

	err := readConfigFile()

	if err != nil {
		test.Errorf("expected\nnil\nbut got\n%s", err)
	}
}

func TestReadConfigFileFailure(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	vt.ConfigErr = errors.New("Cannot read in config")

	err := readConfigFile()

	if err != vt.ConfigErr {
		test.Errorf("expected error\n%s\nbut got\n%s", vt.ConfigErr, err)
	}
}

func TestValidateConfigFileAttributesSuccess(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	vt.GetString = ""
	vt.GetStringFunc = func(in string) string {
		switch in {
		case defaultClientIdConfig:
			return "test-client-id"
		case defaultClientSecretConfig:
			return "test-client-secret"
		case defaultPermissionsConfig:
			return "test.client.permissions"
		case defaultAuthCallbackHost:
			return "http://test-host:3000"
		case defaultAuthCallbackPath:
			return "auth/path"
		default:
			return ""
		}
	}

	err := validateConfigFileAttributes()
	if err != nil {
		test.Errorf("expected no errors\n \n but got\n%s", err)
	}
}

func TestValidateConfigFileAttributesClientIdFailure(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	vt.GetStringFunc = func(key string) string {
		if key == defaultClientIdConfig {
			return ""
		} else {
			return "clientDetail"
		}
	}

	err := validateConfigFileAttributes()
	if err == nil {
		test.Error("\nexpected\nerrors\nbut got\nnil")
	}
}

func TestValidateConfigFileAttributesClientSecretFailure(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	vt.GetStringFunc = func(key string) string {
		if key == defaultClientSecretConfig {
			return ""
		} else {
			return "clientDetail"
		}
	}

	err := validateConfigFileAttributes()
	if err == nil {
		test.Error("\nexpected\nerrors\nbut got\nnil")
	}
}

func TestValidateConfigFileAttributesClientPermissionsFailure(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	vt.GetStringFunc = func(key string) string {
		if key == defaultPermissionsConfig {
			return ""
		} else {
			return "clientDetail"
		}
	}

	err := validateConfigFileAttributes()
	if err == nil {
		test.Error("\nexpected\nerrors\nbut got\nnil")
	}
}

func TestValidateConfigFileAttributesAuthCallbackHostFailure(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	vt.GetStringFunc = func(key string) string {
		if key == defaultAuthCallbackHost {
			return ""
		} else {
			return "clientDetail"
		}
	}

	err := validateConfigFileAttributes()
	if err == nil {
		test.Error("\nexpected\nerrors\nbut got\nnil")
	}
}

func TestValidateConfigFileAttributesAuthCallbackPathFailure(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	vt.GetStringFunc = func(key string) string {
		if key == defaultAuthCallbackPath {
			return ""
		} else {
			return "clientDetail"
		}
	}

	err := validateConfigFileAttributes()
	if err == nil {
		test.Error("\nexpected\nerrors\nbut got\nnil")
	}
}

func TestValidateClientIdConfigPresenceSuccess(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	vt.GetString = "clientId"

	err := validateClientIdConfigPresence()
	if err != nil {
		test.Errorf("expected no errors\nbut got\n%s", err)
	}
}

func TestValidateClientIdConfigPresenceFailure(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}

	vt.GetString = ""
	vt.GetStringFunc = nil

	err := validateClientIdConfigPresence()
	if err == nil {
		test.Error("expected errors\nbut got\nnil")
	}
}

func TestValidateClientSecretConfigPresenceSuccess(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	vt.GetString = "clientSecret"

	err := validateClientSecretConfigPresence()
	if err != nil {
		test.Errorf("expected no errors\n \n but got\n%s", err)
	}
}

func TestValidateClientSecretConfigPresenceFailure(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}

	vt.GetString = ""
	vt.GetStringFunc = nil

	err := validateClientSecretConfigPresence()
	if err == nil {
		test.Error("expected errors\nbut got\nnil")
	}
}

func TestValidateClientConfigPermissionsPresenceSuccess(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}
	vt.GetString = "clientPermissions"

	err := validateClientPermissionsConfigPresence()
	if err != nil {
		test.Errorf("expected no errors\n \n but got\n%s", err)
	}
}

func TestValidateClientPermissionsConfigPresenceFailure(test *testing.T) {
	viper.Client = vt.ViperServiceMock{}

	vt.GetString = ""
	vt.GetStringFunc = nil

	err := validateClientPermissionsConfigPresence()
	if err == nil {
		test.Error("expected errors\nbut got\nnil")
	}
}

func TestHomedirSuccess(test *testing.T) {
	homedir.Client = homedirtest.HomedirServiceMock{}

	_, err := homeDir()
	if err != nil {
		test.Errorf("expected no errors\nbut got\n%s", err)
	}
}

func TestHomedirFailure(test *testing.T) {
	homedir.Client = homedirtest.HomedirServiceMock{}
	homedirtest.HomedirErr = errors.New("No homedir")

	_, err := homeDir()
	if err != homedirtest.HomedirErr {
		test.Errorf("expected \n%s \n but got\n%s", err, homedirtest.HomedirErr)
	}
}

func testAccessorMethodFor(method func() string, expected string, test *testing.T) {
	result := method()
	if result != expected {
		test.Errorf("expected \n%s \n but got\n%s", expected, result)
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
