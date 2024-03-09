package handler

import (
	"github.com/sudeepbatra/alpha-hft/broker/smartapi"
	"github.com/sudeepbatra/alpha-hft/data"
	"github.com/sudeepbatra/alpha-hft/logger"
)

const (
	TickDataTable       = "tick_data"
	BulkInsertThreshold = 2
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
		dataToInsert := make([]interface{}, len(t.ticksData))
		for i, tickData := range t.ticksData {
			dataToInsert[i] = tickData
		}

		err := data.BulkInsert(TickDataTable, dataToInsert)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error while Bulk inserting Tick Data to the DB")
		}

		logger.Log.Debug().Interface("tickData", dataToInsert).Msg("Bulk Stored Tick Data to the DB!")

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
