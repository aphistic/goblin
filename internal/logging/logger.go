// Package logging is the logger the Goblin binary uses internally.
package logging

import "fmt"

// Logger is an interface representint a logger.
type Logger interface {
	Printf(string, ...interface{})
}

// NilLogger is a Logger that does nothing.
type NilLogger struct{}

var _ Logger = &NilLogger{}

// NewNilLogger creates a new nil logger.
func NewNilLogger() *NilLogger {
	return &NilLogger{}
}

// Printf does nothing.
func (l *NilLogger) Printf(format string, args ...interface{}) {}

// PrintfLogger uses fmt.Printf for logging.
type PrintfLogger struct{}

var _ Logger = &PrintfLogger{}

// NewPrintfLogger creates a new fmt.Printf logger.
func NewPrintfLogger() *PrintfLogger {
	return &PrintfLogger{}
}

// Printf uses fmt.Printf to log messages.
func (l *PrintfLogger) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
