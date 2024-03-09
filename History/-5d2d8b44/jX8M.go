package breeze

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type ErrorResponse interface {
	GetError() *string
}

var (
	IntervalTypes = map[string]bool{
		"1minute":  true,
		"5minute":  true,
		"30minute": true,
		"1day":     true,
	}

	IntervalTypesHistV2 = map[string]bool{
		"1second":  true,
		"1minute":  true,
		"5minute":  true,
		"30minute": true,
		"1day":     true,
	}

	IntervalTypesStreamOHLC = map[string]bool{
		"1second":  true,
		"1minute":  true,
		"5minute":  true,
		"30minute": true,
	}

	ProductTypes = map[string]bool{
		"futures":    true,
		"options":    true,
		"futureplus": true,
		"optionplus": true,
		"cash":       true,
		"eatm":       true,
		"margin":     true,
		"mtf":        true,
		"btst":       true,
	}

	ProductTypesHist = map[string]bool{
		"futures":    true,
		"options":    true,
		"futureplus": true,
		"optionplus": true,
	}

	ProductTypesHistV2 = map[string]bool{
		"futures": true,
		"options": true,
		"cash":    true,
	}

	RightTypes = map[string]bool{
		"call":   true,
		"put":    true,
		"others": true,
	}

	ActionTypes = map[string]bool{
		"buy":  true,
		"sell": true,
	}

	OrderTypes = map[string]bool{
		"limit":    true,
		"market":   true,
		"stoploss": true,
	}

	ValidityTypes = map[string]bool{
		"day": true,
		"ioc": true,
		"vtc": true,
	}

	TransactionTypes = map[string]bool{
		"debit":  true,
		"credit": true,
	}

	ExchangeCodeHist = map[string]bool{
		"nse": true,
		"nfo": true,
	}

	ExchangeCodeHistV2 = map[string]bool{
		"nse": true,
		"bse": true,
		"nfo": true,
		"ndx": true,
		"mcx": true,
	}

	FnoExchangeTypes = map[string]bool{
		"nfo": true,
		"mcx": true,
		"ndx": true,
	}
)

type Funds struct {
	FundsData struct {
		BankAccount         string  `json:"bank_account"`
		TotalBankBalance    float64 `json:"total_bank_balance"`
		AllocatedEquity     float64 `json:"allocated_equity"`
		AllocatedFNO        float64 `json:"allocated_fno"`
		BlockByTradeEquity  float64 `json:"block_by_trade_equity"`
		BlockByTradeFNO     float64 `json:"block_by_trade_fno"`
		BlockByTradeBalance float64 `json:"block_by_trade_balance"`
		UnallocatedBalance  string  `json:"unallocated_balance"`
	} `json:"Success"`
	Status int     `json:"Status"`
	Error  *string `json:"Error"`
}

type DematHoldings struct {
	Success []DematHolding `json:"Success"`
	Status  int            `json:"Status"`
	Error   *string        `json:"Error"`
}

type DematHolding struct {
	StockCode              string `json:"stock_code"`
	StockISIN              string `json:"stock_ISIN"`
	Quantity               string `json:"quantity"`
	DematTotalBulkQuantity string `json:"demat_total_bulk_quantity"`
	DematAvailQuantity     string `json:"demat_avail_quantity"`
	BlockedQuantity        string `json:"blocked_quantity"`
	DematAllocatedQuantity string `json:"demat_allocated_quantity"`
}

type PorfolioPositions struct {
	PortfolioPosition []PorfolioPosition `json:"Success"`
	Status            int                `json:"Status"`
	Error             *string            `json:"Error"`
}

type PorfolioPosition struct {
	StockCode             string  `json:"stock_code"`
	ExchangeCode          string  `json:"exchange_code"`
	Quantity              string  `json:"quantity"`
	AveragePrice          string  `json:"average_price"`
	BookedProfitLoss      float64 `json:"booked_profit_loss"`
	CurrentMarketPrice    float64 `json:"current_market_price"`
	ChangePercentage      float64 `json:"change_percentage"`
	AnswerFlag            string  `json:"answer_flag"`
	ProductType           string  `json:"product_type"`
	ExpiryDate            string  `json:"expiry_date"`
	StrikePrice           string  `json:"strike_price"`
	Right                 string  `json:"right"`
	CategoryIndexPerStock string  `json:"category_index_per_stock"`
	Action                string  `json:"action"`
	RealizedProfit        float64 `json:"realized_profit"`
	UnrealizedProfit      float64 `json:"unrealized_profit"`
	OpenPositionValue     float64 `json:"open_position_value"`
	PortfolioCharges      float64 `json:"portfolio_charges"`
	CoverQuantity         string  `json:"cover_quantity"`
	StoplossTrigger       string  `json:"stoploss_trigger"`
}

