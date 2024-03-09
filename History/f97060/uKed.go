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
	ticksData []smartapi.TickParsedDataStorage
	done      chan struct{}
}

func NewTickDataStorage() *TickDataStorage {
	logger.Log.Info().Str("handler", "tick_data_storage").Msg("New TickDataStorage Created!")

	return &TickDataStorage{
		ticksData: make([]smartapi.TickParsedDataStorage, 0, BulkInsertThreshold),
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

	parsedTickData := smartapi.SmartApiDataManager.Subscribe()

	for {
		select {
		case tickData, ok := <-parsedTickData:
			if !ok {
				t.FlushBuffer()
				close(t.done)

				return
			}

			tickParsedDataStorage := smartapi.TickParsedDataStorage{
				SubscriptionMode:             tickData.SubscriptionMode,
				ExchangeType:                 tickData.ExchangeType,
				Token:                        tickData.Token,
				SequenceNumber:               tickData.SequenceNumber,
				ExchangeTimestamp:            tickData.ExchangeTimestamp,
				LastTradedPrice:              tickData.LastTradedPrice,
				SubscriptionModeVal:          tickData.SubscriptionModeVal,
				LastTradedQuantity:           tickData.LastTradedQuantity,
				AverageTradedPrice:           tickData.AverageTradedPrice,
				VolumeTradeForTheDay:         tickData.VolumeTradeForTheDay,
				TotalBuyQuantity:             tickData.TotalBuyQuantity,
				TotalSellQuantity:            tickData.TotalSellQuantity,
				OpenPriceOfTheDay:            tickData.OpenPriceOfTheDay,
				HighPriceOfTheDay:            tickData.HighPriceOfTheDay,
				LowPriceOfTheDay:             tickData.LowPriceOfTheDay,
				ClosedPrice:                  tickData.ClosedPrice,
				LastTradedTimestamp:          tickData.LastTradedTimestamp,
				OpenInterest:                 tickData.OpenInterest,
				OpenInterestChangePercentage: tickData.OpenInterestChangePercentage,
				UpperCircuitLimit:            tickData.UpperCircuitLimit,
				LowerCircuitLimit:            tickData.LowerCircuitLimit,
				Week52HighPrice:              tickData.Week52HighPrice,
				Week52LowPrice:               tickData.Week52LowPrice,
				LastTradedPriceFloat:         tickData.LastTradedPriceFloat,
			}

			logger.Log.Trace().
				Interface("tickParsedDataStorage", tickParsedDataStorage).
				Msg("Storing Tick Data to the buffer...")

			t.ticksData = append(t.ticksData, tickParsedDataStorage)

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
