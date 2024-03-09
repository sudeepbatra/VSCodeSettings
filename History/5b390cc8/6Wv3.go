package handler

import (
	"fmt"

	"github.com/sudeepbatra/alpha-hft/broker/smartapi"
	"github.com/sudeepbatra/alpha-hft/data"
	"github.com/sudeepbatra/alpha-hft/logger"
)

type AlphaSignalProcessor struct{}

func NewAlphaSignalProcessor() *AlphaSignalProcessor {
	logger.Log.Debug().Str("handler", "alpha_signal_processor").Msg("New AlphaSignalProcessor Created!")
	return &AlphaSignalProcessor{}
}

func (p *AlphaSignalProcessor) ProcessAlphaSignals(useProxy bool) {
	alphaSignalChannel := AlphaSignalManager.Subscribe()
	client := data.GetClient(useProxy)

	logger.Log.Debug().Str("handler", "alpha_signal_processor").Msg("AlphaSignalProcessor started!")

	for alphaSignal := range alphaSignalChannel {
		logger.Log.Info().Interface("alphaSignal", alphaSignal).Msg("Processing Alpha Signal")

		// Place order based on the Alpha Signal
		if alphaSignal.Signal == "LONG" {
			tradingAPI := smartapi.SmartApiBrokers["trading"]

			price := fmt.Sprintf("%.2f", alphaSignal.Price)

			order, err := tradingAPI.PlaceOrder(
				client,
				"NORMAL",
				alphaSignal.Token,
				alphaSignal.Symbol,
				alphaSignal.Exchange,
				"BUY", "BUY",
				"LIMIT",
				"INTRADAY",
				"DAY",
				price,
				"1",
				"0",
				"0")
		}
	}
}
