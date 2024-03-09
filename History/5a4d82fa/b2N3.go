package controller

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/markcheno/go-talib"
	"github.com/sudeepbatra/alpha-hft/broker/smartapi"
	"github.com/sudeepbatra/alpha-hft/common"
	"github.com/sudeepbatra/alpha-hft/config"
	"github.com/sudeepbatra/alpha-hft/data"
	"github.com/sudeepbatra/alpha-hft/handler"
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/messaging"
	"github.com/sudeepbatra/alpha-hft/ta"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

const HistoricalPricesLimitRatePerSecond = 3

var nifty50Instruments = map[string]struct{}{
	"ADANIENT-EQ":   {},
	"ADANIPORTS-EQ": {},
	"APOLLOHOSP-EQ": {},
	"ASIANPAINT-EQ": {},
	"AXISBANK-EQ":   {},
	"BAJAJ-AUTO-EQ": {},
	"BAJFINANCE-EQ": {},
	"BAJAJFINSV-EQ": {},
	"BPCL-EQ":       {},
	"BHARTIARTL-EQ": {},
	"BRITANNIA-EQ":  {},
	"CIPLA-EQ":      {},
	"COALINDIA-EQ":  {},
	"DIVISLAB-EQ":   {},
	"DRREDDY-EQ":    {},
	"EICHERMOT-EQ":  {},
	"GRASIM-EQ":     {},
	"HCLTECH-EQ":    {},
	"HDFCBANK-EQ":   {},
	"HDFCLIFE-EQ":   {},
	"HEROMOTOCO-EQ": {},
	"HINDALCO-EQ":   {},
	"HINDUNILVR-EQ": {},
	"ICICIBANK-EQ":  {},
	"ITC-EQ":        {},
	"INDUSINDBK-EQ": {},
	"INFY-EQ":       {},
	"JSWSTEEL-EQ":   {},
	"KOTAKBANK-EQ":  {},
	"LTIM-EQ":       {},
	"LT-EQ":         {},
	"M&M-EQ":        {},
	"MARUTI-EQ":     {},
	"NTPC-EQ":       {},
	"NESTLEIND-EQ":  {},
	"ONGC-EQ":       {},
	"POWERGRID-EQ":  {},
	"RELIANCE-EQ":   {},
	"SBILIFE-EQ":    {},
	"SBIN-EQ":       {},
	"SUNPHARMA-EQ":  {},
	"TCS-EQ":        {},
	"TATACONSUM-EQ": {},
	"TATAMOTORS-EQ": {},
	"TATASTEEL-EQ":  {},
	"TECHM-EQ":      {},
	"TITAN-EQ":      {},
	"UPL-EQ":        {},
	"ULTRACEMCO-EQ": {},
	"WIPRO-EQ":      {},
}

