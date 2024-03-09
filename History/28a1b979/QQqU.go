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

	CreateAlphaSignalsTableQuery = `
		CREATE TABLE IF NOT EXISTS alpha_signals (
			token VARCHAR(255),
			symbol VARCHAR(255),
			exchangecode INT,
			exchange VARCHAR(5),
			interval VARCHAR(255),
			lastbarstartduration  TIMESTAMPTZ,
			strategy VARCHAR(255),
			signal VARCHAR(255),
			signalgenerationtime TIMESTAMPTZ,
			price NUMERIC,
			o NUMERIC,
			h NUMERIC,
			l NUMERIC,
			c NUMERIC,
			v INT,
			alphasignalreason VARCHAR(255),
			ishistorical BOOLEAN
		);
	`

	CreateTickDataTableQuery = `
		CREATE TABLE tick_data (
			id SERIAL PRIMARY KEY,
			subscriptionmode             SMALLINT,
			exchangetype                 SMALLINT,
			token                        VARCHAR(255),
			sequencenumber               BIGINT,
			exchangetimestamp            TIMESTAMP,
			lasttradedprice              BIGINT,
			subscriptionmodeval          VARCHAR(255),
			lasttradedquantity           BIGINT,
			averagetradedprice           BIGINT,
			volumetradefortheday         BIGINT,
			totalbuyquantity             DOUBLE PRECISION,
			totalsellquantity            DOUBLE PRECISION,
			openpriceoftheday            BIGINT,
			highpriceoftheday            BIGINT,
			lowpriceoftheday             BIGINT,
			closedprice                  BIGINT,
			lasttradedtimestamp          BIGINT,
			openinterest                 BIGINT,
			openinterestchangepercentage BIGINT,
			uppercircuitlimit            BIGINT,
			lowercircuitlimit            BIGINT,
			week52highprice              BIGINT,
			week52lowprice               BIGINT,
			lasttradedpricefloat         DOUBLE PRECISION,
			best5buydata                 BIGINT[],
			best5selldata                BIGINT[]
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
