package data

const (
	OneMinute     = "ONE_MINUTE"
	ThreeMinute   = "THREE_MINUTE"
	FiveMinute    = "FIVE_MINUTE"
	TenMinutes    = "TEN_MINUTE"
	FifteenMinute = "FIFTEEN_MINUTE"
	ThirtyMinute  = "THIRTY_MINUTE"
	OneHour       = "ONE_HOUR"
	OneDay        = "ONE_DAY"
)

type Interval struct {
	Interval string
	MaxDays  int
}

var OneMinuteInterval = Interval{OneMinute, 30}
var ThreeMinuteInterval = Interval{ThreeMinute, 90}
var FiveMinuteInterval = Interval{FiveMinute, 90}
var TenMinutesInterval = Interval{TenMinutes, 90}
var FifteenMinuteInterval = Interval{FifteenMinute, 180}
var ThirtyMinuteInterval = Interval{ThirtyMinute, 180}
var OneHourInterval = Interval{OneHour, 365}
var OneDayInterval = Interval{OneDay, 2000}

var HistoricTableIntervals = [...]Interval{OneMinuteInterval, ThreeMinuteInterval, FiveMinuteInterval, TenMinutesInterval, FifteenMinuteInterval, ThirtyMinuteInterval, OneHourInterval, OneDayInterval}

const (
	InstrumentCreateQuery = `
		CREATE TABLE IF NOT EXISTS smartapi_instruments (
			token VARCHAR(255),
			symbol VARCHAR(255),
			name VARCHAR(255),
			expiry VARCHAR(255),
			strike NUMERIC,
			lotsize INT,
			instrumenttype VARCHAR(20),
			exchseg VARCHAR(20),
			ticksize NUMERIC,
			instrumenttypecode VARCHAR(5)
		);
	`
	HistoricalTableQuery = `
		CREATE TABLE IF NOT EXISTS %s (
			token INT,
			exchange VARCHAR(5),
			timestamp TIMESTAMPTZ,
			open DECIMAL(10, 2),
			high DECIMAL(10, 2),
			low DECIMAL(10, 2),
			close DECIMAL(10, 2),
			volume INT
		);
	`

	type AlphaSignal struct {
		Token                string
		Symbol               string
		ExchangeCode         int
		Exchange             string
		Interval             string
		LastBarStartDuration time.Time
		Strategy             string
		Signal               string
		SignalGenerationTime time.Time
		Price                float64
		O                    float64
		H                    float64
		L                    float64
		C                    float64
		V                    int
		AlphaSignalReason    string
		IsHistorical         bool
	}
	CreateAlphaSignalsTableQuery = `
		CREATE TABLE IF NOT EXISTS %s (
			token VARCHAR(255),
			symbol VARCHAR(255),
			exchange_code INT,
			exchange VARCHAR(5),
			interval VARCHAR(255),
			last_bar_start_duration  TIMESTAMPTZ,

			timestamp TIMESTAMPTZ,
			open DECIMAL(10, 2),
			high DECIMAL(10, 2),
			low DECIMAL(10, 2),
			close DECIMAL(10, 2),
			volume INT
		);
	`

	SmartAPIInstrumentGetQuery = `
		SELECT 
		 token,
		 symbol,
		 name,
		 expiry,
		 strike,
		 lotsize,
		 instrumenttype,
		 exchseg,
		 ticksize,
		 instrumenttypecode
		from smartapi_instruments;
	`
	HistoricalCandleDataGetQuery = `
		SELECT * from %s WHERE token = %d ORDER BY timestamp ASC;
	`

	OHLCVTimeSeriesQuery = `
		SELECT token, timestamp, open, high, low, close, volume from %s WHERE token = %d AND exchange = '%s' ORDER BY timestamp ASC;
	`

	HistoricalCandleDataCountQuery = `
		SELECT COUNT(*) from %s WHERE token = %d AND exchange = '%s';
	`

	HistoricalCandleDataIntervalGetQuery = `
		SELECT token, exchange, timestamp, open, high, low, close, volume from %s ORDER BY token, timestamp ASC;
	`

	HistoricalCandleDataIntervalExchangeCodeGetQuery = `
	SELECT token, exchange, timestamp, open, high, low, close, volume from %s where exchange= '%s' ORDER BY token, timestamp ASC;
	`

	HistoricalDeleteQueryRange = `
		DELETE FROM %s WHERE token = %d AND exchange = '%s' AND timestamp >= '%s' AND timestamp <= '%s';
	`
)