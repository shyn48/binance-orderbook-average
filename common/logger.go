package common

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	logger     *logrus.Logger
	loggerOnce sync.Once
)

func Logger() *logrus.Logger {
	loggerOnce.Do(func() {
		logLevel := logrus.TraceLevel
		optionalConfig, err := GetConfig()
		if err == nil {
			logLevel = logrus.AllLevels[optionalConfig.LogLevel]
		}

		logger = logrus.New()
		logger.SetLevel(logLevel)
	})

	return logger
}
