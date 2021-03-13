package internal

import (
	"io"

	"github.com/sirupsen/logrus"
)

// LoggerOpts describes options for the logger
type LoggerOpts struct {
	Out io.Writer
}

// NewLogger returns a new logger instance
func NewLogger(opts *LoggerOpts) *logrus.Logger {
	logger := logrus.New()

	logger.Out = opts.Out
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	return logger
}
