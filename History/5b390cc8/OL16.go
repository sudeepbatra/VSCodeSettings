package handler

import (
	"github.com/sudeepbatra/alpha-hft/broker/smartapi"
	"github.com/sudeepbatra/alpha-hft/logger"
)

type AlphaSignalProcessor struct{}

func NewAlphaSignalProcessor() *AlphaSignalProcessor {
	logger.Log.Debug().Str("handler", "alpha_signal_processor").Msg("New AlphaSignalProcessor Created!")
	return &AlphaSignalProcessor{}
}

func (p *AlphaSignalProcessor) ProcessAlphaSignals() {
	alphaSignalChannel := AlphaSignalManager.Subscribe()

	logger.Log.Debug().Str("handler", "alpha_signal_processor").Msg("AlphaSignalProcessor started!")

	for alphaSignal := range alphaSignalChannel {
		logger.Log.Info().Interface("alphaSignal", alphaSignal).Msg("Processing Alpha Signal")

		// Place order based on the Alpha Signal
		if alphaSignal.Signal == "LONG" {
			tradingAPI := smartapi.SmartApiBrokers["trading"]
			s.PlaceLongOrder(alphaSignal)
			OrderManager.PlaceOrder(alphaSignal.Symbol, alphaSignal.Quantity, alphaSignal.Price, "BUY")
		}
	}
}
