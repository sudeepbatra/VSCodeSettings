package handler

import (
	"sync"

	"github.com/sudeepbatra/alpha-hft/data"
	"github.com/sudeepbatra/alpha-hft/logger"
)

type AlphaSignalStorage struct {
	wg sync.WaitGroup
}

const (
	AlphaSignalsTable = "alpha_signals"
)

func NewAlphaSignalStorage() *AlphaSignalStorage {
	logger.Log.Debug().
		Str("handler", "alpha_signal_storage").
		Msg("New NewAlphaSignalStorage Created!")

	return &AlphaSignalStorage{}
}

func (a *AlphaSignalStorage) StoreAlphaSignals() {
	logger.Log.Debug().
		Str("handler", "alpha_signal_storage").
		Msg("StoreAlphaSignals started! Subscribing to Alpha Signals channel...")

	alphaSignalChannel := AlphaSignalManager.Subscribe()

	for alphaSignal := range alphaSignalChannel {
		logger.Log.Debug().
			Interface("alphaSignal", alphaSignal).
			Msg("Storing Alpha Signal to the DB...")

		err := data.InsertRecord(AlphaSignalsTable, alphaSignal)
		if err != nil {
			logger.Log.Error().
				Err(err).
				Msg("Error storing Alpha Signal to the DB")
		} else {
			logger.Log.Debug().
				Interface("alphaSignal", alphaSignal).
				Msg("Stored Alpha Signal to the DB!")
		}
	}

	a.wg.Done()
}
