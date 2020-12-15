package conf

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
