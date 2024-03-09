package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func InitLogger() {
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetOutput(os.Stderr)
}

func GetLogger() *logrus.Logger {
	if log == nil {
		InitLogger()
	}
	
	return log
}
