package common

import (
	"fmt"
	"time"
)

var ErrInvalidStrategy = fmt.Errorf("invalid strategy")

type OHLCVTimeSeries struct {
	Token     int         `json:"token"`
	Timestamp []time.Time `json:"timestamp"`
	Open      []float64   `json:"open"`
	High      []float64   `json:"high"`
	Low       []float64   `json:"low"`
	Close     []float64   `json:"close"`
	Volume    []int       `json:"volume"`
}

func NewOHLCVTimeSeries(token int) OHLCVTimeSeries {
	return OHLCVTimeSeries{
		Token:     token,
		Timestamp: []time.Time{},
		Open:      []float64{},
		High:      []float64{},
		Low:       []float64{},
		Close:     []float64{},
		Volume:    []int{},
	}
}

func NewOHLCVTimeSeriesWithBarCount(token int, barCount int) OHLCVTimeSeries {
	return OHLCVTimeSeries{
		Token:     token,
		Timestamp: make([]time.Time, barCount),
		Open:      make([]float64, barCount),
		High:      make([]float64, barCount),
		Low:       make([]float64, barCount),
		Close:     make([]float64, barCount),
		Volume:    make([]int, barCount),
	}
}

type AlphaSignal struct {
	Token                int
	ExchangeCode         int
	Exchange             string
	Interval             string
	LastBarStartDuration time.Time
	Strategy             string
	Signal               string
	SignalGenerationTime time.Time
	Price                float64
	O                    float64
	H                    float64
	L                    float64
	C                    float64
	V                    int
	AlphaSignalReason    string
	IsHistorical         bool
}

func NewAlphaSignal(token int, exchangeCode int, exchange, interval string, lastBarStartDuration time.Time,
	strategy, signal string, signalGenerationTime time.Time, price, o, h, l, c float64, v int, alphaSignalReason string,
	isHistorical bool) (*AlphaSignal, error) {
	switch signal {
	case "LONG", "SHORT", "LONG_EXIT", "SHORT_EXIT":
	default:
		return nil, ErrInvalidStrategy
	}

	return &AlphaSignal{
		Token:                token,
		ExchangeCode:         exchangeCode,
		Exchange:             exchange,
		Interval:             interval,
		LastBarStartDuration: lastBarStartDuration,
		Strategy:             strategy,
		Signal:               signal,
		SignalGenerationTime: signalGenerationTime,
		Price:                price,
		O:                    o,
		H:                    h,
		L:                    l,
		C:                    c,
		V:                    v,
		AlphaSignalReason:    alphaSignalReason,
	}, nil
}
