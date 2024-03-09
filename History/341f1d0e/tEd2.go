package handler

import (
	"sync"

	"github.com/sudeepbatra/alpha-hft/common"
	"github.com/sudeepbatra/alpha-hft/logger"
)

type AlphaSignalDistribution struct {
	mu          sync.Mutex
	subscribers []chan common.AlphaSignal
	BufferSize  int
}

var AlphaSignalManager = &AlphaSignalDistribution{
	BufferSize: 1000,
}

func (a *AlphaSignalDistribution) Subscribe() chan common.AlphaSignal {
	a.mu.Lock()
	defer a.mu.Unlock()

	logger.Log.Debug().Str("handler", "AlphaSignalDistribution").Msg("subscribing for alpha signals")

	ch := make(chan common.AlphaSignal, a.BufferSize)
	a.subscribers = append(a.subscribers, ch)

	return ch
}

func (a *AlphaSignalDistribution) Unsubscribe(ch chan common.AlphaSignal) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for i, subscriber := range a.subscribers {
		if subscriber == ch {
			logger.Log.Debug().Str("handler", "AlphaSignalDistribution").Msg("closing one of the alpha signal subscriber")
			close(ch)

			a.subscribers = append(a.subscribers[:i], a.subscribers[i+1:]...)

			break
		}
	}
}

func (a *AlphaSignalDistribution) PushAlphaSignalsForDistribution(alphaSignal common.AlphaSignal) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, subscriber := range a.subscribers {
		select {
		case subscriber <- alphaSignal:
		default:
			logger.Log.Warn().
				Str("handler", "alphasignaldistribution").
				Interface("alphasignal", alphaSignal).
				Interface("subscribers", a.subscribers).
				Int("current subscriber size", len(subscriber)).
				Msg("unable to push alpha signal for distribution")
		}
	}
}
