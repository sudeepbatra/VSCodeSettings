package data

const (
	OneMinute        = "ONE_MINUTE"
	TwoMinute        = "TWO_MINUTE"
	ThreeMinute      = "THREE_MINUTE"
	FiveMinute       = "FIVE_MINUTE"
	TenMinutes       = "TEN_MINUTE"
	FifteenMinute    = "FIFTEEN_MINUTE"
	ThirtyMinute     = "THIRTY_MINUTE"
	FourtyFiveMinute = "FOURTYFIVE_MINUTE"
	OneHour          = "ONE_HOUR"
	TwoHour          = "TWO_HOUR"
	OneDay           = "ONE_DAY"
)

const (
	TwoMinuteTimeBucket        = "2 minute"
	ThreeMinuteTimeBucket      = "3 minute"
	FiveMinuteTimeBucket       = "5 minute"
	TenMinuteTimeBucket        = "10 minute"
	FifteenMinuteTimeBucket    = "15 minute"
	ThirtyMinuteTimeBucket     = "30 minute"
	FourtyFiveMinuteTimeBucket = "45 minute"
	OneHourTimeBucket          = "1 hour"
	TwoHourTimeBucket          = "2 hours"
)

const (
	TwoMinuteOffsetInterval        = "1 minute"
	ThreeMinuteOffsetInterval      = "0 minute"
	FiveMinuteOffsetInterval       = "0 minute"
	TenMinuteOffsetInterval        = "5 minute"
	FifteenMinuteOffsetInterval    = "15 minute"
	ThirtyMinuteOffsetInterval     = "45 minute"
	FourtyFiveMinuteOffsetInterval = "45 minute"
	OneHourOffsetInterval          = "45 minute"
	TwoHourOffsetInterval          = "1 hour 45 minute"
)

type IntervalAggregator struct {
	TimeBucketInterval        string
	TimeBucketOffsetIinterval string
}

type Interval struct {
	Interval   string
	MaxDays    int
	Aggregator IntervalAggregator
}

var OneMinuteInterval = Interval{OneMinute, 30, IntervalAggregator{}}
var TwoMinuteInterval = Interval{TwoMinute, 90, IntervalAggregator{
	TwoHourTimeBucket,
	TwoHourOffsetInterval,
}}
var ThreeMinuteInterval = Interval{ThreeMinute, 90, IntervalAggregator{
	ThreeMinuteTimeBucket,
	ThreeMinuteOffsetInterval,
}}
var FiveMinuteInterval = Interval{FiveMinute, 90, IntervalAggregator{
	FiveMinuteTimeBucket,
	FiveMinuteOffsetInterval,
}}
var TenMinutesInterval = Interval{TenMinutes, 90, IntervalAggregator{
	TenMinuteTimeBucket,
	TenMinuteOffsetInterval,
}}
var FifteenMinuteInterval = Interval{FifteenMinute, 180, IntervalAggregator{
	FifteenMinuteTimeBucket,
	FifteenMinuteOffsetInterval,
}}
var ThirtyMinuteInterval = Interval{ThirtyMinute, 180, IntervalAggregator{
	ThirtyMinuteTimeBucket,
	ThirtyMinuteOffsetInterval,
}}

var FourtyFiveMinuteInterval = Interval{FourtyFiveMinute, 180, IntervalAggregator{
	FourtyFiveMinuteTimeBucket,
	FourtyFiveMinuteOffsetInterval,
}}

var OneHourInterval = Interval{OneHour, 365, IntervalAggregator{
	OneHourTimeBucket,
	OneHourOffsetInterval,
}}

var TwoHourInterval = Interval{TwoHour, 180, IntervalAggregator{
	TwoHourTimeBucket,
	TwoHourOffsetInterval,
}}

var OneDayInterval = Interval{OneDay, 2000, IntervalAggregator{}}