var subsciptionInstruments = map[string]struct{}{
	"5PAISA-EQ":     {},
	"3MINDIA-EQ":    {},
	"ABBOTINDIA-EQ": {},
	"AEGISCHEM-EQ":  {},
	"APOLLOTYRE-EQ": {},
	"ASHOKLEY-EQ":   {},
	"AUROPHARMA-EQ": {},
	"BAJAJCON-EQ":   {},
	"BAJAJELEC-EQ":  {},
	"BATAINDIA-EQ":  {},
	"BERGEPAINT-EQ": {},
	"BHEL-EQ":       {},
	"BIKAJI-EQ":     {},
	"BIOCON-EQ":     {},
	"BLUEDART-EQ":   {},
	"BOSCHLTD-EQ":   {},
	"BPCL-EQ":       {},
	"BRITANNIA-EQ":  {},
	"CAMLINFINE-EQ": {},
	"CANFINHOME-EQ": {},
	"CARTRADE-EQ":   {},
	"CASTROLIND-EQ": {},
	"CEATLTD-EQ":    {},
	"CENTURYPLY-EQ": {},
	"CERA-EQ":       {},
	"CHAMBLFERT-EQ": {},
	"COROMANDEL-EQ": {},
	"CROMPTON-EQ":   {},
	"CUMMINSIND-EQ": {},
	"CYIENT-EQ":     {},
	"DABUR-EQ":      {},
	"DEEPAKFERT-EQ": {},
	"DEEPAKNTR-EQ":  {},
	"DELHIVERY-EQ":  {},
	"DELTACORP-EQ":  {},
	"DMART-EQ":      {},
	"EDELWEISS-EQ":  {},
	"EICHERMOT-EQ":  {},
	"EMAMILTD-EQ":   {},
	"EXIDEIND-EQ":   {},
	"GAIL-EQ":       {},
	"GLAXO-EQ":      {},
	"GLENMARK-EQ":   {},
	"GOKUL-EQ":      {},
	"HATHWAY-EQ":    {},
	"HAVELLS-EQ":    {},
	"HCLTECH-EQ":    {},
	"HDFCLIFE-EQ":   {},
	"HEG-EQ":        {},
	"HEROMOTOCO-EQ": {},
	"HINDUNILVR-EQ": {},
	"IDBI-EQ":       {},
	"IDFC-EQ":       {},
	"IDFCFIRSTB-EQ": {},
	"IMAGICAA-EQ":   {},
	"INDUSINDBK-EQ": {},
	"INFIBEAM-EQ":   {},
	"INOXINDIA-EQ":  {},
	"IRCTC-EQ":      {},
	"ITC-EQ":        {},
	"JAMNAAUTO-EQ":  {},
	"JINDALSTEL-EQ": {},
	"JKCEMENT-EQ":   {},
	"JKTYRE-EQ":     {},
	"JUBLFOOD-EQ":   {},
	"JUBLINDS-EQ":   {},
	"JUSTDIAL-EQ":   {},
	"KOLTEPATIL-EQ": {},
	"KPITTECH-EQ":   {},
	"LALPATHLAB-EQ": {},
	"LEMONTREE-EQ":  {},
	"LODHA-EQ":      {},
	"LT-EQ":         {},
	"LUPIN-EQ":      {},
	"MAPMYINDIA":    {},
	"MARICO-EQ":     {},
	"MARUTI-EQ":     {},
	"MOTHERSON-EQ":  {},
	"MOTILALOFS-EQ": {},
	"MPHASIS-EQ":    {},
	"MRF-EQ":        {},
	"MUTHOOTFIN-EQ": {},
	"NESTLEIND-EQ":  {},
	"NILKAMAL-EQ":   {},
	"NTPC-EQ":       {},
	"NYKAA-EQ":      {},
	"OBEROIRLTY-EQ": {},
	"OMAXE-EQ":      {},
	"ONGC-EQ":       {},
	"PAGEIND-EQ	":   {},
	"PAYTM-EQ":      {},
	"PCJEWELLER-EQ": {},
	"PERSISTENT-EQ": {},
	"PIDILITIND-EQ": {},
	"PNB-EQ":        {},
	"PVRINOX-EQ":    {},
	"RAIN-EQ":       {},
	"RECLTD-EQ":     {},
	"SAIL-EQ":       {},
	"SBIN-EQ":       {},
	"SIEMENS-EQ":    {},
	"SKFINDIA-EQ":   {},
	"SOBHA-EQ":      {},
	"SRF-EQ":        {},
	"STAR-EQ":       {},
	"SUNTV-EQ":      {},
	"TASTYBITE-EQ":  {},
	"TATACHEM-EQ":   {},
	"TATACOFFEE-EQ": {},
	"TATACOMM-EQ":   {},
	"TATACONSUM-EQ": {},
	"TATAELXSI-EQ":  {},
	"TATAINVEST-EQ": {},
	"TATAMETALI-EQ": {},
	"TATAMOTORS-EQ": {},
	"TATAMTRDVR-EQ": {},
	"TATAPOWER-EQ":  {},
	"TATASTEEL-EQ":  {},
	"TATATECH-EQ":   {},
	"TORNTPOWER-EQ": {},
	"TRENT-EQ":      {},
	"TV18BRDCST-EQ": {},
	"TVSMOTOR-EQ":   {},
	"UNIONBANK-EQ":  {},
	"UNITEDTEA-EQ":  {},
	"VENKEYS-EQ":    {},
	"VGUARD-EQ":     {},
	"VIPIND-EQ":     {},
	"VOLTAS-EQ":     {},
	"WHIRLPOOL-EQ":  {},
	"YATRA-EQ":      {},
	"YESBANK-EQ":    {},
	"ZEEL-EQ":       {},
	"ZEEMEDIA-EQ":   {},
	"ZOMATO-EQ":     {},
}

var indiceInstruments = map[string]struct{}{
	"Nifty 50":          {},
	"Nifty Consumption": {},
	"Nifty IT":          {},
	"Nifty Pharma":      {},
	"Nifty PSU Bank":    {},
	"Nifty Metal":       {},
	"Nifty Bank":        {},
	"Nifty 100":         {},
	"Nifty 500":         {},
	"Nifty Energy":      {},
	"Nifty Realty":      {},
	"Nifty Next 50":     {},
	"Nifty Auto":        {},
	"Nifty Infra":       {},
	"SENSEX":            {},
}

