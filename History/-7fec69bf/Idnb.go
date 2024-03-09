package data

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sudeepbatra/alpha-hft/config"
	"github.com/sudeepbatra/alpha-hft/logger"
)

var (
	AlphaHftDbConn     *pgx.Conn
	AlphaHftDbConnPool *pgxpool.Pool
)

func InitializeTables() {
	_, err := AlphaHftDbConn.Exec(context.Background(), InstrumentCreateQuery)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error in creating instrument table")
	}

	logger.Log.Trace().Str("table", InstrumentTable).Msg("Table created or it already exists.")

	_, err = AlphaHftDbConn.Exec(context.Background(), CreateAlphaSignalsTableQuery)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error in creating alpha signals table")
	}

	err = createTickDataTable()
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error in creating tick data table")
		return
	}

	for _, interval := range HistoricTableIntervals {
		tableName := interval.Interval
		hisotricalTableCreateQuery := fmt.Sprintf(HistoricalTableQuery, tableName)
		hypertableQuery := fmt.Sprintf("SELECT create_hypertable('%s','%s');", tableName, "timestamp")
		indexQuery := fmt.Sprintf("CREATE INDEX ix_symbol_time_%s ON %s (%s, %s, %s DESC);",
			interval.Interval, tableName, "token", "exchange", "timestamp")

		_, err := AlphaHftDbConn.Exec(context.Background(), hisotricalTableCreateQuery)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error in creating historical table")
			continue
		}

		logger.Log.Info().Str("table", tableName).Msg("Table created or it already exists.")

		_, err = AlphaHftDbConn.Exec(context.Background(), hypertableQuery)

		if err != nil {
			logger.Log.Error().Err(err).Msg("Error in creating hypertable on historical tables")
			continue
		}

		_, err = AlphaHftDbConn.Exec(context.Background(), indexQuery)

		if err != nil {
			logger.Log.Error().Err(err).Msg("Error in creating index on historical tables")
			continue
		}

	}
}

func createTickDataTable() error {
	_, err := AlphaHftDbConn.Exec(context.Background(), CreateLiveCandlestickTableQuery)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error in creating live candlestick table")
		return err
	}

	liveCandlesticksHyperTableQuery := "SELECT create_hypertable('live_candlesticks', 'timestamp');"

	_, err = AlphaHftDbConn.Exec(context.Background(), liveCandlesticksHyperTableQuery)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error in creating hypertable on live candlesticks table")
		return err
	}

	tickDataIndexQuery := fmt.Sprintf("CREATE INDEX idx_token_exchange_timestamp ON live_candlesticks (token, exchangetype, exchangetimestamp DESC);")

	_, err = AlphaHftDbConn.Exec(context.Background(), tickDataIndexQuery)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error in creating index on tick data tables")
		return err
	}

	return nil
}

func createLiveCandlestickTable() error {
	_, err := AlphaHftDbConn.Exec(context.Background(), CreateTickDataTableQuery)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error in creating tick data table")
		return err
	}

	tickDataHyperTableQuery := "SELECT create_hypertable('tick_data', 'exchangetimestamp');"

	_, err = AlphaHftDbConn.Exec(context.Background(), tickDataHyperTableQuery)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error in creating hypertable on tick data table")
		return err
	}

	tickDataIndexQuery := fmt.Sprintf("CREATE INDEX idx_token_exchange_timestamp ON tick_data (token, exchangetype, exchangetimestamp DESC);")

	_, err = AlphaHftDbConn.Exec(context.Background(), tickDataIndexQuery)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error in creating index on tick data tables")
		return err
	}

	return nil
}

func init() {
	var err error

	AlphaHftDbConn, err = pgx.Connect(context.Background(), config.Config.DatabaseConfig.MainTsdbUri)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Unable to initialize timescale db for our bot")
		os.Exit(1)
	}

	AlphaHftDbConnPool, err = pgxpool.New(context.Background(), config.Config.DatabaseConfig.MainTsdbUri)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Unable to initialize AlphaHftDbConnPool")
		os.Exit(1)
	}
}
