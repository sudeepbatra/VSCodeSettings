package smartapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sudeepbatra/alpha-hft/logger"
)

var (
	ErrAccessTokenEmpty        = errors.New("access token is empty")
	ErrOrderTypeNotValid       = errors.New("order type is not valid")
	ErrTransactionTypeNotValid = errors.New("transaction type is not valid")
	ErrProductTypeNotValid     = errors.New("product type is not valid")
)

const (
	InstrumentScripURL  = "https://margincalculator.angelbroking.com/OpenAPI_File/files/OpenAPIScripMaster.json"
	HistoricalPricesURL = "https://apiconnect.angelbroking.com/rest/secure/angelbroking/historical/v1/getCandleData"
	placeOrderURL       = "https://apiconnect.angelbroking.com/rest/secure/angelbroking/order/v1/placeOrder"
	modifyOrderURL      = "https://apiconnect.angelbroking.com/rest/secure/angelbroking/order/v1/modifyOrder"
	cancelOrderURL      = "https://apiconnect.angelbroking.com/rest/secure/angelbroking/order/v1/cancelOrder"
	orderBookURL        = "https://apiconnect.angelbroking.com/rest/secure/angelbroking/order/v1/getOrderBook"
)

func (s *SmartApiApplication) CheckAccessTokenNotEmpty() error {
	if *s.AccessToken == "" {
		logger.Log.Error().Msg("access token is empty")
		return ErrAccessTokenEmpty
	}

	return nil
}

func (s *SmartApiApplication) setDefaultHeaders(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *s.AccessToken))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-UserType", "USER")
	req.Header.Add("X-SourceID", "WEB")
	req.Header.Add("X-ClientLocalIP", "CLIENT_LOCAL_IP")
	req.Header.Add("X-ClientPublicIP", "CLIENT_PUBLIC_IP")
	req.Header.Add("X-MACAddress", "MAC_ADDRESS")
	req.Header.Add("X-PrivateKey", s.AppKey)

	logger.Log.Debug().Interface("headers", req.Header).Msg("Request Headers")
}

func GetInstrumentsScripMasterData(client *http.Client) ([]InstrumentRecord, error) {
	req, err := http.NewRequest(http.MethodGet, InstrumentScripURL, nil)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in GetInstrumentsScripMasterData")
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in GetInstrumentsScripMasterData")
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in GetInstrumentsScripMasterData")
		return nil, err
	}

	var instrumentRecords []InstrumentRecord

	err = json.Unmarshal(body, &instrumentRecords)

	if err != nil {
		logger.Log.Error().Err(err).Msg("error in deocding get instrument scrips response")
		return nil, err
	}

	for i := range instrumentRecords {
		if instrumentRecords[i].ExchSeg == "NSE" {
			symbolParts := strings.Split(instrumentRecords[i].Symbol, "-")
			if len(symbolParts) > 1 {
				instrumentRecords[i].InstrumentTypeCode = symbolParts[1]
			} else {
				instrumentRecords[i].InstrumentTypeCode = ""
			}
		}
	}

	return instrumentRecords, nil
}

func (s *SmartApiApplication) GetCandleDataForIntervalRange(client *http.Client,
	token, exchange, interval, from_date, to_date string,
) ([]CandleData, error) {
	if *s.AccessToken == "" {
		logger.Log.Error().Msg("access token is empty")
		return nil, ErrAccessTokenEmpty
	}

	method := "POST"

	payload := map[string]string{
		"exchange":    exchange,
		"symboltoken": token,
		"interval":    strings.ToUpper(interval),
		"fromdate":    from_date,
		"todate":      to_date,
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		logger.Log.Error().Str("endpoint", "historical/v1/getCandleData").Err(err).Msg("error in marshaling request body")
		return nil, err
	}

	logger.Log.Info().Str("body", string(requestBody)).Msg("payload for historic data")

	logger.Log.Debug().Msg("Request URL:" + HistoricalPricesURL)
	logger.Log.Debug().Msg("Request Method:" + method)
	logger.Log.Debug().Msg("Request Body:" + string(requestBody))

	req, err := http.NewRequest(method, HistoricalPricesURL, bytes.NewBuffer(requestBody))
	if err != nil {
		logger.Log.Error().Str("endpoint", "historical/v1/getCandleData").Err(err).Msg("error in construct request")
		return nil, err
	}

	req.Header.Add("X-PrivateKey", s.AppKey)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-SourceID", "WEB")
	req.Header.Add("X-ClientLocalIP", "CLIENT_LOCAL_IP")
	req.Header.Add("X-ClientPublicIP", "CLIENT_PUBLIC_IP")
	req.Header.Add("X-MACAddress", "MAC_ADDRESS")
	req.Header.Add("X-UserType", "USER")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *s.AccessToken))
	req.Header.Add("Content-Type", "application/json")

	logger.Log.Debug().Interface("headers", req.Header).Msg("Request Headers")

	res, err := client.Do(req)
	if err != nil {
		logger.Log.Error().Str("endpoint", "historical/v1/getCandleData").Err(err).Msg("error in tirggering request")
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Log.Error().Str("endpoint", "historical/v1/getCandleData").Err(err).Msg("error in decoding response")
		return nil, err
	}

	logger.Log.Debug().Str("data", string(body)).Msg("historical result")

	return parseHistoricalCandleData(err, body, token, interval, exchange)
}