func isNifty50(instrument smartapi.InstrumentRecord) bool {
	_, present := nifty50Instruments[instrument.Symbol]
	return instrument.ExchSeg == "NSE" && present
}

func isSubscribed(instrument smartapi.InstrumentRecord) bool {
	_, present := subsciptionInstruments[instrument.Symbol]
	return instrument.ExchSeg == "NSE" && present
}

func StartAlphaHft(placeOrdersFlag bool) {
	logger.Log.Info().Msg("booting up alpha hft")

	testOHLCVRsiForToken()

	historicalOHLCVTimeSeriesData, err := LoadHistoricalOHLCVTimeSeriesData()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("error in loading historical ohlcv time series data")
		return
	}

	messaging.InitializeFIFOFile()

	if config.Config.RunTimeConfig.SimulationMode {
		logger.Log.Info().Str("simulationDate", config.Config.RunTimeConfig.SimulationDate).
			Msg("Simulation mode enabled")
		go handler.SimulateOldTickData(config.Config.RunTimeConfig.SimulationDate)
	} else {
		InitializeBrokersFromCurrentState(true)
		go handler.NewTickDataStorage().StoreTickData()
		go handler.NewLiveCandlesticksStorage().StoreCandlesticks()
	}

	alphaSignalGenerator := ta.NewAlphaSignalGenerator(historicalOHLCVTimeSeriesData)
	go alphaSignalGenerator.GenerateAlphaSignals()

	if !config.Config.RunTimeConfig.SimulationMode {
		go handler.NewAlphaSignalProcessor().ProcessAlphaSignals(false, placeOrdersFlag)
		go handler.NewAlphaSignalStorage().StoreAlphaSignals()
	}

	go messaging.ReadFromFIFO()
	go handler.NewCandlestickHandler().ProcessTicks()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan
}

func StartAlphaHftHistoricalBot() {
	logger.Log.Info().Msg("booting up alpha hft for historical data")

	historicalOHLCVTimeSeriesData, err := LoadHistoricalOHLCVTimeSeriesData()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("error in loading historical ohlcv time series data")
		return
	}

	InitializeBrokersFromCurrentState(true)

	messaging.InitializeFIFOFile()

	alphaSignalGenerator := ta.NewAlphaSignalGenerator(historicalOHLCVTimeSeriesData)
	go alphaSignalGenerator.GenerateHistoricalAlphaSignals()

	go handler.NewAlphaSignalStorage().StoreAlphaSignals()

	go messaging.ReadFromFIFO()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan
}

func testOHLCVRsiForToken() {
	olhcvtimeSeriesForToken, _ := data.GetOHLCVTimeSeriesIntervalToken("one_day", 25, "NSE")
	logger.Log.Debug().Msg("olhcvtimeSeriesForToken: " + fmt.Sprintf("%+v", olhcvtimeSeriesForToken))
	rsi14Token25 := talib.Rsi(olhcvtimeSeriesForToken.Close, 14)
	logger.Log.Debug().Msg("rsi14Token25: " + fmt.Sprintf("%+v", rsi14Token25))

	isHighRule := rules.NewIsHighestRule(rsi14Token25, 25)
	isHighRuleSatisfied := isHighRule.IsSatisfied(len(rsi14Token25) - 1)
	logger.Log.Debug().Msg("isHighRuleSatisfied: " + fmt.Sprintf("%+v", isHighRuleSatisfied))
}

