package handler

import (
	"sync"

	"github.com/sudeepbatra/alpha-hft/logger"
)

type CandlestickDistribution struct {
	mu          sync.Mutex
	subscribers []chan Candlestick
}

var CandlestickDataManager = &CandlestickDistribution{}

func (c *CandlestickDistribution) Subscribe() chan Candlestick {
	c.mu.Lock()
	defer c.mu.Unlock()

	logger.Log.Debug().Str("handler", "CandlestickDistribution").Msg("subscribing for live candlestick data")

	ch := make(chan Candlestick, 10)
	c.subscribers = append(c.subscribers, ch)

	return ch
}
