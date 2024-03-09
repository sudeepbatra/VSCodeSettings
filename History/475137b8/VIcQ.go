package handler

import "github.com/sudeepbatra/alpha-hft/logger"

type CandlestickStorage struct{}

const (
	CandlesticksTable = "live_candlesticks"
)

func NewCandlestickStorage() *CandlestickStorage {
	logger.Log.Debug().Str("handler", "alpha_signal_storage").Msg("New NewAlphaSignalStorage Created!")
	return &CandlestickStorage{}
}