func LoadHistoricalCandleData() (map[string]map[int]map[string][]smartapi.CandleData, error) {
	logger.Log.Info().Msg("Starting reading historical data for all instruments and intervals...")

	var historicalCandleDataMap = make(map[string]map[int]map[string][]smartapi.CandleData)
	var intervalLocks = make(map[string]*sync.Mutex)
	var historicalCandleDataMapMutex = &sync.Mutex{}
	var historicalCandleDataWaitGroup sync.WaitGroup
	var connLock sync.Mutex

	for _, interval := range data.HistoricTableIntervals {
		intervalLocks[strings.ToLower(interval.Interval)] = &sync.Mutex{}
	}

	for _, interval := range data.HistoricTableIntervals {
		intervalName := strings.ToLower(interval.Interval)

		if historicalCandleDataMap[intervalName] == nil {
			historicalCandleDataMapMutex.Lock()
			historicalCandleDataMap[intervalName] = make(map[int]map[string][]smartapi.CandleData)
			historicalCandleDataMapMutex.Unlock()
		}

		for _, exchgCode := range smartapi.CodeExchangeTypes {
			historicalCandleDataWaitGroup.Add(1)
			go func(interval data.Interval, exchangeCode string) {
				defer historicalCandleDataWaitGroup.Done()

				connLock.Lock()
				historicalCandleData, err := data.GetHistoricalDataForInterval(intervalName, exchangeCode)
				connLock.Unlock()

				if err != nil {
					logger.Log.Fatal().Str("interval", intervalName).Err(err).Msg("Error fetching historical data")
					return
				}

				intervalLock := intervalLocks[intervalName]
				intervalLock.Lock()
				defer intervalLock.Unlock()

				historicalCandleDataMapMutex.Lock()
				for _, candleData := range historicalCandleData {

					token := candleData.Token

					if historicalCandleDataMap[intervalName][token] == nil {
						historicalCandleDataMap[intervalName][token] = make(map[string][]smartapi.CandleData)
					}

					historicalCandleDataMap[intervalName][token][exchangeCode] =
						append(historicalCandleDataMap[intervalName][token][exchangeCode], candleData)
				}
				historicalCandleDataMapMutex.Unlock()
			}(interval, exchgCode)
		}
	}

	historicalCandleDataWaitGroup.Wait()

	logger.Log.Info().Msg("Finished reading historical data for all instruments and intervals!")

	return historicalCandleDataMap, nil
}

func LoadHistoricalOHLCVTimeSeriesData() (map[string]map[int]map[string]*common.OHLCVTimeSeries, error) {
	historicalCandleDataMap, err := LoadHistoricalCandleData()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error fetching historical data")
		return nil, err
	}

	logger.Log.Info().Msg("Converting the historical data to ohlcv times series data for all intervals and tokens...")

	ohlcvTimeSeriesForIntervalToken := make(map[string]map[int]map[string]*common.OHLCVTimeSeries)

	for interval, tokenData := range historicalCandleDataMap {
		intervalData := make(map[int]map[string]*common.OHLCVTimeSeries)

		for token, tokenExchangeData := range tokenData {
			ohlcvTimeSeries := common.NewOHLCVTimeSeries(token)

			intervalData[token] = make(map[string]*common.OHLCVTimeSeries)

			for exchangeCode, candleDataSlice := range tokenExchangeData {

				for _, candleData := range candleDataSlice {
					ohlcvTimeSeries.Timestamp = append(ohlcvTimeSeries.Timestamp, candleData.Timestamp)
					ohlcvTimeSeries.Open = append(ohlcvTimeSeries.Open, candleData.Open)
					ohlcvTimeSeries.High = append(ohlcvTimeSeries.High, candleData.High)
					ohlcvTimeSeries.Low = append(ohlcvTimeSeries.Low, candleData.Low)
					ohlcvTimeSeries.Close = append(ohlcvTimeSeries.Close, candleData.Close)
					ohlcvTimeSeries.Volume = append(ohlcvTimeSeries.Volume, candleData.Volume)
				}

				intervalData[token][exchangeCode] = &ohlcvTimeSeries
			}
		}

		ohlcvTimeSeriesForIntervalToken[interval] = intervalData
	}

	logger.Log.Info().
		Msg("Finished Converting the historical data to ohlcv times series data for all intervals and tokens!")

	return ohlcvTimeSeriesForIntervalToken, nil
}