func parseHistoricalCandleData(err error, body []byte, token, interval, exchange string) ([]CandleData, error) {
	var historicalAPIResponse HistoricalApiResponse

	parseErr := json.Unmarshal(body, &historicalAPIResponse)
	if parseErr != nil {
		logger.Log.Error().Str("endpoint", "historical/v1/getCandleData").Err(parseErr).Msg("error in parsing response")
		return nil, parseErr
	}

	if historicalAPIResponse.Data == nil {
		logger.Log.Info().Str("symbolToken", token).Str("interval", interval).Msg("unable to fetch historical price")
	}

	var candlesData []CandleData

	for _, candlePoint := range historicalAPIResponse.Data {
		if len(candlePoint) == 6 {
			parsedTime, _ := time.Parse(time.RFC3339, candlePoint[0].(string))
			parsedToken, _ := strconv.Atoi(token)

			candle := CandleData{
				Token:     parsedToken,
				Exchange:  exchange,
				Timestamp: parsedTime,
				Open:      candlePoint[1].(float64),
				High:      candlePoint[2].(float64),
				Low:       candlePoint[3].(float64),
				Close:     candlePoint[4].(float64),
				Volume:    int(candlePoint[5].(float64)),
			}
			candlesData = append(candlesData, candle)
		}
	}

	return candlesData, nil
}

func (s *SmartApiApplication) requestSmartAPI(client *http.Client,
	httpMethod, requestURL string,
	queryParams map[string]string,
	payload interface{},
	headers map[string]string,
) ([]byte, error) {
	req, err := http.NewRequest(httpMethod, requestURL, nil)
	if err != nil {
		return nil, err
	}

	if queryParams != nil {
		q := req.URL.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}

		req.URL.RawQuery = q.Encode()
	}

	s.setDefaultHeaders(req)

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	if payload != nil {
		requestBody, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		req.Body = io.NopCloser(bytes.NewReader(requestBody))
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *SmartApiApplication) PlaceOrder(client *http.Client, variety, token, tradingSymbol, exchange,
	transactionType, orderType, productType, duration, price, quantity, squareOff, stopLoss string,
) (*PlaceOrderResponse, error) {
	err := s.CheckAccessTokenNotEmpty()
	if err != nil {
		return nil, err
	}

	validatorsWithTransactionType := []func(string) error{
		validateVariety,
		validateTransactionType,
		validateOrderType,
		validateProductType,
		validateDuration,
		validateExchange,
	}

	err = validateOrderValues([]string{variety, transactionType, orderType, productType, duration, exchange}, validatorsWithTransactionType)
	if err != nil {
		return nil, err
	}

	httpMethod := "POST"

	orderRequestPayload := map[string]interface{}{
		"variety":         variety,
		"symboltoken":     token,
		"tradingsymbol":   tradingSymbol,
		"transactiontype": transactionType,
		"exchange":        exchange,
		"ordertype":       orderType,
		"producttype":     productType,
		"duration":        duration,
		"price":           price,
		"squareoff":       squareOff,
		"stoploss":        stopLoss,
		"quantity":        quantity,
	}

	resBody, err := s.requestSmartAPI(client, httpMethod, placeOrderURL, nil, orderRequestPayload, nil)
	if err != nil {
		return nil, err
	}

	logger.Log.Debug().Str("data", string(resBody)).Msg("place order resBody")

	var placeOrderResponse PlaceOrderResponse

	err = json.Unmarshal(resBody, &placeOrderResponse)
	if err != nil {
		return nil, err
	}

	return &placeOrderResponse, nil
}

func (s *SmartApiApplication) ModifyOrder(client *http.Client, variety, token, tradingSymbol, exchange,
	orderID, orderType, productType, duration, price, quantity string,
) (*ModifyOrderResponse, error) {
	err := s.CheckAccessTokenNotEmpty()
	if err != nil {
		return nil, err
	}

	validators := []func(string) error{
		validateVariety,
		validateOrderType,
		validateProductType,
		validateDuration,
		validateExchange,
	}

	err = validateOrderValues([]string{variety, orderType, productType, duration, exchange}, validators)
	if err != nil {
		return nil, err
	}

	httpMethod := "POST"

	orderRequestPayload := map[string]interface{}{
		"variety":       variety,
		"orderid":       orderID,
		"ordertype":     orderType,
		"producttype":   productType,
		"duration":      duration,
		"price":         price,
		"quantity":      quantity,
		"symboltoken":   token,
		"tradingsymbol": tradingSymbol,
		"exchange":      exchange,
	}

	resBody, err := s.requestSmartAPI(client, httpMethod, modifyOrderURL, nil, orderRequestPayload, nil)
	if err != nil {
		return nil, err
	}

	logger.Log.Debug().Str("data", string(resBody)).Msg("modify order resBody")

	var modifyOrderResponse ModifyOrderResponse

	err = json.Unmarshal(resBody, &modifyOrderResponse)
	if err != nil {
		return nil, err
	}

	return &modifyOrderResponse, nil
}

