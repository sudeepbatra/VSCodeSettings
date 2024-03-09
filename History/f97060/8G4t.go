package handler

import (
	"github.com/sudeepbatra/alpha-hft/broker/smartapi"
	"github.com/sudeepbatra/alpha-hft/data"
	"github.com/sudeepbatra/alpha-hft/logger"
)

const (
	TickDataTable       = "tick_data"
	BulkInsertThreshold = 2 // Adjust based on performance testing
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
		// Create a new slice to avoid modifying the original buffer
		dataToInsert := make([]interface{}, len(t.ticksData))
		for i, tickData := range t.ticksData {
			dataToInsert[i] = tickData
		}

		err := data.BulkInsert(TickDataTable, dataToInsert)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error while Bulk inserting Tick Data to the DB")
		}

		logger.Log.Debug().Interface("tickData", dataToInsert).Msg("Bulk Stored Tick Data to the DB!")

		// Clear the buffer without modifying the original slice
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
				// Channel closed, exit the loop
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
			// Handle any cleanup if needed
			return
		}
	}
}

func (t *TickDataStorage) Stop() {
	close(t.done)
	<-t.done
}