func InitializeAlphaHftSOD(smartAPITotp string, createTables bool, populateOldHistoricData bool,
	useBrokersSavedState bool, useProxy bool, instrumentFilterFlag, intervalFilterFlag string) {
	// TODO this will take in the totp and rest,
	// To run before startalphahft
	// initalizeing table if it does not exists
	startTime := time.Now()
	logger.Log.Info().Msg("Initializing alpha hft started at: " + startTime.String())

	if useBrokersSavedState {
		InitializeBrokersFromCurrentState(false)
	} else {
		InitializeBrokers(smartAPITotp)
	}

	if createTables {
		data.InitializeTables()
	}

	client := common.GetClient(useProxy)

	instrumentsData, err := smartapi.GetInstrumentsScripMasterData(client)
	totalInstruments := len(instrumentsData)

	if err != nil {
		logger.Log.Error().Err(err).Msg("error in initalizng instruments data for sod")
		return
	}

	logger.Log.Info().Msg("Received the instruments data saving it in db")
	data.RefreshInstrumentsData(instrumentsData)

	logger.Log.Info().Msg("Fetching Nifty Indice Last Value")

	niftyIndiceLastTradePrice := data.GetNiftyIndiceLastValue()
	logger.Log.Info().
		Float64("niftyIndiceLastTradePrice", niftyIndiceLastTradePrice).
		Msg("Nifty Indice Last Value: ")

	logger.Log.Info().Msg("Fetching historic data for all instruments")

	shouldSkip := func(instrument smartapi.InstrumentRecord) bool {
		// Process all the instruments
		if instrumentFilterFlag == "ALL" {
			return false
		}

		// Process all the relavant instruments.
		if instrumentFilterFlag == "ALL-EQ" {
			if instrument.ExchSeg == "NSE" &&
				(instrument.InstrumentTypeCode == "EQ" || instrument.InstrumentTypeCode == "") {
				return false
			}
		}

		if instrumentFilterFlag == "NIFTY50" {
			return !isNifty50(instrument)
		}

		if instrumentFilterFlag == "SUBSCRIPTION" {
			return !isSubscribed(instrument)
		}

		if instrumentFilterFlag == "SUBSCRIPTION_NIFTY50" {
			return !isNifty50(instrument) && !isSubscribed(instrument)
		}

		if instrumentFilterFlag == "NIFTYOPTNEAR" &&
			instrument.ExchSeg == "NFO" &&
			instrument.InstrumentType == "OPTIDX" &&
			instrument.Name == "NIFTY" {
			maxStrike := niftyIndiceLastTradePrice + 2000
			minStrike := niftyIndiceLastTradePrice - 2000

			instrumentStrike, err := strconv.ParseFloat(instrument.Strike, 64)
			if err != nil {
				logger.Log.
					Err(err).
					Msg("Error parsing instrument strike price")

				return true
			}

			if instrumentStrike > maxStrike || instrumentStrike < minStrike {
				return false
			}

			return true
		}

		// For NSE Exchange process EQ or indices for other exchanges process all instruments.
		if instrument.ExchSeg != instrumentFilterFlag {
			return true
		} else if instrument.ExchSeg == "NSE" &&
			instrument.InstrumentTypeCode != "EQ" &&
			instrument.InstrumentTypeCode != "" {
			return true
		}

		return false
	}

	for i, instrument := range instrumentsData {
		if shouldSkip(instrument) {
			logger.Log.Debug().
				Str("exchange", instrument.ExchSeg).
				Str("symbol", instrument.Symbol).
				Msg("Skipping")

			continue
		}

		logger.Log.Info().
			Str("Percent Complete", fmt.Sprintf("%.3f%%", float64((i+1)*100)/float64(totalInstruments))).
			Str("Total Instruments", fmt.Sprint(totalInstruments)).
			Str("instrument number", fmt.Sprint(i+1)).
			Str("symbol", instrument.Symbol).
			Str("exchseg", instrument.ExchSeg).
			Msg("Processing instrument: ")
		logger.Log.Debug().Msg("*****************************************")

		interval := time.Second / time.Duration(HistoricalPricesLimitRatePerSecond)
		ticker := time.NewTicker(interval)

		for _, interval := range data.HistoricTableIntervals {
			if intervalFilterFlag == "ALL" || intervalFilterFlag == interval.Interval {
				logger.Log.Debug().
					Str("symbol", instrument.Symbol).
					Str("interval", interval.Interval).
					Msg("Fetching historical data for")

				fromDate, toDate := smartapi.GetFromDateToDateForHistoricalData(populateOldHistoricData, interval.MaxDays)

				<-ticker.C

				intervalRecords, err := smartapi.SmartApiBrokers["historical"].GetCandleDataForIntervalRange(client,
					instrument.Token, instrument.ExchSeg, interval.Interval, fromDate, toDate)
				if err != nil {
					logger.Log.Error().Msg("error in getting record skipping")
					continue
				}

				if len(intervalRecords) == 0 {
					logger.Log.Error().Msg("No record for interval. Not processing any further intervals for the instrument")
					break
				}

				logger.Log.Info().Str("symbolToken", instrument.Token).Str("interval", interval.Interval).
					Int("size", len(intervalRecords)).Msg("received interval records saving")

				parsedToken, _ := strconv.Atoi(instrument.Token)
				data.SaveHistoricalCandleData(interval.Interval, parsedToken, instrument.ExchSeg, fromDate, toDate, intervalRecords)
			} else {
				logger.Log.Debug().
					Str("interval", interval.Interval).
					Msg("Skipping the interval based on the filter flag")
			}
		}

		logger.Log.Debug().Msg("*****************************************")
	}

	endTime := time.Now()
	logger.Log.Info().Msg("Initializing alpha hft finished at: " + endTime.String())

	elapsedTime := endTime.Sub(startTime)
	logger.Log.Info().Msg("Execution time taken: " + elapsedTime.String())
}
