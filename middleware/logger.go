package middleware

import (
	"context"
)

// Logger interface provides full minimal necessary list to log data
type Logger interface {
	// WithContext necessary to update Logger entity with any useful information from the context.Context
	WithContext(ctx context.Context) Logger
	// WithField adds any value into the Logger entity with the given name
	WithField(name string, value interface{}) Logger
	// Printf is used to print the message
	Printf(format string, args ...interface{})
}

type noopLogger struct{}

var defaultLogger Logger = (*noopLogger)(nil)

func getLogger(logger Logger) Logger {
	if logger == nil {
		return defaultLogger
	}
	return logger
}

func (l *noopLogger) WithContext(_ context.Context) Logger {
	return l
}

func (l *noopLogger) WithField(_ string, _ interface{}) Logger {
	return l
}

func (l *noopLogger) Printf(_ string, _ ...interface{}) {}
