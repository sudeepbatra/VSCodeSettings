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
	buffer []*smartapi.TickParsedData
	done   chan struct{}
}

func NewTickDataStorage() *TickDataStorage {
	logger.Log.Debug().Str("handler", "tick_data_storage").Msg("New TickDataStorage Created!")

	return &TickDataStorage{
		buffer: make([]*smartapi.TickParsedData, 0, BulkInsertThreshold),
		done:   make(chan struct{}),
	}
}

func (t *TickDataStorage) FlushBuffer() {
	if len(t.buffer) > 0 {
		// Create a new slice to avoid modifying the original buffer
		dataToInsert := make([]interface{}, len(t.buffer))
		for i, tickData := range t.buffer {
			dataToInsert[i] = *tickData
		}

		err := data.BulkInsert(TickDataTable, dataToInsert)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error while Bulk inserting Tick Data to the DB")
		}

		logger.Log.Debug().Interface("tickData", dataToInsert).Msg("Bulk Stored Tick Data to the DB!")

		// Clear the buffer without modifying the original slice
		t.buffer = t.buffer[:0]
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

			t.buffer = append(t.buffer, tickData)

			if len(t.buffer) >= BulkInsertThreshold {
				t.FlushBuffer()
			}

		case <-t.done:
			// Handle any cleanup if needed
			return
		}
	}
}

// Stop processing and wait for any remaining data to be stored
func (t *TickDataStorage) Stop() {
	close(t.done)
	// Wait for the StoreTickData goroutine to finish
	<-t.done
}
