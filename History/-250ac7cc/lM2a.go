package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func InitLogger() {
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
}

func GetLogger() *zerolog.Logger {
	if Log == nil {
		InitLogger()
	}

	return &Log
}
