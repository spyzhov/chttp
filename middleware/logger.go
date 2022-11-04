package transport

import (
	"context"
)

type Logger interface {
	WithContext(ctx context.Context) Logger
	WithField(name string, value interface{}) Logger
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
