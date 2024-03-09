package handler

import (
	"github.com/sudeepbatra/alpha-hft/broker/smartapi"
	"github.com/sudeepbatra/alpha-hft/data"
	"github.com/sudeepbatra/alpha-hft/logger"
)

const (
	TickDataTable       = "tick_data"
	BulkInsertThreshold = 1000
)

type TickDataStorage struct {
	ticksData []smartapi.TickParsedData
	done      chan struct{}
}

func NewTickDataStorage() *TickDataStorage {
	logger.Log.Debug().Str("handler", "tick_data_storage").Msg("New TickDataStorage Created!")

	return &TickDataStorage{
		ticksData: make([]smartapi.TickParsedData, 0, BulkInsertThreshold),
		done:      make(chan struct{}),
	}
}

func (t *TickDataStorage) FlushBuffer() {
	if len(t.ticksData) > 0 {
		err := data.BulkInsert(TickDataTable, t.ticksData)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error while Bulk inserting Tick Data to the DB")
		}

		logger.Log.Debug().Interface("tickData", t.ticksData).Msg("Bulk Stored Tick Data to the DB!")

		t.ticksData = t.ticksData[:0]
	}
}

func (t *TickDataStorage) StoreTickData() {
	parsedTickData := smartapi.SmartApiDataManager.Subscribe()

	logger.Log.Debug().Str("handler", "tick_data_processor").Msg("TickDataProcessor started!")

	for {
		select {
		case tickData, ok := <-parsedTickData:
			if !ok {
				t.FlushBuffer()
				close(t.done)

				return
			}

			// tickData.ExchangeTimestamp = tickData.ExchangeTimestamp
			// tickData.LastTradedTimestamp = tickData.LastTradedTimestamp / 1000

			logger.Log.Info().Interface("tickData", tickData).Msg("Storing Tick Data to the buffer...")

			t.ticksData = append(t.ticksData, *tickData)

			if len(t.ticksData) >= BulkInsertThreshold {
				t.FlushBuffer()
			}

		case <-t.done:
			return
		}
	}
}

func (t *TickDataStorage) Stop() {
	close(t.done)
	<-t.done
}