type Quotes struct {
	Quote  []Quote `json:"Success"`
	Status int     `json:"Status"`
	Error  *string `json:"Error"`
}

type Quote struct {
	ExchangeCode        string      `json:"exchange_code"`
	ProductType         string      `json:"product_type"`
	StockCode           string      `json:"stock_code"`
	ExpiryDate          interface{} `json:"expiry_date"`
	Right               interface{} `json:"right"`
	StrikePrice         float64     `json:"strike_price"`
	LTP                 float64     `json:"ltp"`
	LTT                 string      `json:"ltt"`
	BestBidPrice        float64     `json:"best_bid_price"`
	BestBidQuantity     string      `json:"best_bid_quantity"`
	BestOfferPrice      float64     `json:"best_offer_price"`
	BestOfferQuantity   string      `json:"best_offer_quantity"`
	Open                float64     `json:"open"`
	High                float64     `json:"high"`
	Low                 float64     `json:"low"`
	PreviousClose       float64     `json:"previous_close"`
	LTPPercentChange    float64     `json:"ltp_percent_change"`
	UpperCircuit        float64     `json:"upper_circuit"`
	LowerCircuit        float64     `json:"lower_circuit"`
	TotalQuantityTraded string      `json:"total_quantity_traded"`
	SpotPrice           interface{} `json:"spot_price"`
}

type PortfolioHolding struct {
	StockCode             string  `json:"stock_code"`
	ExchangeCode          string  `json:"exchange_code"`
	Quantity              string  `json:"quantity"`
	AveragePrice          string  `json:"average_price"`
	BookedProfitLoss      string  `json:"booked_profit_loss"`
	CurrentMarketPrice    string  `json:"current_market_price"`
	ChangePercentage      string  `json:"change_percentage"`
	AnswerFlag            string  `json:"answer_flag"`
	ProductType           *string `json:"product_type"`             // Use pointer since it can be null
	ExpiryDate            *string `json:"expiry_date"`              // Use pointer since it can be null
	StrikePrice           *string `json:"strike_price"`             // Use pointer since it can be null
	Right                 *string `json:"right"`                    // Use pointer since it can be null
	CategoryIndexPerStock *string `json:"category_index_per_stock"` // Use pointer since it can be null
	Action                *string `json:"action"`                   // Use pointer since it can be null
	RealizedProfit        *string `json:"realized_profit"`          // Use pointer since it can be null
	UnrealizedProfit      *string `json:"unrealized_profit"`        // Use pointer since it can be null
	OpenPositionValue     *string `json:"open_position_value"`      // Use pointer since it can be null
	PortfolioCharges      *string `json:"portfolio_charges"`        // Use pointer since it can be null
}

type PortfolioHoldings struct {
	PortfolioHolding []PortfolioHolding `json:"Success"`
	Status           int                `json:"Status"`
	Error            *string            `json:"Error"` // Use pointer since it can be null
}

type HistoricalChartsV2Params struct {
	Interval     string
	FromDate     string
	ToDate       string
	StockCode    string
	ExchangeCode string
	ProductType  string
	ExpiryDate   string
	Right        string
	StrikePrice  string
}

type HistoricalChartDataV2 struct {
	HistoricalChartEntry []HistoricalChartEntry `json:"Success"`
	Status               int                    `json:"Status"`
	Error                *string                `json:"Error"` // Use pointer since it can be null
}

type HistoricalChartEntry struct {
	Close        float64 `json:"close"`
	Datetime     string  `json:"datetime"`
	ExchangeCode string  `json:"exchange_code"`
	ExpiryDate   string  `json:"expiry_date"`
	High         float64 `json:"high"`
	Low          float64 `json:"low"`
	Open         float64 `json:"open"`
	OpenInterest int     `json:"open_interest"`
	ProductType  string  `json:"product_type"`
	Right        string  `json:"right"`
	StockCode    string  `json:"stock_code"`
	StrikePrice  string  `json:"strike_price"`
	Volume       int     `json:"volume"`
}

