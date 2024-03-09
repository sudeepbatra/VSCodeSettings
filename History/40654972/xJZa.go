package handler

import (
	"sort"
	"sync"
	"time"

	"github.com/sudeepbatra/alpha-hft/broker/smartapi"
	"github.com/sudeepbatra/alpha-hft/logger"
)

type Candlestick struct {
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    int64
	Token     string
	Exchange  int
	Duration  string
}

type CandlestickInterval struct {
	Duration      time.Duration
	Name          string
	IntervalStart time.Time
	Candlesticks  map[time.Time]Candlestick
}

type CandlestickHandler struct {
	mu              sync.Mutex
	intervals       []CandlestickInterval
	marketOpenTime  time.Time
	marketCloseTime time.Time
}

type TokenCandlestickHandlers struct {
	mu       sync.Mutex
	handlers map[string]*CandlestickHandler
}

func NewTokenCandlestickHandlers() *TokenCandlestickHandlers {
	logger.Log.Debug().Str("parser", "candlestick_from_tick").Msg("New token Candlestick Handler Created!")

	return &TokenCandlestickHandlers{
		handlers: make(map[string]*CandlestickHandler),
	}
}

var handlers = NewTokenCandlestickHandlers()

func NewCandlestickHandler() *CandlestickHandler {
	handler := &CandlestickHandler{
		intervals: []CandlestickInterval{
			{1 * time.Minute, "one_minute", time.Time{}, make(map[time.Time]Candlestick)},
			{2 * time.Minute, "two_minute", time.Time{}, make(map[time.Time]Candlestick)},
			{5 * time.Minute, "five_minute", time.Time{}, make(map[time.Time]Candlestick)},
			{10 * time.Minute, "ten_minute", time.Time{}, make(map[time.Time]Candlestick)},
			{15 * time.Minute, "fifteen_minute", time.Time{}, make(map[time.Time]Candlestick)},
			{45 * time.Minute, "fortyfive_minute", time.Time{}, make(map[time.Time]Candlestick)},
			{1 * time.Hour, "one_hour", time.Time{}, make(map[time.Time]Candlestick)}},
	}
	handler.InitializeMarketOpenCloseTime()

	return handler
}

func (handler *CandlestickHandler) InitializeMarketOpenCloseTime() {
	currentDate := time.Now().Local()
	year, month, day := currentDate.Date()
	handler.marketOpenTime = time.Date(year, month, day, 9, 15, 0, 0, time.Local)
	handler.marketCloseTime = time.Date(year, month, day, 15, 30, 0, 0, time.Local)
}

func (handler *CandlestickHandler) ProcessTicks() {
	logger.Log.Info().Str("parser", "candlestick_from_tick").Msg("ProcessTicks started! Subscribing to tick channel...")

	parsedTickData := smartapi.SmartApiDataManager.Subscribe()

	for data := range parsedTickData {
		if data == nil {
			logger.Log.Error().
				Str("parser", "candlestick_from_tick").
				Msg("channel is closed. returning from ProcessTicks!")

			break
		}

		handler := handlers.getOrCreateCandlestickTokenHandler(data.Token)
		go handler.processTick(*data)
	}
}

func (handlers *TokenCandlestickHandlers) getOrCreateCandlestickTokenHandler(token string) *CandlestickHandler {
	handlers.mu.Lock()
	defer handlers.mu.Unlock()

	handler, exists := handlers.handlers[token]
	if !exists {
		logger.Log.Debug().Str("parser", "candlestick_from_tick").Msg("New handler for token: " + token)

		handler = NewCandlestickHandler()
		handlers.handlers[token] = handler
	}

	return handler
}

