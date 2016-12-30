// Package logging provides structured logging instance
package logging

import (
	"github.com/uber-go/zap"
)

var logger zap.Logger

// Logger returns the logger instance
func Logger() zap.Logger {
	if logger == nil {
		logger = zap.New(zap.NewJSONEncoder())
	}
	return logger
}