type TradeBookItem struct {
	MatchAccount         string `json:"match_account"`
	OrderTradeDate       string `json:"order_trade_date"`
	OrderStockCode       string `json:"order_stock_code"`
	OrderFlow            string `json:"order_flow"`
	OrderQuantity        string `json:"order_quantity"`
	OrderAvgExecutedRate string `json:"order_average_executed_rate"`
	OrderTransValue      string `json:"order_trans_value"`
	OrderBrokerage       string `json:"order_brokerage"`
	OrderProduct         string `json:"order_product"`
	OrderExchangeCode    string `json:"order_exchange_code"`
	OrderReference       string `json:"order_reference"`
	OrderSegmentCode     string `json:"order_segment_code"`
	OrderSettlement      string `json:"order_settlement"`
	DpID                 string `json:"dp_id"`
	ClientID             string `json:"client_id"`
	LTP                  string `json:"LTP"`
	OrderEATMWithheld    string `json:"order_eATM_withheld"`
	OrderCshWithheld     string `json:"order_csh_withheld"`
	OrderTotalTaxes      string `json:"order_total_taxes"`
	OrderType            string `json:"order_type"`
}

type TradeBookResponse struct {
	TradeBook []TradeBookItem `json:"trade_book"`
}

type TradeListResponse struct {
	TradeBookResponse TradeBookResponse `json:"Success"`
	Status            int               `json:"Status"`
	Error             *string           `json:"Error"` // Use pointer since it can be null
}

type TradeDetail struct {
	OrderSettlement          string `json:"order_settlement"`
	OrderExchangeTradeNumber string `json:"order_exchange_trade_number"`
	OrderExecutedQuantity    string `json:"order_executed_quantity"`
	OrderFlow                string `json:"order_flow"`
	OrderBrokerage           string `json:"order_brokerage"`
	OrderPureBrokerage       string `json:"order_pure_brokerage"`
	OrderTaxes               string `json:"order_taxes"`
	OrderEATMWithheld        string `json:"order_eATM_withheld"`
	OrderCashWithheld        string `json:"order_cash_withheld"`
	OrderExecutedRate        string `json:"order_executed_rate"`
	OrderStockCode           string `json:"order_stock_code"`
	OrderExchangeCode        string `json:"order_exchange_code"`
	MatchAccount             string `json:"match_account"`
	OrderTradeReference      string `json:"order_trade_reference"`
	OrderExchangeTradeTm     string `json:"order_exchange_trade_tm"`
	OrderSegmentDis          string `json:"order_segment_dis"`
	OrderSegmentCode         string `json:"order_segment_code"`
}

type TradeDetails struct {
	TradeDetail []TradeDetail `json:"Success"`
	Status      int           `json:"Status"`
	Error       *string       `json:"Error"` // Use pointer since it can be null
}

type BrokerageCharge struct {
	Brokerage                   float64 `json:"brokerage"`
	ExchangeTurnoverCharges     float64 `json:"exchange_turnover_charges"`
	StampDuty                   float64 `json:"stamp_duty"`
	STT                         float64 `json:"stt"`
	SEBICharges                 float64 `json:"sebi_charges"`
	GST                         float64 `json:"gst"`
	TotalTurnoverAndSEBICharges float64 `json:"total_turnover_and_sebi_charges"`
	TotalOtherCharges           float64 `json:"total_other_charges"`
	TotalBrokerage              float64 `json:"total_brokerage"`
}

type BrokerageCharges struct {
	BrokerageCharge BrokerageCharge `json:"Success"`
	Status          int             `json:"Status"`
	Error           *string         `json:"Error"` // Use pointer since it can be null
}

type OrderResponse struct {
	Order  Order   `json:"Success"`
	Status int     `json:"Status"`
	Error  *string `json:"Error"`
}

type Order struct {
	OrderID string  `json:"order_id"`
	Message *string `json:"message"`
}

func (o *OrderResponse) getOrderID() (string, error) {
	if o.Status == 200 && o.Order.Message != nil {
		pattern := `reference no (\d{8}[A-Z]\d{9})`

		regExp, err := regexp.Compile(pattern)
		if err != nil {
			log.Error("Error compiling regex:", err)
			return "", err
		}

		match := regExp.FindStringSubmatch(*o.Order.Message)
		if len(match) >= 2 {
			orderID := match[1]
			log.Error("Order ID:", orderID)
			return orderID, nil
		} else {
			log.Error("Order ID not found in the text.")
			return "", errors.New("could not get the order id from the message")
		}
	}

	return "", errors.New("http status not ok or message is nil")
}

