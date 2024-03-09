package handler

import "sync"

type CandlestickDistribution struct {
	mu          sync.Mutex
	subscribers []chan Candlestick
}

var CandlestickDataManager = &CandlestickDistribution{}

func (c *CandlestickDistribution) Subscribe() chan Candlestick {
	c.mu.Lock()

}
