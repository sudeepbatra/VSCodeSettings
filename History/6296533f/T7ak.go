package handler

import (
	"github.com/sudeepbatra/alpha-hft/data"
	"github.com/sudeepbatra/alpha-hft/logger"
)

type AlphaSignalStorage struct{}

const (
	AlphaSignalsTable = "alpha_signals"
)

func NewAlphaSignalStorage() *AlphaSignalStorage {
	logger.Log.Debug().Str("handler", "alpha_signal_storage").Msg("New NewAlphaSignalStorage Created!")
	return &AlphaSignalStorage{}
}

func (p *AlphaSignalStorage) StoreAlphaSignals() {
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