func (handler *CandlestickHandler) processTick(data smartapi.TickParsedData) {
	token := data.Token
	exchange := smartapi.CodeExchangeTypes[int(data.ExchangeType)]
	logger.Log.Trace().
		Str("parser", "candlestick_from_tick").
		Str("token", token).
		Str("exchange", exchange).
		Interface("tick", data).
		Str("timestamp", data.ExchangeTimestamp.Format("2006-01-02 15:04:05")).
		Msg("Tick received: ")

	timestamp := data.ExchangeTimestamp

	handler.mu.Lock()
	defer handler.mu.Unlock()

	if timestamp.Before(handler.marketOpenTime) || timestamp.After(handler.marketCloseTime) {
		logger.Log.
			Info().
			Str("parser", "candlestick_from_tick").
			Msg("Tick outside the market hours: " + timestamp.Format("2006-01-02 15:04:05"))

		return
	} else {
		logger.Log.Trace().
			Str("parser", "candlestick_from_tick").
			Interface("tick", data).
			Str("timestamp", data.ExchangeTimestamp.Format("2006-01-02 15:04:05")).
			Msg("Candlestick from tick: Processing the tick: ")
	}

	// Process each candlestick interval
	for i := range handler.intervals {
		interval := &handler.intervals[i]

		if interval.IntervalStart.IsZero() {
			interval.IntervalStart = handler.getIntervalStartTime(timestamp, interval.Duration)
			logger.Log.Info().
				Str("parser", "candlestick_from_tick").
				Interface("interval", interval).
				Msg("Interval Start time set to: " + interval.IntervalStart.Format("2006-01-02 15:04:05") + " which was zero for interval: " + interval.Name)
		}

		// Check if tick belongs to the current interval
		if timestamp.Before(interval.IntervalStart.Add(interval.Duration)) {
			logger.Log.Debug().
				Str("parser", "candlestick_from_tick").
				Msg("Updating Candlestick data for interval: " + interval.Name)
			handler.updateCandlesticks(interval.IntervalStart, interval.Candlesticks, data, interval)
		} else {
			logger.Log.Info().
				Str("parser", "candlestick_from_tick").
				Interface("Completed Candlestick", interval.Candlesticks[interval.IntervalStart]).
				Msg("Pushing the completed candle for distribution for interval: " + interval.Name)

			// If the tick belongs to the next interval, send the complete candlestick and set the next interval start time
			CandlestickDataManager.PushCandlestickForDistribution(interval.Candlesticks[interval.IntervalStart])

			logger.Log.Debug().
				Str("parser", "candlestick_from_tick").
				Interface("Candlesticks", interval.Candlesticks) /

				logger.Log.Debug().
					Str("parser", "candlestick_from_tick").
					Msg("Setting the Interval Start time for the next interval: "+interval.Name)

			interval.IntervalStart = handler.getIntervalStartTime(timestamp, interval.Duration)
			handler.printCandlesticks(interval.Candlesticks)
		}
	}
}

func (handler *CandlestickHandler) getIntervalStartTime(timestamp time.Time, duration time.Duration) time.Time {
	intervalCount := int(timestamp.Sub(handler.marketOpenTime) / duration)

	if intervalCount == 0 {
		return handler.marketOpenTime
	}

	return handler.marketOpenTime.Add(duration * time.Duration(intervalCount))
}

func (handler *CandlestickHandler) updateCandlesticks(intervalStart time.Time, candlestickMap map[time.Time]Candlestick, tick smartapi.TickParsedData, interval *CandlestickInterval) {
	candlestick, exists := candlestickMap[intervalStart]
	if !exists {
		candlestick = Candlestick{
			Timestamp: intervalStart,
			Open:      tick.LastTradedPriceFloat,
			High:      tick.LastTradedPriceFloat,
			Low:       tick.LastTradedPriceFloat,
			Close:     tick.LastTradedPriceFloat,
			Volume:    tick.LastTradedQuantity,
			Token:     tick.Token,
			Exchange:  int(tick.ExchangeType),
			Duration:  interval.Name,
		}
	}

	if tick.LastTradedPriceFloat > candlestick.High {
		candlestick.High = tick.LastTradedPriceFloat
	}

	if tick.LastTradedPriceFloat < candlestick.Low {
		candlestick.Low = tick.LastTradedPriceFloat
	}

	candlestick.Close = tick.LastTradedPriceFloat
	candlestick.Volume += tick.LastTradedQuantity

	candlestickMap[intervalStart] = candlestick
}

func (handler *CandlestickHandler) printCandlesticks(candlestickMap map[time.Time]Candlestick) {
	var keys []time.Time
	for k := range candlestickMap {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	logger.Log.Trace().Str("parser", "candlestick_from_tick").Msg("Printing Candlesticks")
	logger.Log.Trace().Str("parser", "candlestick_from_tick").Msg("-----------------------------------------------------")

	for _, key := range keys {
		candlestick := candlestickMap[key]
		logger.Log.Trace().
			Str("parser", "candlestick_from_tick").
			Msgf("Timestamp: %s, O: %.2f, H: %.2f, L: %.2f, C: %.2f, V: %d\n", key.Format(time.RFC3339), candlestick.Open, candlestick.High, candlestick.Low, candlestick.Close, candlestick.Volume)
	}

	logger.Log.Trace().Str("parser", "candlestick_from_tick").Msg("-----------------------------------------------------")
}