func (b *Breeze) getStockScriptList() error {
	// Fetch CSV content
	response, err := http.Get(stockScriptCSVURL)
	if err != nil {
		return fmt.Errorf("http request failed for stockScriptCSVURL: %w", err)
	}
	defer response.Body.Close()

	// Decode CSV content and process rows
	decodedContent, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll failed for getStockScriptList: %w", err)
	}

	cr := csv.NewReader(strings.NewReader(string(decodedContent)))

	myList, err := cr.ReadAll()
	if err != nil {
		return fmt.Errorf("cr.ReadAll() failed for getStockScriptList: %w", err)
	}

	// Initialize the stock script and token script maps
	stockScriptDictList, tokenScriptDictList := getStockScriptMaps(myList)

	// Update the stock script and token script maps in the Breeze struct
	b.stockScriptDictList = stockScriptDictList
	b.tokenScriptDictList = tokenScriptDictList

	return nil
}

func getStockScriptMaps(myList [][]string) (map[string]map[string]string, map[string]map[string][]string) {
	stockScriptDictList := make(map[string]map[string]string)
	tokenScriptDictList := make(map[string]map[string][]string)

	for _, row := range myList {
		exchange := row[2]
		symbol := row[3]
		token := row[5]
		instrument := row[7]
		companyName := row[1]

		switch exchange {
		case "BSE":
			if stockScriptDictList["BSE"] == nil {
				stockScriptDictList["BSE"] = make(map[string]string)
				tokenScriptDictList["BSE"] = make(map[string][]string)
			}
			stockScriptDictList["BSE"][symbol] = token
			tokenScriptDictList["BSE"][token] = []string{symbol, companyName}
		case "NSE":
			if stockScriptDictList["NSE"] == nil {
				stockScriptDictList["NSE"] = make(map[string]string)
				tokenScriptDictList["NSE"] = make(map[string][]string)
			}
			stockScriptDictList["NSE"][symbol] = token
			tokenScriptDictList["NSE"][token] = []string{symbol, companyName}
		case "NDX":
			if stockScriptDictList["NDX"] == nil {
				stockScriptDictList["NDX"] = make(map[string]string)
				tokenScriptDictList["NDX"] = make(map[string][]string)
			}
			stockScriptDictList["NDX"][instrument] = token
			tokenScriptDictList["NDX"][token] = []string{instrument, companyName}
		case "MCX":
			if stockScriptDictList["MCX"] == nil {
				stockScriptDictList["MCX"] = make(map[string]string)
				tokenScriptDictList["MCX"] = make(map[string][]string)
			}
			stockScriptDictList["MCX"][instrument] = token
			tokenScriptDictList["MCX"][token] = []string{instrument, companyName}
		case "NFO":
			if stockScriptDictList["NFO"] == nil {
				stockScriptDictList["NFO"] = make(map[string]string)
				tokenScriptDictList["NFO"] = make(map[string][]string)
			}
			stockScriptDictList["NFO"][instrument] = token
			tokenScriptDictList["NFO"][token] = []string{instrument, companyName}
		}
	}

	return stockScriptDictList, tokenScriptDictList
}

func (a *ApificationBreeze) getCustomerDetails(sessionToken string, useProxy bool) (map[string]interface{}, error) {
	if sessionToken == "" {
		return a.validationErrorResponse("Empty sessionToken value received in getCustomerDetails"), nil
	}

	body := map[string]interface{}{
		"SessionToken": sessionToken,
		"AppKey":       a.Breeze.apiKey,
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		a.errorException("getCustomerDetails", err)
	}

	headers := a.generateHeaders(string(bodyJSON))

	responseBody, err := a.makeRequest(http.MethodGet, customerDetailsEndpoint, string(bodyJSON), headers, useProxy)
	if err != nil {
		a.errorException("getCustomerDetails", err)

		return nil, err
	}

	var response map[string]interface{}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		a.errorException("getCustomerDetails", err)

		return response, err
	}

	if success, ok := response["Success"].(map[string]interface{}); ok && success != nil {
		delete(success, "session_token")
		response["Success"] = success
	}

	return response, nil
}

