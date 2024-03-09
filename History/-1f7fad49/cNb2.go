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
	OptionType         string `json:"option_type"`
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
	ExchangeTimestamp            time.Time
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
	LastTradedTimestamp          time.Time
	OpenInterest                 int64
	OpenInterestChangePercentage int64
	UpperCircuitLimit            int64
	LowerCircuitLimit            int64
	Week52HighPrice              int64
	Week52LowPrice               int64
	LastTradedPriceFloat         float64
	Best5BuyData, Best5SellData  []int64 // Modify the data types as required
}

type BestFivePacket struct {
	BuySellFlag int16
	Quantity    int64
	Price       int64
	NumOfOrders int16
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

type HoldingsResponse struct {
	Status    bool   `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"errorcode"`
	Data      []struct {
		TradingSymbol      string      `json:"tradingsymbol"`
		Exchange           string      `json:"exchange"`
		ISIN               string      `json:"isin"`
		T1Quantity         float64     `json:"t1quantity"`
		RealisedQuantity   float64     `json:"realisedquantity"`
		Quantity           float64     `json:"quantity"`
		AuthorisedQuantity float64     `json:"authorisedquantity"`
		Product            string      `json:"product"`
		CollateralQuantity interface{} `json:"collateralquantity"` // Use interface{} to handle null values
		CollateralType     interface{} `json:"collateraltype"`     // Use interface{} to handle null values
		Haircut            float64     `json:"haircut"`
		AveragePrice       float64     `json:"averageprice"`
		LTP                float64     `json:"ltp"`
		Close              float64     `json:"close"`
		ProfitAndLoss      float64     `json:"profitandloss"`
		PNLPercentage      float64     `json:"pnlpercentage"`
	} `json:"data"`
}

type PositionsResponse struct {
	Status    bool       `json:"status"`
	Message   string     `json:"message"`
	ErrorCode string     `json:"errorcode"`
	Positions []Position `json:"data"`
}

type Position struct {
	Exchange          string `json:"exchange"`
	SymbolToken       string `json:"symboltoken"`
	ProductType       string `json:"producttype"`
	TradingSymbol     string `json:"tradingsymbol"`
	SymbolName        string `json:"symbolname"`
	InstrumentType    string `json:"instrumenttype"`
	PriceDen          string `json:"priceden"`
	PriceNum          string `json:"pricenum"`
	GenDen            string `json:"genden"`
	GenNum            string `json:"gennum"`
	Precision         string `json:"precision"`
	Multiplier        string `json:"multiplier"`
	BoardLotSize      string `json:"boardlotsize"`
	BuyQty            string `json:"buyqty"`
	SellQty           string `json:"sellqty"`
	BuyAmount         string `json:"buyamount"`
	SellAmount        string `json:"sellamount"`
	SymbolGroup       string `json:"symbolgroup"`
	StrikePrice       string `json:"strikeprice"`
	OptionType        string `json:"optiontype"`
	ExpiryDate        string `json:"expirydate"`
	LotSize           string `json:"lotsize"`
	CFBuyQty          string `json:"cfbuyqty"`
	CFSellQty         string `json:"cfsellqty"`
	CFBuyAmount       string `json:"cfbuyamount"`
	CFSellAmount      string `json:"cfsellamount"`
	BuyAvgPrice       string `json:"buyavgprice"`
	SellAvgPrice      string `json:"sellavgprice"`
	AvgNetPrice       string `json:"avgnetprice"`
	NetValue          string `json:"netvalue"`
	NetQty            string `json:"netqty"`
	TotalBuyValue     string `json:"totalbuyvalue"`
	TotalSellValue    string `json:"totalsellvalue"`
	CFBuyAvgPrice     string `json:"cfbuyavgprice"`
	CFSellAvgPrice    string `json:"cfsellavgprice"`
	TotalBuyAvgPrice  string `json:"totalbuyavgprice"`
	TotalSellAvgPrice string `json:"totalsellavgprice"`
	NetPrice          string `json:"netprice"`
}

type FundsAndMargin struct {
	Status    bool              `json:"status"`
	Message   string            `json:"message"`
	ErrorCode string            `json:"errorcode"`
	Data      FundAndMarginData `json:"data"`
}

type FundAndMarginData struct {
	Net                    float64  `json:"net,string"`
	AvailableCash          float64  `json:"availablecash,string"`
	AvailableIntradayPayin float64  `json:"availableintradaypayin,string"`
	AvailableLimitMargin   float64  `json:"availablelimitmargin,string"`
	Collateral             float64  `json:"collateral,string"`
	M2MUnrealized          float64  `json:"m2munrealized,string"`
	M2MRealized            float64  `json:"m2mrealized,string"`
	UtilisedDebits         float64  `json:"utiliseddebits,string"`
	UtilisedSpan           *float64 `json:"utilisedspan,string"`
	UtilisedOptionPremium  *float64 `json:"utilisedoptionpremium,string"`
	UtilisedHoldingSales   *float64 `json:"utilisedholdingsales,string"`
	UtilisedExposure       *float64 `json:"utilisedexposure,string"`
	UtilisedTurnover       *float64 `json:"utilisedturnover,string"`
	UtilisedPayout         float64  `json:"utilisedpayout,string"`
}

type InstrumentSymbolKey struct {
	Symbol   string
	Exchange string
}

type InstrumentTokenKey struct {
	Token    string
	Exchange string
}