var HistoricTableIntervals = [...]Interval{OneMinuteInterval, TwoMinuteInterval, ThreeMinuteInterval, FiveMinuteInterval, TenMinutesInterval, FifteenMinuteInterval, ThirtyMinuteInterval,
	FourtyFiveMinuteInterval, OneHourInterval, TwoHourInterval, OneDayInterval}

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
			instrumenttypecode VARCHAR(5),
			optiontype VARCHAR(20)
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
			subscriptionmode             SMALLINT,
			exchangetype                 SMALLINT,
			token                        VARCHAR(255),
			sequencenumber               BIGINT,
			exchangetimestamp            TIMESTAMPTZ,
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
			lasttradedtimestamp          TIMESTAMPTZ,
			openinterest                 BIGINT,
			openinterestchangepercentage BIGINT,
			uppercircuitlimit            BIGINT,
			lowercircuitlimit            BIGINT,
			week52highprice              BIGINT,
			week52lowprice               BIGINT,
			lasttradedpricefloat         DOUBLE PRECISION
		);
`

	CreateLiveCandlesticksTableQuery = `
		CREATE TABLE live_candlesticks (
			timestamp TIMESTAMPTZ NOT NULL,
			unixtimestamp bigint,
			open DOUBLE PRECISION NOT NULL,
			high DOUBLE PRECISION NOT NULL,
			low DOUBLE PRECISION NOT NULL,
			close DOUBLE PRECISION NOT NULL,
			volume BIGINT NOT NULL,
			token VARCHAR(255) NOT NULL,
			exchange INTEGER NOT NULL,
			duration VARCHAR(50) NOT NULL
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

	NiftyIndiceLastTradePriceQuery = `
	SELECT lasttradedpricefloat FROM tick_data WHERE token = '99926000' ORDER BY exchangetimestamp DESC LIMIT 1
	`

	HistoricTickDataQuery = `
	WITH RankedData AS (
		SELECT
		  exchangetimestamp,
		  lasttradedpricefloat,
		  lasttradedquantity,
		  exchangetype,
		  token,
		  ROW_NUMBER() OVER (PARTITION BY token ORDER BY exchangetimestamp) AS row_num
		FROM
		  tick_data
		WHERE
		 exchangetimestamp::date = '%s' and
		 token in (%s)
	  )
	  SELECT
		  exchangetimestamp,
		  lasttradedpricefloat,
		  lasttradedquantity,
		  exchangetype,
		  token
	  FROM
		RankedData
	  ORDER BY
		row_num,
		exchangetimestamp;
		`

	// candle time frame aggregator queries
	OneMinCandleAggregatorQuery = `
	SELECT 
		time_bucket('%s', timestamp, '%s'::INTERVAL) AS bucket, 
		(
			array_agg(
			OPEN 
			order by 
				timestamp
			)
		) [1] AS open, 
		max(high) AS high, 
		min(low) AS low, 
		(
			array_agg(
			CLOSE 
			order by 
				timestamp DESC
			)
		) [1] AS close, 
		sum(volume) AS volume 
		FROM 
		one_minute 
		WHERE 
		token = %d 
		AND exchange = '%s'
		AND timestamp::date >= '%s'
		AND timestamp::date <= '%s' 
		GROUP BY 
		bucket 
		ORDER BY 
		bucket 
		`
)

CreateCnbcFairValueQuoteTableQuery = `
CREATE TABLE IF NOT EXISTS fair_value_quotes (
	id SERIAL PRIMARY KEY,
	symbol VARCHAR(255),
	code VARCHAR(255),
	name VARCHAR(255),
	last VARCHAR(255),
	last_time TIMESTAMP,
	last_time_string VARCHAR(255),
	last_time_msec VARCHAR(255),
	exchange VARCHAR(255),
	provider VARCHAR(255),
	todays_closing VARCHAR(255),
	provider_symbol VARCHAR(255),
	index_close VARCHAR(255),
	fv_close VARCHAR(255),
	fv_change VARCHAR(255),
	fv_spread VARCHAR(255),
	fv_raw VARCHAR(255),
	last_time_date VARCHAR(255),
	realtime VARCHAR(255),
	shortname VARCHAR(255),
	alt_symbol VARCHAR(255),
	issue_id VARCHAR(255),
	fmt_last VARCHAR(255),
	fmt_change VARCHAR(255),
	fv_change_pct VARCHAR(255),
	change_pct VARCHAR(255),
	change VARCHAR(255)
);
`

var SupportedCandlePeriodsAggregator = map[Interval]string{
	OneMinuteInterval: OneMinCandleAggregatorQuery,
}
