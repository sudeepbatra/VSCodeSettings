package ta

import (
	"fmt"
	"strconv"
	"time"

	"github.com/markcheno/go-talib"
	"github.com/sudeepbatra/alpha-hft/broker/smartapi"
	"github.com/sudeepbatra/alpha-hft/common"
	"github.com/sudeepbatra/alpha-hft/data"
	"github.com/sudeepbatra/alpha-hft/handler"
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/indicators"
)

const (
	minCandlesticksRequiredForAlphaSignalGen = 80
	traceRule                                = true
)

type AlphaSignalGenerator struct {
	historicalOHLCVTimeSeriesData map[string]map[int]map[string]*common.OHLCVTimeSeries
}

func NewAlphaSignalGenerator(
	historicalOHLCVTimeSeriesData map[string]map[int]map[string]*common.OHLCVTimeSeries,
) *AlphaSignalGenerator {
	return &AlphaSignalGenerator{
		historicalOHLCVTimeSeriesData: historicalOHLCVTimeSeriesData,
	}
}

func (a *AlphaSignalGenerator) AddToken(interval string, token int, exchgCode string, data *common.OHLCVTimeSeries) {
	if _, exists := a.historicalOHLCVTimeSeriesData[interval]; !exists {
		a.historicalOHLCVTimeSeriesData[interval] = make(map[int]map[string]*common.OHLCVTimeSeries)
	}

	a.historicalOHLCVTimeSeriesData[interval][token][exchgCode] = data
}

func (a *AlphaSignalGenerator) GenerateHistoricalAlphaSignals() {
	logger.Log.Info().Msg("Starting generating Historical Alpha signals...")

	_, tokenToSymbolMap, err := data.GetInstrumentsSymbolAndTokenMapping()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("error in getting symbol to token and token to symbol mapping")
		return
	}

	for interval, intervalData := range a.historicalOHLCVTimeSeriesData {
		for token, tokenData := range intervalData {
			for exchange, ohlcvTimeSeries := range tokenData {
				instrumentTokenKey := smartapi.InstrumentTokenKey{Token: strconv.Itoa(token), Exchange: exchange}
				instrument := tokenToSymbolMap[instrumentTokenKey]

				if len(ohlcvTimeSeries.Close) < minCandlesticksRequiredForAlphaSignalGen {
					logger.Log.Debug().
						Str("interval", interval).
						Str("symbol", instrument.Symbol).
						Str("exchange", instrument.Exchange).
						Int("token", token).
						Msg("Not enough data to generate Alpha signals")

					continue
				}

				logger.Log.Debug().Msg("olhcvtimeSeries: " + fmt.Sprintf("%+v", ohlcvTimeSeries))

				rsi14, sar, ichimoku, adx, plusDI, minusDI, sma20, ema10, ema20, ema60 := calculateIndicators(ohlcvTimeSeries)

				strategies := createStrategies(ohlcvTimeSeries, sar, rsi14, ichimoku, adx, plusDI, minusDI,
					sma20, ema10, ema20, ema60)

				for _, strategy := range strategies {
					enter := strategy.ShouldEnter(len(ohlcvTimeSeries.Close) - 1)

					if enter {
						logger.Log.Info().
							Str("interval", interval).
							Str("symbol", instrument.Symbol).
							Str("exchange", instrument.Exchange).
							Int("token", token).
							Str("strategy", strategy.GetName()).
							Str("strategy description", strategy.GetDescription()).
							Msg("Enter position based on strategy")

						if adxStrategy, ok := strategy.(*ADXBasicLongStrategy); ok {
							logger.Log.Info().
								Interface("adx", adxStrategy.Parameters["adx"]).
								Interface("plusDI", adxStrategy.Parameters["plusDI"]).
								Interface("minusDI", adxStrategy.Parameters["minusDI"]).
								Msg("ADX strategy parameters")
						}

						signal := "LONG"
						if strategy.GetStrategyType() == "SHORT" {
							signal = "SHORT"
						}

						lastIndex := len(ohlcvTimeSeries.Close) - 1

						entryAlphaSinal, err := common.NewAlphaSignal(
							strconv.Itoa(token),
							instrument.Symbol,
							smartapi.ExchangeCodeTypes[exchange],
							exchange,
							interval,
							ohlcvTimeSeries.Timestamp[lastIndex],
							strategy.GetName(),
							strategy.IsStrategyLive(),
							signal,
							time.Now(),
							ohlcvTimeSeries.Close[lastIndex],
							ohlcvTimeSeries.Open[lastIndex],
							ohlcvTimeSeries.High[lastIndex],
							ohlcvTimeSeries.Low[lastIndex],
							ohlcvTimeSeries.Close[lastIndex],
							ohlcvTimeSeries.Volume[lastIndex],
							"",
							true)
						if err != nil {
							logger.Log.Error().Err(err).Msg("Error creating AlphaSignal")
						}

						handler.AlphaSignalManager.PushAlphaSignalsForDistribution(*entryAlphaSinal)
						time.Sleep(time.Millisecond * 10)
					}

					exit := strategy.ShouldExit(len(ohlcvTimeSeries.Close) - 1)

					if exit {
						logger.Log.Info().
							Str("interval", interval).
							Str("symbol", instrument.Symbol).
							Str("exchange", instrument.Exchange).
							Int("token", token).
							Msg("Exit position based on strategy: " + strategy.GetName())

						lastIndex := len(ohlcvTimeSeries.Close) - 1

						signal := "LONG_EXIT"
						if strategy.GetStrategyType() == "SHORT" {
							signal = "SHORT_EXIT"
						} else if strategy.GetStrategyType() == "BOTH" {
							signal = "SHORT"
						}

						// TODO: Hardcoded exchange code. Needs to be fixed
						exchangeCode := 1

						exitAlphaSinal, err := common.NewAlphaSignal(
							strconv.Itoa(token),
							instrument.Symbol,
							exchangeCode,
							smartapi.CodeExchangeTypes[exchangeCode],
							interval,
							ohlcvTimeSeries.Timestamp[lastIndex],
							strategy.GetName(),
							strategy.IsStrategyLive(),
							signal,
							time.Now(),
							ohlcvTimeSeries.Close[lastIndex],
							ohlcvTimeSeries.Open[lastIndex],
							ohlcvTimeSeries.High[lastIndex],
							ohlcvTimeSeries.Low[lastIndex],
							ohlcvTimeSeries.Close[lastIndex],
							ohlcvTimeSeries.Volume[lastIndex],
							"",
							true)
						if err != nil {
							logger.Log.Error().Err(err).Msg("Error creating AlphaSignal")
						}

						handler.AlphaSignalManager.PushAlphaSignalsForDistribution(*exitAlphaSinal)
						time.Sleep(time.Millisecond * 10)
					}
				}
			}
		}
	}

	logger.Log.Info().Msg("Finished generating Historical Alpha signals!")
}

