package handler

import (
	"github.com/sudeepbatra/alpha-hft/data"
	"github.com/sudeepbatra/alpha-hft/logger"
)

type CandlestickStorage struct{}

const (
	CandlesticksTable = "live_candlesticks"
)

func NewCandlestickStorage() *CandlestickStorage {
	logger.Log.Debug().Str("handler", "candlestick_storage").Msg("New CandlestickStorage Created!")
	return &CandlestickStorage{}
}

func (p *CandlestickStorage) StoreCandlesticks() {
	logger.Log.Trace().Str("handler", "alpha_signal_processor").Msg("CandlestickStorage started!")

	completedCandles := CandlestickDataManager.Subscribe()

	for completedCandle := range completedCandles {
		logger.Log.Trace().Interface("completedCandle", completedCandle).Msg("Storing Completed Candle to the DB...")

		err := data.InsertRecord(CandlesticksTable, completedCandle)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error storing Alpha Signal to the DB")
		}

		logger.Log.Info().Interface("completedCandle", completedCandle).Msg("Storeds Alpha Signal to the DB!")
	}

	for alphaSignal := range alphaSignalChannel {
		logger.Log.Info().Interface("alphaSignal", alphaSignal).Msg("Storing Alpha Signal to the DB...")

		err := data.InsertRecord(AlphaSignalsTable, alphaSignal)

	}
}