func (a *ApificationBreeze) executeBreezeAPIRESTRequest(method, endpoint string, requestBody interface{}, useProxy bool) ([]byte, error) {
	if a.sessionToken == "" {
		return nil, fmt.Errorf("empty sessionToken value in makeBreezeRESTCall")
	}

	bodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	headers := a.generateHeaders(string(bodyJSON))
	responseBody, err := a.makeRequest(method, endpoint, string(bodyJSON), headers, useProxy)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (a *ApificationBreeze) getAPIResponse(method string, useProxy bool, endpoint string, requestBody interface{}, responseData interface{}, errorHandler func(error)) error {
	jsonResponse, err := a.executeBreezeAPIRESTRequest(method, endpoint, requestBody, useProxy)

	if err != nil {
		log.Errorf("Error while executing breeze api rest request for %s: %s", endpoint, err)
		errorHandler(err)
		return err
	}

	if err = json.Unmarshal(jsonResponse, responseData); err != nil {
		log.Errorf("Error parsing JSON for %s: %s", endpoint, err)
		errorHandler(err)
		return err
	}

	log.Debugf("%s: %+v", endpoint, responseData)
	return nil
}

func (a *ApificationBreeze) getFunds(useProxy bool) (funds Funds, err error) {
	requestBody := make(map[string]string)
	getFunds, err := a.executeBreezeAPIRESTRequest(http.MethodGet, fundsEndpoint, requestBody, useProxy)
	if err != nil {
		log.Errorf("Error while executing breeze api rest request for getfunds: %s", err)

		return
	}

	if err = json.Unmarshal(getFunds, &funds); err != nil {
		log.Errorf("Error parsing JSON for getFunds: %s", err)

		return
	}

	log.Debugf("Funds: %+v", funds)

	return
}

func (a *ApificationBreeze) getDematHoldings(useProxy bool) (dematHoldings DematHoldings, err error) {
	requestBody := make(map[string]string)
	dematHoldingsJsonResponse, err := a.executeBreezeAPIRESTRequest(http.MethodGet, dematHoldingsEndpoint, requestBody, useProxy)
	if err != nil {
		log.Errorf("Error while executing breeze api rest request for getDematHoldings: %s", err)

		return
	}

	if err = json.Unmarshal(dematHoldingsJsonResponse, &dematHoldings); err != nil {
		log.Errorf("Error parsing JSON for getDematHoldings: %s", err)

		return
	}

	if dematHoldings.Error != nil {
		log.Debugf("DematHoldings Error: %s", *dematHoldings.Error)
	}
	log.Debugf("DematHoldings: %+v", dematHoldings)

	return
}

func (a *ApificationBreeze) getPortfolioHoldings(useProxy bool, requestBody map[string]string) (porfolioHoldings PortfolioHoldings, err error) {
	porfolioHoldingsJsonResponse, err := a.executeBreezeAPIRESTRequest(http.MethodGet, portfolioHoldingsEndpoint, requestBody, useProxy)
	if err != nil {
		log.Errorf("Error while executing breeze api rest request for getPortfolioHoldings: %s", err)

		return
	}

	if err = json.Unmarshal(porfolioHoldingsJsonResponse, &porfolioHoldings); err != nil {
		log.Errorf("Error parsing JSON for getPortfolioHoldings: %s", err)

		return
	}

	if porfolioHoldings.Error != nil {
		log.Debugf("Portfolio Holdings Error: %s", *porfolioHoldings.Error)
	}

	log.Debugf("Portfolio Holdings: : %+v", porfolioHoldings)

	return
}

func (a *ApificationBreeze) getPortfolioPositions(useProxy bool) (porfolioPositions PorfolioPositions, err error) {
	requestBody := make(map[string]string)
	porfolioPositionsJsonResponse, err := a.executeBreezeAPIRESTRequest(http.MethodGet, portfolioPositionsEndpoint, requestBody, useProxy)
	log.Debugf("porfolioPositionsJsonResponse: %s", porfolioPositionsJsonResponse)
	if err != nil {
		log.Errorf("Error while executing breeze api rest request for getPortfolioPositions: %s", err)

		return
	}

	if err = json.Unmarshal(porfolioPositionsJsonResponse, &porfolioPositions); err != nil {
		log.Errorf("Error parsing JSON for getPortfolioPositions: %s", err)

		return
	}

	if porfolioPositions.Error != nil {
		log.Debugf("Portfolio Positions Error: %s", *porfolioPositions.Error)
	}

	log.Debugf("Portfolio Positions: : %+v", porfolioPositions)

	return
}

func (a *ApificationBreeze) getQuotes(useProxy bool, requestBody map[string]string) (quotes Quotes, err error) {
	errorHandler := func(err error) {
		log.Debugf("Quotes Error: %s", err)
	}

	err = a.getAPIResponse(http.MethodGet, useProxy, quotesEndpoint, requestBody, &quotes, errorHandler)

	return
}

func (a *ApificationBreeze) getTradeList(useProxy bool, requestBody map[string]string) (tradeListResponse TradeListResponse, err error) {
	errorHandler := func(err error) {
		log.Debugf("tradeListResponse Error: %s", err)
	}

	err = a.getAPIResponse(http.MethodGet, useProxy, tradesEndpoint, requestBody, &tradeListResponse, errorHandler)

	return
}

func (a *ApificationBreeze) getTradeDetail(useProxy bool, requestBody map[string]string) (tradeDetails TradeDetails, err error) {
	errorHandler := func(err error) {
		log.Debugf("TradeDetails Error: %s", err)
	}

	err = a.getAPIResponse(http.MethodGet, useProxy, tradesEndpoint, requestBody, &tradeDetails, errorHandler)

	return
}

func (a *ApificationBreeze) getBrokerageCharges(useProxy bool, requestBody map[string]string) (brokerageCharges BrokerageCharges, err error) {
	errorHandler := func(err error) {
		log.Debugf("BrokerageCharges Error: %s", err)
	}

	err = a.getAPIResponse(http.MethodGet, useProxy, brokerageEndpoint, requestBody, &brokerageCharges, errorHandler)

	return
}

func (a *ApificationBreeze) placeOrder(useProxy bool, requestBody interface{}) (orderResponse OrderResponse, err error) {
	errorHandler := func(err error) {
		log.Debugf("placeOrder Error: %s", err)
	}

	err = a.getAPIResponse(http.MethodPost, useProxy, orderEndpoint, requestBody, &orderResponse, errorHandler)

	return
}

func (a *ApificationBreeze) getHistoricalChartsV2DataStruct(historicalChartDataV2Params HistoricalChartsV2Params, useProxy bool) (historicalChartDataV2 HistoricalChartDataV2, err error) {
	return a.getHistoricalChartsV2Data(
		historicalChartDataV2Params.Interval,
		historicalChartDataV2Params.FromDate,
		historicalChartDataV2Params.ToDate,
		historicalChartDataV2Params.StockCode,
		historicalChartDataV2Params.ExchangeCode,
		historicalChartDataV2Params.ProductType,
		historicalChartDataV2Params.ExpiryDate,
		historicalChartDataV2Params.StrikePrice,
		historicalChartDataV2Params.Right,
		useProxy)
}

func (a *ApificationBreeze) getHistoricalChartsV2Data(interval, from_date, to_date, stock_code, exchange_code, product_type, expiry_date, strike_price, right string, useProxy bool) (historicalChartDataV2 HistoricalChartDataV2, err error) {
	urlParams := url.Values{
		"interval":   {interval},
		"from_date":  {from_date},
		"to_date":    {to_date},
		"stock_code": {stock_code},
		"exch_code":  {exchange_code},
	}

	if product_type != "" {
		urlParams.Add("product_type", product_type)
	}
	if expiry_date != "" {
		urlParams.Add("expiry_date", expiry_date)
	}
	if strike_price != "" {
		urlParams.Add("strike_price", strike_price)
	}
	if right != "" {
		urlParams.Add("right", right)
	}

	url := apiBaseURLV2 + historicalChartsV2Endpoint

	req, err := http.NewRequest(http.MethodGet, url+"?"+urlParams.Encode(), nil)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-SessionToken", a.base64SessionToken)
	req.Header.Set("apikey", a.apiKey)

	client := getClient(useProxy)
	historicalChartsV2DataResponse, err := client.Do(req)
	if err != nil {
		return
	}
	defer historicalChartsV2DataResponse.Body.Close()

	historicalChartsV2DataResponseBody, err := io.ReadAll(historicalChartsV2DataResponse.Body)
	if err != nil {
		log.Errorf("Error while executing breeze api rest request for getHistoricalChartsV2Data: %s", err)

		return
	}

	log.Debugf("responseBody: %s", historicalChartsV2DataResponseBody)

	if err = json.Unmarshal(historicalChartsV2DataResponseBody, &historicalChartDataV2); err != nil {
		log.Errorf("Error parsing JSON for getHistoricalChartsV2Data: %s", err)

		return
	}

	if historicalChartDataV2.Error != nil {
		log.Debugf("Historical Charts V2 Error: %s", *historicalChartDataV2.Error)
	}

	log.Debugf("Historical Charts V2 Error: %+v", historicalChartDataV2)

	return
}

func IsValidValue(value string, validValues map[string]bool) bool {
	return validValues[strings.ToLower(value)]
}

func createTradeListRequestBody(exchangeCode, fromDate, toDate, productType, action, stockCode string) map[string]string {
	quotesRequestBody := make(map[string]string)
	quotesRequestBody["exchange_code"] = exchangeCode

	if fromDate != "" {
		quotesRequestBody["from_date"] = fromDate
	}
	if toDate != "" {
		quotesRequestBody["to_date"] = toDate
	}
	if productType != "" {
		quotesRequestBody["product_type"] = productType
	}
	if action != "" {
		quotesRequestBody["action"] = action
	}
	if stockCode != "" {
		quotesRequestBody["stock_code"] = stockCode
	}

	return quotesRequestBody
}

func createTradeDetailRequestBody(exchangeCode, order_id string) map[string]string {
	quotesRequestBody := make(map[string]string)
	quotesRequestBody["exchange_code"] = exchangeCode
	quotesRequestBody["order_id"] = order_id

	return quotesRequestBody
}

func createBrokerageChargesRequestBody(stockCode, exchangeCode, product, orderType, price, action, quantity, expiryDate, right, strikePrice, stopLoss, orderRateFresh string) map[string]string {
	quotesRequestBody := make(map[string]string)
	quotesRequestBody["stock_code"] = stockCode
	quotesRequestBody["exchange_code"] = exchangeCode
	quotesRequestBody["product"] = product
	quotesRequestBody["order_type"] = orderType
	quotesRequestBody["price"] = price
	quotesRequestBody["action"] = action
	quotesRequestBody["quantity"] = quantity
	quotesRequestBody["expiry_date"] = expiryDate
	quotesRequestBody["right"] = right
	quotesRequestBody["strike_price"] = strikePrice
	quotesRequestBody["stoploss"] = stopLoss
	quotesRequestBody["order_rate_fresh"] = orderRateFresh

	return quotesRequestBody
}

func createOrderRequestBody(stockCode, exchangeCode, productType, action, orderType string, stopLoss float64, quantity int, price float64, validity string, validityDate string, disclosedQuantity int, expiryDate string, right string, strikePrice int, userRemark string) (interface{}, error) {
	orderRequestBody := make(map[string]interface{})

	if stockCode == "" || exchangeCode == "" || productType == "" || action == "" || orderType == "" || validity == "" {
		return orderRequestBody, errors.New("validation error: Missing required fields")
	}

	if !IsValidValue(productType, ProductTypes) {
		return orderRequestBody, errors.New("validation error: invalid product type")
	}

	if !IsValidValue(action, ActionTypes) {
		return orderRequestBody, errors.New("validation error: invalid action type")
	}

	if !IsValidValue(orderType, OrderTypes) {
		return orderRequestBody, errors.New("validation error: invalid order type")
	}

	if !IsValidValue(validity, ValidityTypes) {
		return orderRequestBody, errors.New("validation error: invalid validity type")
	}

	orderRequestBody["stock_code"] = stockCode
	orderRequestBody["exchange_code"] = exchangeCode
	orderRequestBody["product"] = productType
	orderRequestBody["action"] = action
	orderRequestBody["order_type"] = orderType

	if stopLoss != 0.0 {
		orderRequestBody["stoploss"] = stopLoss
	}
	orderRequestBody["quantity"] = quantity
	orderRequestBody["price"] = price
	orderRequestBody["validity"] = validity

	if validityDate != "" {
		orderRequestBody["validity_date"] = validityDate
	}
	if disclosedQuantity != 0 {
		orderRequestBody["disclosed_quantity"] = disclosedQuantity
	}
	if expiryDate != "" {
		orderRequestBody["expiry_date"] = expiryDate
	}
	if right != "" {
		orderRequestBody["right"] = right
	}
	if strikePrice != 0 {
		orderRequestBody["strike_price"] = strikePrice
	}
	if userRemark != "" {
		orderRequestBody["user_remark"] = userRemark
	}

	return orderRequestBody, nil
}
