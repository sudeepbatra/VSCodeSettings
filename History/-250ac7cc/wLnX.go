package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

var Log zerolog.Logger

func InitLogger() {
	log = logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.SetOutput(os.Stderr)
}

func GetLogger() *logrus.Logger {
	if log == nil {
		InitLogger()
	}

	return log
}

func init() {
	logFile, err := os.OpenFile("alpha-hft.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	// zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimeFieldFormat = time.RFC3339

	fileWriter := zerolog.New(logFile).With().Timestamp().Logger()

	// Set up console logger with color output
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: false}

	// Decide default log level
	// defaultLogLevel := config.Config.LoggingConfig.LogLevel // Use the appropriate default level

	// Create a multi-level logger that writes to both file and console
	multiLogger := zerolog.MultiLevelWriter(
		consoleWriter,
		fileWriter,
	)

	// Set up global logger
	Log = zerolog.New(multiLogger).With().Timestamp().Logger()

}
