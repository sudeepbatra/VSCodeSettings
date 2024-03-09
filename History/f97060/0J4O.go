package handler

import (
	"github.com/sudeepbatra/alpha-hft/broker/smartapi"
	"github.com/sudeepbatra/alpha-hft/data"
	"github.com/sudeepbatra/alpha-hft/logger"
)

const (
	TickDataTable       = "tick_data"
	BulkInsertThreshold = 1
)

type TickDataStorage struct {
	ticksData []smartapi.TickParsedData
	done      chan struct{}
}

func NewTickDataStorage() *TickDataStorage {
	logger.Log.Info().Str("handler", "tick_data_storage").Msg("New TickDataStorage Created!")

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

		logger.Log.Trace().Interface("tickData", t.ticksData).Msg("Bulk Stored Tick Data to the DB!")

		t.ticksData = t.ticksData[:0]
	}
}

func (t *TickDataStorage) StoreTickData() {
	logger.Log.Info().Str("handler", "tick_data_processor").Msg("StoreTickData started! Subscribing to tick channel...")

	excludedFields := map[string]bool{
		"Best5Data":     true,
		"Best5BuyData":  true,
		"Best5SellData": true,
	}

	parsedTickData := smartapi.SmartApiDataManager.Subscribe()

	for {
		select {
		case tickData, ok := <-parsedTickData:
			if !ok {
				t.FlushBuffer()
				close(t.done)

				return
			}

			modifiedTickData := ModifiedTickParsedData{}
			copyStructExcludingFields(tickData, &modifiedTickData, excludedFields)

			modifiedTickData = data.CopyStructExcludingFields()

			logger.Log.Trace().Interface("tickData", tickData).Msg("Storing Tick Data to the buffer...")

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
