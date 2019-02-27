package logging

import "fmt"

type Logger interface {
	Printf(string, ...interface{})
}

type NilLogger struct {}

var _ Logger = &NilLogger{}

func NewNilLogger() *NilLogger {
	return &NilLogger{}
}

func (l *NilLogger) Printf(format string, args ...interface{}) {

}

type PrintfLogger struct {}

var _ Logger = &PrintfLogger{}

func NewPrintfLogger() *PrintfLogger {
	return &PrintfLogger{}
}

func (l *PrintfLogger) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}