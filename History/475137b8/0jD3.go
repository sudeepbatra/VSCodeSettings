package handler

import (
	"github.com/sudeepbatra/alpha-hft/data"
	"github.com/sudeepbatra/alpha-hft/logger"
)

type LiveCandlesticksStorage struct{}

const (
	LiveCandlesticksTable = "live_candlesticks"
)

func NewCandlestickStorage() *LiveCandlesticksStorage {
	logger.Log.Trace().Str("handler", "live_candlestick_storage").Msg("New LiveCandlesticksStorage Created!")
	return &LiveCandlesticksStorage{}
}

func (p *LiveCandlesticksStorage) StoreCandlesticks() {
	logger.Log.Trace().Str("handler", "live_candlestick_storage").Msg("LiveCandlesticksStorage started!")

	completedCandles := CandlestickDataManager.Subscribe()

	for completedCandle := range completedCandles {
		logger.Log.Trace().Interface("completedCandle", completedCandle).Msg("Storing Completed Candle to the DB...")

		err := data.InsertRecord(LiveCandlesticksTable, completedCandle)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error storing Alpha Signal to the DB")
		}

		logger.Log.Info().Interface("completedCandle", completedCandle).Msg("Stored Completed Candle to the DB!")
	}
}
