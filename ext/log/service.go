package log

import (
	"log"
)

var Client LoggerService = Logger{}

type LoggerService interface {
	Printf(format string, in ...interface{})
	Println(v ...interface{})
	Fatal(in ...interface{})
	Fatalf(format string, in ...interface{})
}

type Logger struct{}

func (l Logger) Printf(s string, v ...interface{}) {
	log.Printf(s, v...)
}

func (l Logger) Println(v ...interface{}) {
	log.Println(v...)
}

func (l Logger) Fatal(in ...interface{}) {
	log.Fatal(in...)
}

func (l Logger) Fatalf(format string, in ...interface{}) {
	log.Fatalf(format, in...)
}