func createStrategies(ohlcvTimeSeries *common.OHLCVTimeSeries, sar []float64, rsi14 []float64,
	ichimoku *indicators.IchimokuCloud, adx []float64, plusDI []float64, minusDI []float64, sma20 []float64,
	ema10 []float64, ema20 []float64, ema60 []float64,
) []Strategy {
	strategies := []Strategy{
		NewPSARStrategy(ohlcvTimeSeries.Close, sar),
		NewRSIAboveStrategy(rsi14, 70, 30),
		NewIchimokuChikouCrossoverStrategy(ichimoku.Chikou, ohlcvTimeSeries.Close),
		NewIchimokuChikouHighLowCrossoverStrategy(ichimoku.Chikou, ohlcvTimeSeries.High, ohlcvTimeSeries.Low),
		NewIchimokuChikouHighLowNBarsCrossoverStrategy(ichimoku.Chikou, ohlcvTimeSeries.High, ohlcvTimeSeries.Low),
		NewADXBasicLongStrategy(adx, plusDI, minusDI, 20),
		NewADXEMIRSILongStrategy(ohlcvTimeSeries.Close, adx, plusDI, minusDI, sma20, ema10,
			ema20, ema60, rsi14, sar, 20, 55),
		NewADXEMIPSARLongStrategy(ohlcvTimeSeries.Close, plusDI, minusDI, ema10, ema20, sar),
		NewIchimokuEmaLongStrategy(ohlcvTimeSeries.Close, ema10, ema20, ichimoku.Chikou, sar),
		NewIchimokuLowHighEmaLongStrategy(ohlcvTimeSeries.Low, ohlcvTimeSeries.High, ohlcvTimeSeries.Close,
			ema10, ema20, ichimoku.Chikou, sar),
	}

	return strategies
}