func (s *SmartApiApplication) CancelOrder(client *http.Client, variety, orderID string) (*CancelOrderResponse, error) {
	err := s.CheckAccessTokenNotEmpty()
	if err != nil {
		return nil, err
	}

	validators := []func(string) error{
		validateVariety,
	}

	err = validateOrderValues([]string{variety}, validators)
	if err != nil {
		return nil, err
	}

	httpMethod := "POST"

	orderRequestPayload := map[string]interface{}{
		"variety": variety,
		"orderid": orderID,
	}

	resBody, err := s.requestSmartAPI(client, httpMethod, cancelOrderURL, nil, orderRequestPayload, nil)
	if err != nil {
		return nil, err
	}

	logger.Log.Debug().Str("data", string(resBody)).Msg("cancel order result")

	var cancelOrderResponse CancelOrderResponse

	err = json.Unmarshal(resBody, &cancelOrderResponse)
	if err != nil {
		return nil, err
	}

	return &cancelOrderResponse, nil
}

func (s *SmartApiApplication) GetOrderBook(client *http.Client) (*OrderBookResponse, error) {
	err := s.CheckAccessTokenNotEmpty()
	if err != nil {
		return nil, err
	}

	httpMethod := "GET"

	resBody, err := s.requestSmartAPI(client, httpMethod, orderBookURL, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	logger.Log.Debug().Str("data", string(resBody)).Msg("get order book result")

	var orderBookResponse OrderBookResponse

	err = json.Unmarshal(resBody, &orderBookResponse)
	if err != nil {
		return nil, err
	}

	return &orderBookResponse, nil
}

func validateOrderValues(values []string, validators []func(string) error) error {
	if len(values) != len(validators) {
		return errors.New("number of values and validators must match")
	}

	for i, v := range values {
		if err := validators[i](v); err != nil {
			return err
		}
	}

	return nil
}

func validateVariety(variety string) error {
	if strings.ToUpper(variety) != "NORMAL" &&
		strings.ToUpper(variety) != "STOPLOSS" &&
		strings.ToUpper(variety) != "AMO" &&
		strings.ToUpper(variety) != "ROBO" {
		logger.Log.Error().
			Str("variety", variety).
			Msg("variety is not valid. Only variety of NORMAL, STOPLOSS, AMO and ROBO supported.")

		return ErrOrderTypeNotValid
	}

	return nil
}

func validateTransactionType(transactionType string) error {
	if strings.ToUpper(transactionType) != "BUY" &&
		strings.ToUpper(transactionType) != "SELL" {
		logger.Log.Error().
			Str("transactionType", transactionType).
			Msg("transaction type is not valid. Only transaction types of BUY and SELL supported.")

		return ErrTransactionTypeNotValid
	}

	return nil
}

func validateOrderType(orderType string) error {
	if strings.ToUpper(orderType) != "MARKET" &&
		strings.ToUpper(orderType) != "LIMIT" &&
		strings.ToUpper(orderType) != "STOPLOSS_LIMIT" &&
		strings.ToUpper(orderType) != "STOPLOSS_MARKET" {
		logger.Log.Error().
			Str("orderType", orderType).
			Msg("order type is not valid. Only order types of MARKET, LIMIT, STOPLOSS_LIMIT and STOPLOSS_MARKET supported.")

		return ErrOrderTypeNotValid
	}

	return nil
}

func validateProductType(productType string) error {
	if strings.ToUpper(productType) != "DELIVERY" &&
		strings.ToUpper(productType) != "CARRYFORWARD" &&
		strings.ToUpper(productType) != "MARGIN" &&
		strings.ToUpper(productType) != "INTRADAY" &&
		strings.ToUpper(productType) != "BO" {
		logger.Log.Error().
			Str("productType", productType).
			Msg("product type is not valid. Only product types of INTRADAY and CNC supported.")

		return ErrOrderTypeNotValid
	}

	return nil
}

func validateDuration(duration string) error {
	if strings.ToUpper(duration) != "DAY" &&
		strings.ToUpper(duration) != "IOC" {
		logger.Log.Error().
			Str("duration", duration).
			Msg("duration is not valid. Only duration of DAY and IOC supported.")

		return ErrOrderTypeNotValid
	}

	return nil
}

func validateExchange(exchange string) error {
	if strings.ToUpper(exchange) != "BSE" &&
		strings.ToUpper(exchange) != "NSE" &&
		strings.ToUpper(exchange) != "NFO" &&
		strings.ToUpper(exchange) != "MCX" {
		logger.Log.Error().
			Str("exchange", exchange).
			Msg("exchange is not valid. Only exchange type of BSE, NSE, NFO and MCX supported.")

		return ErrOrderTypeNotValid
	}

	return nil
}
