package smartapi

import (
	"time"
)

type InstrumentRecord struct {
	Token              string `json:"token"`
	Symbol             string `json:"symbol"`
	Name               string `json:"name"`
	Expiry             string `json:"expiry"`
	Strike             string `json:"strike"`
	Lotsize            string `json:"lotsize"`
	InstrumentType     string `json:"instrumenttype"`
	ExchSeg            string `json:"exch_seg"`
	TickSize           string `json:"tick_size"`
	InstrumentTypeCode string `json:"instrument_type_code"`
}

type Subscription struct {
	CorrelationID string        `json:"correlationID"`
	Action        int           `json:"action"`
	Params        RequestParams `json:"params"`
}

type RequestParams struct {
	Mode       int                `json:"mode"`
	TokenLists []RequestTokenList `json:"tokenList"`
}

type RequestTokenList struct {
	ExchangeType int      `json:"exchangeType"`
	Tokens       []string `json:"tokens"`
}

type ErrorResponse struct {
	CorrelationID string `json:"correlationID"`
	ErrorCode     string `json:"errorCode"`
	ErrorMessage  string `json:"errorMessage"`
}

type TickParsedData struct {
	SubscriptionMode             uint8
	ExchangeType                 uint8
	Token                        string
	SequenceNumber               int64
	ExchangeTimestamp            int64
	LastTradedPrice              int64
	SubscriptionModeVal          string
	LastTradedQuantity           int64
	AverageTradedPrice           int64
	VolumeTradeForTheDay         int64
	TotalBuyQuantity             float64
	TotalSellQuantity            float64
	OpenPriceOfTheDay            int64
	HighPriceOfTheDay            int64
	LowPriceOfTheDay             int64
	ClosedPrice                  int64
	LastTradedTimestamp          int64
	OpenInterest                 int64
	OpenInterestChangePercentage int64
	UpperCircuitLimit            int64
	LowerCircuitLimit            int64
	Week52HighPrice              int64
	Week52LowPrice               int64
	LastTradedPriceFloat         float64
	Best5BuyData, Best5SellData  []int64 // Modify the data types as required
}

type CandleData struct {
	Token     int
	Exchange  string
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    int
}

type HistoricalApiResponse struct {
	Status    bool            `json:"status"`
	Message   string          `json:"message"`
	ErrorCode string          `json:"errorcode"`
	Data      [][]interface{} `json:"data"`
}

type PlaceOrderResponse struct {
	Status    bool   `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"errorcode"`
	Data      struct {
		Script        string `json:"script"`
		OrderID       string `json:"orderid"`
		UniqueOrderID string `json:"uniqueorderid"`
	} `json:"data"`
}

type ModifyOrderResponse struct {
	Status    bool   `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"errorcode"`
	Data      struct {
		OrderID       string `json:"orderid"`
		UniqueOrderID string `json:"uniqueorderid"`
	} `json:"data"`
}

type CancelOrderResponse struct {
	Status    bool   `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"errorcode"`
	Data      struct {
		OrderID       string `json:"orderid"`
		UniqueOrderID string `json:"uniqueorderid"`
	} `json:"data"`
}

type OrderBookResponse struct {
	Status    bool            `json:"status"`
	Message   string          `json:"message"`
	ErrorCode string          `json:"errorcode"`
	Data      []OrderBookItem `json:"data"`
}

type OrderBookItem struct {
	Variety             string  `json:"variety"`
	OrderType           string  `json:"ordertype"`
	ProductType         string  `json:"producttype"`
	Duration            string  `json:"duration"`
	Price               float64 `json:"price"`
	TriggerPrice        float64 `json:"triggerprice"`
	Quantity            string  `json:"quantity"`
	DisclosedQuantity   string  `json:"disclosedquantity"`
	SquareOff           float64 `json:"squareoff"`
	StopLoss            float64 `json:"stoploss"`
	TrailingStopLoss    float64 `json:"trailingstoploss"`
	TradingSymbol       string  `json:"tradingsymbol"`
	TransactionType     string  `json:"transactiontype"`
	Exchange            string  `json:"exchange"`
	SymbolToken         string  `json:"symboltoken"`
	InstrumentType      string  `json:"instrumenttype"`
	StrikePrice         float64 `json:"strikeprice"`
	OptionType          string  `json:"optiontype"`
	ExpiryDate          string  `json:"expirydate"`
	LotSize             string  `json:"lotsize"`
	CancelSize          string  `json:"cancelsize"`
	AveragePrice        float64 `json:"averageprice"`
	FilledShares        string  `json:"filledshares"`
	UnfilledShares      string  `json:"unfilledshares"`
	OrderID             string  `json:"orderid"`
	Text                string  `json:"text"`
	Status              string  `json:"status"`
	OrderStatus         string  `json:"orderstatus"`
	UpdateTime          string  `json:"updatetime"`
	ExchTime            string  `json:"exchtime"`
	ExchOrderUpdateTime string  `json:"exchorderupdatetime"`
	FillID              string  `json:"fillid"`
	FillTime            string  `json:"filltime"`
	ParentOrderID       string  `json:"parentorderid"`
}

type TradeBookResponse struct {
	Status    bool            `json:"status"`
	Message   string          `json:"message"`
	ErrorCode string          `json:"errorcode"`
	Data      []TradeBookItem `json:"data"`
}

type TradeBookItem struct {
	Exchange        string `json:"exchange"`
	ProductType     string `json:"producttype"`
	TradingSymbol   string `json:"tradingsymbol"`
	InstrumentType  string `json:"instrumenttype"`
	SymbolGroup     string `json:"symbolgroup"`
	StrikePrice     string `json:"strikeprice"`
	OptionType      string `json:"optiontype"`
	ExpiryDate      string `json:"expirydate"`
	MarketLot       string `json:"marketlot"`
	Precision       string `json:"precision"`
	Multiplier      string `json:"multiplier"`
	TradeValue      string `json:"tradevalue"`
	TransactionType string `json:"transactiontype"`
	FillPrice       string `json:"fillprice"`
	FillSize        string `json:"fillsize"`
	OrderID         string `json:"orderid"`
	FillID          string `json:"fillid"`
	FillTime        string `json:"filltime"`
}

type LTPDataResponse struct {
	Status    bool       `json:"status"`
	Message   string     `json:"message"`
	ErrorCode string     `json:"errorcode"`
	Data      LtpDetails `json:"data"`
}

type LtpDetails struct {
	Exchange      string  `json:"exchange"`
	TradingSymbol string  `json:"tradingsymbol"`
	SymbolToken   string  `json:"symboltoken"`
	Open          float64 `json:"open"`
	High          float64 `json:"high"`
	Low           float64 `json:"low"`
	Close         float64 `json:"close"`
	Ltp           float64 `json:"ltp"`
}
