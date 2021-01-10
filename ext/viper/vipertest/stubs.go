package vipertest

var AddConfigPathFunc, SetConfigNameFunc func(in string)
var AutomaticEnvFunc = func() {}
var ConfigFileUsed, GetString string
var GetStringFunc func(in string) string
var GetInt64Func func(in string) int64
var SetKeyValue func(key string, value interface{})
var SetCfgFilePathFunc = func(in string) {}
var ConfigErr error
var WriteConfigFunc func() error
var GetInt64 int64

type ViperServiceMock struct{}

func (v ViperServiceMock) AddConfigPath(in string) { AddConfigPathFunc(in) }
func (v ViperServiceMock) AutomaticEnv()           { AutomaticEnvFunc() }
func (v ViperServiceMock) ConfigFileUsed() string  { return ConfigFileUsed }
func (v ViperServiceMock) GetInt64(key string) int64 {
	if GetInt64Func != nil {
		return GetInt64Func(key)
	} else {
		return GetInt64
	}
}
func (v ViperServiceMock) GetString(key string) string {
	if GetStringFunc != nil {
		return GetStringFunc(key)
	} else {
		return GetString
	}
}
func (v ViperServiceMock) Set(key string, value interface{}) { SetKeyValue(key, value) }
func (v ViperServiceMock) SetConfigFile(in string)           { SetCfgFilePathFunc(in) }
func (v ViperServiceMock) SetConfigName(in string)           { SetConfigNameFunc(in) }
func (v ViperServiceMock) ReadInConfig() error               { return ConfigErr }
func (v ViperServiceMock) WriteConfig() error                { return WriteConfigFunc() }