func calculateIndicators(ohlcvTimeSeries *common.OHLCVTimeSeries) ([]float64, []float64, *indicators.IchimokuCloud,
	[]float64, []float64, []float64, []float64, []float64, []float64, []float64,
) {
	rsi14 := talib.Rsi(ohlcvTimeSeries.Close, 14)

	sar := talib.Sar(ohlcvTimeSeries.High, ohlcvTimeSeries.Low, 0.02, 0.2)

	ichimoku := indicators.CalculateIchimokuCloud(ohlcvTimeSeries.High, ohlcvTimeSeries.Low, ohlcvTimeSeries.Close, 9, 26, 52)

	adx := talib.Adx(ohlcvTimeSeries.High, ohlcvTimeSeries.Low, ohlcvTimeSeries.Close, 14)

	plusDI := talib.PlusDI(ohlcvTimeSeries.High, ohlcvTimeSeries.Low, ohlcvTimeSeries.Close, 14)

	minusDI := talib.MinusDI(ohlcvTimeSeries.High, ohlcvTimeSeries.Low, ohlcvTimeSeries.Close, 14)

	sma20 := talib.Sma(ohlcvTimeSeries.Close, 20)

	ema10 := talib.Ema(ohlcvTimeSeries.Close, 10)

	ema20 := talib.Ema(ohlcvTimeSeries.Close, 20)

	ema60 := talib.Ema(ohlcvTimeSeries.Close, 60)

	logger.Log.Trace().
		Interface("rsi14", rsi14).
		Interface("sar", sar).
		Interface("ichimoku", ichimoku).
		Interface("adx", adx).
		Interface("plusDI", plusDI).
		Interface("minusDI", minusDI).
		Interface("sma20", sma20).
		Interface("ema10", ema10).
		Interface("ema20", ema20).
		Interface("ema60", ema60).
		Msg("Indicators calculated")

	return rsi14, sar, ichimoku, adx, plusDI, minusDI, sma20, ema10, ema20, ema60
}

