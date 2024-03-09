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

func (c *CandlestickDistribution) Unsubscribe(ch chan Candlestick) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, subscriber := range c.subscribers {
		if subscriber == ch {
			logger.Log.Debug().Str("handler", "CandlestickDistribution").Msg("closing one of the candlestick data subscribers")
			close(ch)

			c.subscribers = append(c.subscribers[:i], c.subscribers[i+1:]...)

			break
		}
	}
}

func (c *CandlestickDistribution) PushCandlestickForDistribution(candlestick Candlestick) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, subscriber := range c.subscribers {
		select {
		case subscriber <- candlestick:
		default:
			logger.Log.Info().
				Str("handler", "CandlestickDistribution").
				Interface("subscribers", subscriber).
				Int("Current Size", len(subscriber)).
				Msg("unable to push candlestick data for distribution")
		}
	}

	for _, subscriber := range dd.subscribers {
		select {
		case subscriber <- data:
		default:
			logger.Log.Info().
				Str("broker", "SmartApi").
				Interface("data", data).
				Interface("subscribers", subscriber).
				Int("Current Size", len(subscriber)).
				Msg("unable to push message for distribution")
		}
	}
}
