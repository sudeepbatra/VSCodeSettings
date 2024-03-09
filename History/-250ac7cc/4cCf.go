package logger

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var (
	logInstance *zerolog.Logger
	logOnce     sync.Once
)

func InitLogger() {
	logOnce.Do(func() {
		logFile, err := os.OpenFile("alpha-hft.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			panic(err)
		}

		zerolog.TimeFieldFormat = time.RFC3339

		fileWriter := zerolog.New(logFile).With().Timestamp().Logger()

		// Create a multi-level logger that writes to both file and console
		multiLogger := zerolog.MultiLevelWriter(fileWriter)

		// Set up global logger
		logger := zerolog.New(multiLogger).With().Timestamp().Logger()
		logInstance = &logger
	})
}

func GetLogger() *zerolog.Logger {
	if logInstance == nil {
		InitLogger()
	}

	return logInstance
}
