package conf

var configFileUsed, getString string
var getStringFunc func(in string) string
var getInt64Func func(in string) int64
var addConfigPathFunc, setConfigNameFunc, setCfgFilePathFunc func(in string)
var automaticEnvFunc func()
var setKeyValue func(key string, value interface{})
var configErr error
var writeConfigFunc func() error
var getInt64 int64

type ViperServiceMock struct{}

func (v ViperServiceMock) AddConfigPath(in string) { addConfigPathFunc(in) }
func (v ViperServiceMock) AutomaticEnv()           { automaticEnvFunc() }
func (v ViperServiceMock) ConfigFileUsed() string  { return configFileUsed }
func (v ViperServiceMock) GetInt64(key string) int64 {
	if getInt64Func != nil {
		return getInt64Func(key)
	} else {
		return getInt64
	}
}
func (v ViperServiceMock) GetString(key string) string {
	if getStringFunc != nil {
		return getStringFunc(key)
	} else {
		return getString
	}
}
func (v ViperServiceMock) Set(key string, value interface{}) { setKeyValue(key, value) }
func (v ViperServiceMock) SetConfigFile(in string)           { setCfgFilePathFunc(in) }
func (v ViperServiceMock) SetConfigName(in string)           { setConfigNameFunc(in) }
func (v ViperServiceMock) ReadInConfig() error               { return configErr }
func (v ViperServiceMock) WriteConfig() error                { return writeConfigFunc() }

type LoggerServiceMock struct{}

var logMock = func(in string) {}
var fatalMock = func(in ...interface{}) {}
var fatalfMock = func(format string, in ...interface{}) {}

func (l LoggerServiceMock) Log(in string)                           { logMock(in) }
func (l LoggerServiceMock) Fatal(in ...interface{})                 { fatalMock(in) }
func (l LoggerServiceMock) Fatalf(format string, in ...interface{}) { fatalfMock(format, in) }

type HomedirServiceMock struct{}

var homedirErr error
var homedirPath string

func (h HomedirServiceMock) Dir() (string, error) {
	homedirPath = "/home/user/"
	return homedirPath, homedirErr
}
