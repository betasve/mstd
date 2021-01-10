package loggertest

type LoggerServiceMock struct{}

var PrintfMock = func(s string, v ...interface{}) {}
var PrintlnMock = func(v ...interface{}) {}
var FatalMock = func(in ...interface{}) {}
var FatalfMock = func(format string, in ...interface{}) {}

func (l LoggerServiceMock) Printf(s string, v ...interface{})       { PrintfMock(s, v...) }
func (l LoggerServiceMock) Println(v ...interface{})                { PrintlnMock(v...) }
func (l LoggerServiceMock) Fatal(in ...interface{})                 { FatalMock(in) }
func (l LoggerServiceMock) Fatalf(format string, in ...interface{}) { FatalfMock(format, in) }
