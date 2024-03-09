package smartapi

import (
	"sync"

	"github.com/sudeepbatra/alpha-hft/logger"
)

type DataDistribution struct {
	mu          sync.Mutex
	subscribers []chan *TickParsedData
	BufferSize  int
}

var SmartApiDataManager = &DataDistribution{
	BufferSize: 1000,
}

func (dd *DataDistribution) Subscribe() chan *TickParsedData {
	dd.mu.Lock()
	defer dd.mu.Unlock()

	logger.Log.Info().Str("broker", "smart_api").Msg("Subscribe received a subscription for live ticker market data")

	ch := make(chan *TickParsedData, BufferSize)
	dd.subscribers = append(dd.subscribers, ch)

	return ch
}

func (dd *DataDistribution) Unsubscribe(ch chan *TickParsedData) {
	dd.mu.Lock()
	defer dd.mu.Unlock()

	for i, subscriber := range dd.subscribers {
		if subscriber == ch {
			logger.Log.Info().Str("broker", "smart_api").Msg("closing one of the subscriber for live market data")
			close(ch)

			dd.subscribers = append(dd.subscribers[:i], dd.subscribers[i+1:]...)

			break
		}
	}
}

func (dd *DataDistribution) PushMessageForDistribution(data *TickParsedData) {
	dd.mu.Lock()
	defer dd.mu.Unlock()

	for _, subscriber := range dd.subscribers {
		select {
		case subscriber <- data:
		default:
			logger.Log.Info().Str("broker", "SmartApi").Msg("unable to push message for distribution")
		}
	}
}
