package data

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sudeepbatra/alpha-hft/broker/smartapi"
	"github.com/sudeepbatra/alpha-hft/common"
	"github.com/sudeepbatra/alpha-hft/logger"
)

const (
	InstrumentTable         = "smartapi_instruments"
	last_nifty_indice_value = 21349.40
)

func RefreshInstrumentsData(instrumentsData []smartapi.InstrumentRecord) {
	_, err := AlphaHftDbConn.Exec(context.Background(), "TRUNCATE TABLE "+InstrumentTable)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Unable to truncate table failed to refresh instruments data")
		return
	}

	err = BulkInsert(InstrumentTable, instrumentsData)

	if err != nil {
		logger.Log.Error().Err(err).Msg("error in refreshing instruments data")
		return
	}

	logger.Log.Info().Msg("successfully refreshed instruments data")
}

func SaveHistoricalCandleData(interval string, token int, exchange, fromDate, toDate string, candlesData []smartapi.CandleData) {
	interval = strings.ToLower(interval)

	_, err := AlphaHftDbConn.Exec(context.Background(), fmt.Sprintf(HistoricalDeleteQueryRange, interval,
		token, exchange, fromDate, toDate))
	if err != nil {
		logger.Log.Error().Err(err).Msg("Unable to delete previous record for the range")
		return
	}

	err = BulkInsert(interval, candlesData)

	if err != nil {
		logger.Log.Error().Err(err).Msg("error in saving candle data")
		return
	}

	logger.Log.Info().Msg("successfully saved candle data")
}

func GetInstrumentsData() ([]smartapi.InstrumentRecord, error) {
	rows, err := AlphaHftDbConn.Query(context.Background(), SmartAPIInstrumentGetQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Create a slice to hold the result.
	var result []smartapi.InstrumentRecord

	// Iterate through the rows and scan data into the custom structure.
	for rows.Next() {
		var instrument smartapi.InstrumentRecord

		err := rows.Scan(&instrument.Token, &instrument.Symbol, &instrument.Name, &instrument.Expiry, &instrument.Strike, &instrument.Lotsize, &instrument.InstrumentType, &instrument.ExchSeg, &instrument.TickSize, &instrument.InstrumentTypeCode)
		if err != nil {
			return nil, err
		}

		result = append(result, instrument)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func GetInstrumentsSymbolAndTokenMapping() (
	map[smartapi.InstrumentSymbolKey]smartapi.InstrumentTokenKey,
	map[smartapi.InstrumentTokenKey]smartapi.InstrumentSymbolKey,
	error,
) {
	instruments, err := GetInstrumentsData()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error fetching instruments data")
		return nil, nil, err
	}

	symbolToTokenMap := make(map[smartapi.InstrumentSymbolKey]smartapi.InstrumentTokenKey)
	tokenToSymbolMap := make(map[smartapi.InstrumentTokenKey]smartapi.InstrumentSymbolKey)

	for _, instrument := range instruments {
		instrumentSymbolKey := smartapi.InstrumentSymbolKey{Symbol: instrument.Symbol, Exchange: instrument.ExchSeg}
		instrumentTokenKey := smartapi.InstrumentTokenKey{Token: instrument.Token, Exchange: instrument.ExchSeg}

		symbolToTokenMap[instrumentSymbolKey] = instrumentTokenKey
		tokenToSymbolMap[instrumentTokenKey] = instrumentSymbolKey
	}

	return symbolToTokenMap, tokenToSymbolMap, nil
}

func GetHistoricalDataForIntervalForToken(interval string, token int) ([]smartapi.CandleData, error) {
	query := fmt.Sprintf(HistoricalCandleDataGetQuery, strings.ToLower(interval), token)

	conn, err := AlphaHftDbConnPool.Acquire(context.Background())
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in saving candle data")
		return nil, err
	}

	defer conn.Release()

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		logger.Log.Error().
			Str("interval", interval).
			Int("token", token).
			Err(err).
			Msg("Unable to get historical data for interval and token")

		return nil, err
	}

	defer rows.Close()

	var candleDataList []smartapi.CandleData

	for rows.Next() {
		var candleData smartapi.CandleData

		// Scan the values from the current row into the variables
		err := rows.Scan(&candleData.Token,
			&candleData.Timestamp,
			&candleData.Open,
			&candleData.High,
			&candleData.Low,
			&candleData.Close,
			&candleData.Volume)
		if err != nil {
			logger.Log.Error().
				Err(err).
				Msg("Error scanning row")

			return nil, err
		}

		candleDataList = append(candleDataList, candleData)
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error().
			Err(err).
			Msg("Error iterating rows")

		return nil, err
	}

	return candleDataList, nil
}

