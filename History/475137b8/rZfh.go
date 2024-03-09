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
	completedCandles := handler.CandlestickDataManager.Subscribe()
	for completedCandle := range completedCandles {
		logger.Log.Debug().Msg("Receieved Completed Candle: " + fmt.Sprintf("%+v", completedCandle))


	alphaSignalChannel := AlphaSignalManager.Subscribe()

	logger.Log.Debug().Str("handler", "alpha_signal_processor").Msg("AlphaSignalProcessor started!")

	for alphaSignal := range alphaSignalChannel {
		logger.Log.Info().Interface("alphaSignal", alphaSignal).Msg("Storing Alpha Signal to the DB...")

		err := data.InsertRecord(AlphaSignalsTable, alphaSignal)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error storing Alpha Signal to the DB")
		}

		logger.Log.Info().Interface("alphaSignal", alphaSignal).Msg("Storeds Alpha Signal to the DB!")
	}
}
