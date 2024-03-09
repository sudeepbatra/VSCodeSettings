package handler

import (
	"fmt"

	"github.com/sudeepbatra/alpha-hft/broker/smartapi"
	"github.com/sudeepbatra/alpha-hft/logger"
)

type AlphaSignalStorage struct{}

func NewAlphaSignalStorage() *AlphaSignalProcessor {
	logger.Log.Debug().Str("handler", "alpha_signal_storage").Msg("New NewAlphaSignalStorage Created!")
	return &AlphaSignalProcessor{}
}

func (p *AlphaSignalProcessor) StoreAlphaSignals(useProxy bool) {
	alphaSignalChannel := AlphaSignalManager.Subscribe()
	logger.Log.Debug().Str("handler", "alpha_signal_processor").Msg("AlphaSignalProcessor started!")

	}
}
