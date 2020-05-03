package logger

import (
	"os"
	"time"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log"
)

var (
	timestampFormat = log.TimestampFormat(
		func() time.Time { return time.Now().UTC() },
		"2006-01-02T15:04:05.000Z07:00",
	)
)

// GetLoggerInterface generates a logging interface for exporter
func GetLoggerInterface() log.Logger{
	var logger log.Logger
	logger = log.NewJSONLogger(os.Stderr)
	logger = level.NewFilter(logger, level.AllowInfo())
	logger = log.With(logger, "ts", timestampFormat)
	return logger
}
