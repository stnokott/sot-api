// Package log provides logging utilities
package log

import "go.uber.org/zap"

var rootLogger = mustLogger()

func mustLogger() *zap.Logger {
	// zap.NewProduction()
	if logger, err := zap.NewDevelopment(); err != nil {
		panic(err)
	} else {
		return logger
	}
}

// ForModule creates a new child logger with the specified module name
func ForModule(name string) *zap.Logger {
	return rootLogger.With(zap.String("module", name))
}

// Sync flushes the logger buffer.
// It should be called at the end of main.
func Sync() error {
	return rootLogger.Sync()
}
