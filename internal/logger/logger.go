package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogrusLogger() *logrus.Logger {
	logger := &logrus.Logger{
		Out: os.Stdout,
		Formatter: &logrus.TextFormatter{
			ForceColors:            true,
			FullTimestamp:          true,
			DisableLevelTruncation: true,
			PadLevelText:           true,
			TimestampFormat:        "2006.01.02 15:04:05", // default go time format
		},
		// ReportCaller: true,
		Level: logrus.DebugLevel,
	}
	return logger
}