func (a *AlphaSignalGenerator) GenerateAlphaSignals() {
	logger.Log.Info().Msg("Starting generating Live Alpha signals...")

	_, tokenToSymbolMap, err := data.GetInstrumentsSymbolAndTokenMapping()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("error in getting symbol to token and token to symbol mapping")
		return
	}

	completedCandles := handler.CandlestickDataManager.Subscribe()

	logger.Log.Info().Msg("Subscribed to completed candles channel")

	for completedCandle := range completedCandles {
		logger.Log.Trace().
			Interface("completedCandle", completedCandle).
			Msg("Received completed candle in GenerateAlphaSignals")

		token, _ := strconv.Atoi(completedCandle.Token)
		exchange := smartapi.CodeExchangeTypes[completedCandle.Exchange]

		instrumentTokenKey := smartapi.InstrumentTokenKey{Token: strconv.Itoa(token), Exchange: exchange}
		instrument := tokenToSymbolMap[instrumentTokenKey]
		interval := completedCandle.Duration

		if ohlcvTimeSeries, ok := a.historicalOHLCVTimeSeriesData[interval][token][exchange]; ok {
			logger.Log.Debug().Msg("olhcvtimeSeries: " + fmt.Sprintf("%+v", ohlcvTimeSeries))

			verifyOhlcvTimeSeriesSize := common.CheckSameSize(ohlcvTimeSeries.Timestamp,
				ohlcvTimeSeries.Open, ohlcvTimeSeries.High, ohlcvTimeSeries.Low,
				ohlcvTimeSeries.Close, ohlcvTimeSeries.Volume)

			if !verifyOhlcvTimeSeriesSize {
				logger.Log.
					Error().
					Str("interval", interval).
					Int("token", token).
					Str("exchange", exchange).
					Msg("OHLCV TimeSeries size is not same. Skipping Alpha signal generation")

				continue
			}

			if len(ohlcvTimeSeries.Open) < minCandlesticksRequiredForAlphaSignalGen {
				logger.Log.
					Warn().
					Str("interval", interval).
					Int("token", token).
					Str("exchange", exchange).
					Int("candlesticks", len(ohlcvTimeSeries.Open)).
					Msg("Not enough candlesticks. Skipping Alpha signal generation")

				continue
			}

			ohlcvTimeSeries.Timestamp = append(ohlcvTimeSeries.Timestamp, completedCandle.Timestamp)
			ohlcvTimeSeries.Open = append(ohlcvTimeSeries.Open, completedCandle.Open)
			ohlcvTimeSeries.High = append(ohlcvTimeSeries.High, completedCandle.High)
			ohlcvTimeSeries.Low = append(ohlcvTimeSeries.Low, completedCandle.Low)
			ohlcvTimeSeries.Close = append(ohlcvTimeSeries.Close, completedCandle.Close)
			ohlcvTimeSeries.Volume = append(ohlcvTimeSeries.Volume, int(completedCandle.Volume))

			logger.Log.Debug().
				Interface("ohlcvTimeSeries", ohlcvTimeSeries).
				Msg("Updated ohlcvTimeSeries in GenerateAlphaSignals before generating Alpha signals")

			rsi14, sar, ichimoku, adx, plusDI, minusDI, sma20, ema10, ema20, ema60 := calculateIndicators(ohlcvTimeSeries)

			strategies := createStrategies(ohlcvTimeSeries, sar, rsi14, ichimoku, adx, plusDI, minusDI,
				sma20, ema10, ema20, ema60)

			for _, strategy := range strategies {
				enter := strategy.ShouldEnter(len(ohlcvTimeSeries.Close) - 1)

				if enter {
					logger.Log.Info().
						Str("interval", interval).
						Str("symbol", instrument.Symbol).
						Str("exchange", instrument.Exchange).
						Int("token", token).
						Str("strategy", strategy.GetName()).
						Str("strategy description", strategy.GetDescription()).
						Msg("Enter position based on strategy")

					signal := "LONG"
					if strategy.GetStrategyType() == "SHORT" {
						signal = "SHORT"
					}

					lastIndex := len(ohlcvTimeSeries.Close) - 1

					entryAlphaSinal, err := common.NewAlphaSignal(
						strconv.Itoa(token),
						instrument.Symbol,
						completedCandle.Exchange,
						smartapi.CodeExchangeTypes[completedCandle.Exchange],
						interval,
						completedCandle.Timestamp,
						strategy.GetName(),
						strategy.IsStrategyLive(),
						signal,
						time.Now(),
						ohlcvTimeSeries.Close[lastIndex],
						ohlcvTimeSeries.Open[lastIndex],
						ohlcvTimeSeries.High[lastIndex],
						ohlcvTimeSeries.Low[lastIndex],
						ohlcvTimeSeries.Close[lastIndex],
						ohlcvTimeSeries.Volume[lastIndex],
						"",
						false)
					if err != nil {
						logger.Log.Error().Err(err).Msg("Error creating AlphaSignal")
					}

					handler.AlphaSignalManager.PushAlphaSignalsForDistribution(*entryAlphaSinal)
				}

				exit := strategy.ShouldExit(len(ohlcvTimeSeries.Close) - 1)

				if exit {
					logger.Log.Info().
						Str("interval", interval).
						Str("symbol", instrument.Symbol).
						Str("exchange", instrument.Exchange).
						Int("token", token).
						Msg("Exit position based on strategy: " + strategy.GetName())

					lastIndex := len(ohlcvTimeSeries.Close) - 1

					signal := "LONG_EXIT"
					if strategy.GetStrategyType() == "SHORT" {
						signal = "SHORT_EXIT"
					} else if strategy.GetStrategyType() == "BOTH" {
						signal = "SHORT"
					}

					exitAlphaSinal, err := common.NewAlphaSignal(
						strconv.Itoa(token),
						instrument.Symbol,
						completedCandle.Exchange,
						smartapi.CodeExchangeTypes[completedCandle.Exchange],
						interval,
						completedCandle.Timestamp,
						strategy.GetName(),
						strategy.IsStrategyLive(),
						signal,
						time.Now(),
						ohlcvTimeSeries.Close[lastIndex],
						ohlcvTimeSeries.Open[lastIndex],
						ohlcvTimeSeries.High[lastIndex],
						ohlcvTimeSeries.Low[lastIndex],
						ohlcvTimeSeries.Close[lastIndex],
						ohlcvTimeSeries.Volume[lastIndex],
						"",
						false)
					if err != nil {
						logger.Log.Error().Err(err).Msg("Error creating AlphaSignal")
					}

					handler.AlphaSignalManager.PushAlphaSignalsForDistribution(*exitAlphaSinal)
				}
			}
		} else {
			logger.Log.Warn().
				Str("interval", interval).
				Str("symbol", instrument.Symbol).
				Str("exchange", instrument.Exchange).
				Int("token", token).
				Int("minimum candlesticks required for alpha signal generation", minCandlesticksRequiredForAlphaSignalGen).
				Msg("Skipping instrument for alpha signal generation as minimum candlesticks required is not met")
		}
	}

	logger.Log.Info().Msg("Finished generating Live Alpha signals!")
}
