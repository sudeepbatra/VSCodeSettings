package handler

import (
	"github.com/sudeepbatra/alpha-hft/data"
	"github.com/sudeepbatra/alpha-hft/logger"
)

type LiveCandlesticksStorage struct{}

const (
	LiveCandlesticksTable = "live_candlesticks"
)

func NewLiveCandlesticksStorage() *LiveCandlesticksStorage {
	logger.Log.Info().Str("handler", "live_candlesticks_storage").Msg("New LiveCandlesticksStorage Created!")
	return &LiveCandlesticksStorage{}
}

func (p *LiveCandlesticksStorage) StoreCandlesticks() {
	logger.Log.Info().Str("handler", "live_candlesticks_storage").Msg("LiveCandlesticksStorage started!")

	completedCandles := CandlestickDataManager.Subscribe()

	for completedCandle := range completedCandles {
		logger.Log.Info().Interface("completedCandle", completedCandle).Msg("Storing Live Candlestick to the DB...")

		err := data.InsertRecord(LiveCandlesticksTable, completedCandle)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error storing live candlestick to the DB")
		}

		logger.Log.Info().Interface("completedCandle", completedCandle).Msg("Stored Live Candlestick to the DB!")
	}
}