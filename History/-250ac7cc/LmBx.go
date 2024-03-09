package logger

import (
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
)

var (
	logOnce     sync.Once
	zerologOnce sync.Once
	Log         zerolog.Logger
	logrusLog   *logrus.Logger
)

func InitLogger() {
	logOnce.Do(func() {
		logFile, err := os.OpenFile("alpha-hft.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}

		zerolog.TimeFieldFormat = time.RFC3339

		fileWriter := zerolog.New(logFile).With().Timestamp().Logger()

		// Create a multi-level logger that writes to both file and console
		multiLogger := zerolog.MultiLevelWriter(fileWriter)

		// Set up global logger
		Log = zerolog.New(multiLogger).With().Timestamp().Logger()

		// Logrus initialization
		logrusLog = logrus.New()
		logrusLog.SetLevel(logrus.DebugLevel)
		logrusLog.SetOutput(ioutil.Discard) // Set output to discard to avoid double printing
	})
}

func GetLogger() *zerolog.Logger {
	if Log == (zerolog.Logger{}) {
		InitLogger()
	}

	return &Log
}

// You can add similar functions for Logrus as needed
func GetLogrusLogger() *logrus.Logger {
	if logrusLog == nil {
		InitLogger()
	}

	return logrusLog
}
