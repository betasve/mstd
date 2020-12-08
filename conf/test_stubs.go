package conf

var configFileUsed, getString, homedirPath string
var configErr error = nil
var homedirErr error = nil

type ViperServiceMock struct{}

func (v ViperServiceMock) AddConfigPath(in string)     {}
func (v ViperServiceMock) AutomaticEnv()               {}
func (v ViperServiceMock) ConfigFileUsed() string      { return configFileUsed }
func (v ViperServiceMock) GetString(key string) string { return getString }
func (v ViperServiceMock) SetConfigFile(in string)     {}
func (v ViperServiceMock) SetConfigName(in string)     {}
func (v ViperServiceMock) ReadInConfig() error         { return configErr }

type LoggerServiceMock struct{}

var logMock = func(in string) {}
var fatalMock = func(in ...interface{}) {}
var fatalfMock = func(format string, in ...interface{}) {}

func (l LoggerServiceMock) Log(in string)                           { logMock(in) }
func (l LoggerServiceMock) Fatal(in ...interface{})                 { fatalMock(in) }
func (l LoggerServiceMock) Fatalf(format string, in ...interface{}) { fatalfMock(format, in) }

type HomedirServiceMock struct{}

func (h HomedirServiceMock) Dir() (string, error) {
	homedirPath = "/home/user/"
	return homedirPath, homedirErr
}
