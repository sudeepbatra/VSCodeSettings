package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/sudeepbatra/alpha-hft/config"
)

var Log zerolog.Logger

func init() {
	logFile, err := os.OpenFile("alpha-hft.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	webSocketLogFile, err := os.OpenFile("websocket.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	zerolog.SetGlobalLevel(config.Config.LoggingConfig.LogLevel)

	zerolog.TimeFieldFormat = time.RFC3339

	fileWriter := zerolog.New(logFile).With().Timestamp().Logger()

	webSocketFileWriter := zerolog.New(webSocketLogFile).With().Timestamp().Logger()

	// Set up console logger with color output
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: false}

	// Create a multi-level logger that writes to both file and console
	multiLogger := zerolog.MultiLevelWriter(
		consoleWriter,
		fileWriter,
		webSocketFileWriter,
	)

	// Set up global logger
	Log = zerolog.New(multiLogger).With().Timestamp().Logger()
}
