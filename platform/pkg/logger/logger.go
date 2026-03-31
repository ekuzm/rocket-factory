package logger

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	logger   *logrus.Logger
	initOnce sync.Once
)

func Init(level string, asJSON bool) {
	initOnce.Do(func() {
		logger = logrus.New()

		logger.SetReportCaller(true)

		var formatter logrus.Formatter
		if asJSON {
			formatter = &logrus.JSONFormatter{
				TimestampFormat: "02-01-2006 15:04:05",
				CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
					return "", filepath.Base(f.File) + ":" + strconv.Itoa(f.Line)
				},
				PrettyPrint: true,
			}
		} else {
			formatter = &logrus.TextFormatter{
				TimestampFormat:        "02-01-2006 15:04:05",
				FullTimestamp:          true,
				DisableLevelTruncation: true,
				CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
					return "", filepath.Base(f.File)
				},
			}
		}

		logger.SetFormatter(formatter)

		logger.SetLevel(parseLevel(level))
		logger.SetOutput(os.Stdout)
	})
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(fields)
}

func parseLevel(level string) logrus.Level {
	switch strings.ToUpper(level) {
	case "INFO":
		return logrus.InfoLevel
	case "ERROR":
		return logrus.ErrorLevel
	case "WARNING", "WARN":
		return logrus.WarnLevel
	case "FATAL":
		return logrus.FatalLevel
	default:
		return logrus.DebugLevel
	}
}
