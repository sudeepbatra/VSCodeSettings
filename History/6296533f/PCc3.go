package handler

import (
	"github.com/sudeepbatra/alpha-hft/data"
	"github.com/sudeepbatra/alpha-hft/logger"
)

type AlphaSignalStorage struct{}

const (
	InstrumentTable = "alpha_signals"
)

func NewAlphaSignalStorage() *AlphaSignalProcessor {
	logger.Log.Debug().Str("handler", "alpha_signal_storage").Msg("New NewAlphaSignalStorage Created!")
	return &AlphaSignalProcessor{}
}

func (p *AlphaSignalProcessor) StoreAlphaSignals(useProxy bool) {
	alphaSignalChannel := AlphaSignalManager.Subscribe()

	logger.Log.Debug().Str("handler", "alpha_signal_processor").Msg("AlphaSignalProcessor started!")

	for alphaSignal := range alphaSignalChannel {
		logger.Log.Info().Interface("alphaSignal", alphaSignal).Msg("Storing Alpha Signal to the DB")
		err = BulkInsert(InstrumentTable, instrumentsData)
		data.InsertRecord()

	}
}