func GetHistoricalDataForInterval(interval string, exchangeCode string) ([]smartapi.CandleData, error) {
	query := fmt.Sprintf(HistoricalCandleDataIntervalExchangeCodeGetQuery, strings.ToLower(interval), exchangeCode)

	conn, err := AlphaHftDbConnPool.Acquire(context.Background())
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in saving candle data")
		return nil, err
	}

	defer conn.Release()

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		logger.Log.Error().
			Str("interval", interval).
			Err(err).
			Msg("Unable to get historical data for interval")

		return nil, err
	}

	defer rows.Close()

	var candleDataList []smartapi.CandleData

	for rows.Next() {
		var candleData smartapi.CandleData

		// Scan the values from the current row into the variables
		err := rows.Scan(&candleData.Token,
			&candleData.Exchange,
			&candleData.Timestamp,
			&candleData.Open,
			&candleData.High,
			&candleData.Low,
			&candleData.Close,
			&candleData.Volume)
		if err != nil {
			logger.Log.Error().
				Err(err).
				Msg("Error scanning row")

			return nil, err
		}

		candleDataList = append(candleDataList, candleData)
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error().
			Err(err).
			Msg("Error iterating rows")

		return nil, err
	}

	return candleDataList, nil
}

func GetOHLCVTimeSeriesIntervalToken(interval string, token int, exchange string) (common.OHLCVTimeSeries, error) {
	var rowCount int

	queryCount := fmt.Sprintf(HistoricalCandleDataCountQuery, strings.ToLower(interval), token, exchange)

	err := AlphaHftDbConnPool.QueryRow(context.Background(), queryCount).Scan(&rowCount)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in getting ohlcv time series interval token count")
		return common.OHLCVTimeSeries{}, err
	}

	query := fmt.Sprintf(OHLCVTimeSeriesQuery, strings.ToLower(interval), token, exchange)

	conn, err := AlphaHftDbConnPool.Acquire(context.Background())
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in saving candle data")
		return common.OHLCVTimeSeries{}, err
	}

	defer conn.Release()

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		logger.Log.Error().
			Str("interval", interval).
			Int("token", token).
			Err(err).
			Msg("Unable to get historical data for interval and token")

		return common.OHLCVTimeSeries{}, err
	}

	defer rows.Close()

	ohlcvTimeSeriesForToken := common.NewOHLCVTimeSeriesWithBarCount(token, rowCount)

	rowNumber := 0

	for rows.Next() {
		var Token string
		var Timestamp time.Time
		var Open, High, Low, Close float64
		var Volume int

		err := rows.Scan(&Token, &Timestamp, &Open, &High, &Low, &Close, &Volume)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error scanning row")
			return common.OHLCVTimeSeries{}, err
		}

		ohlcvTimeSeriesForToken.Timestamp[rowNumber] = Timestamp
		ohlcvTimeSeriesForToken.Open[rowNumber] = Open
		ohlcvTimeSeriesForToken.High[rowNumber] = High
		ohlcvTimeSeriesForToken.Low[rowNumber] = Low
		ohlcvTimeSeriesForToken.Close[rowNumber] = Close
		ohlcvTimeSeriesForToken.Volume[rowNumber] = Volume
		rowNumber++
	}

	return ohlcvTimeSeriesForToken, nil
}

func GetNiftyIndexLastValue() float64 {
	var lastTradedPrice float64

	err := AlphaHftDbConn.QueryRow(context.Background(), NiftyIndiceLastTradePriceQuery).Scan(&lastTradedPrice)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Msg("error in getting last traded price")

		return last_nifty_indice_value
	}

	return lastTradedPrice
}
