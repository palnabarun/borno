package internal

import (
	"io"
	"strings"
	"unicode"

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

func slugify(title, location string) string {
	slug := ""
	for _, c := range strings.Join([]string{title, location}, " ") {
		if !unicode.IsDigit(c) && !unicode.IsLetter(c) && !unicode.IsSpace(c) && string(c) != "-" {
			continue
		}

		if unicode.IsSpace(c) || string(c) == "-" {
			slug += "-"
			continue
		}

		slug += string(unicode.ToLower(c))
	}

	return slug
}
