package logger

import "github.com/sirupsen/logrus"

var log *logrus.Logger

func InitLogger() {
	log = logrus.New()

	log.SetLevel(logrus.InfoLevel)
}

func GetLogger() *logrus.Logger {
	if log == nil {
		InitLogger()
	}

	return log
}
