package logger

import (
	"log"
)

var Client LoggerService = Logger{}

type LoggerService interface {
	Log(in string)
	Fatal(in ...interface{})
	Fatalf(format string, in ...interface{})
}

type Logger struct{}

func (l Logger) Log(in string) {
	log.Println(in)
}

func (l Logger) Fatal(in ...interface{}) {
	log.Fatal(in...)
}

func (l Logger) Fatalf(format string, in ...interface{}) {
	log.Fatalf(format, in...)
}
