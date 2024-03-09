package breeze

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"

	"github.com/sudeepbatra/alpha-hft/config"
)

const (
	stockScriptCSVURL          = "https://traderweb.icicidirect.com/Content/File/txtFile/ScripFile/StockScriptNew.csv"
	apiBaseURLV1               = "https://api.icicidirect.com/breezeapi/api/v1/"
	apiBaseURLV2               = "https://breezeapi.icicidirect.com/api/v2/"
	customerDetailsEndpoint    = "customerdetails"
	dematHoldingsEndpoint      = "dematholdings"
	fundsEndpoint              = "funds"
	portfolioHoldingsEndpoint  = "portfolioholdings"
	portfolioPositionsEndpoint = "portfoliopositions"
	quotesEndpoint             = "quotes"
	historicalChartsV2Endpoint = "historicalcharts"
	tradesEndpoint             = "trades"
	brokerageEndpoint          = "preview_order"
)

func generateSession(appName, apiKey, secretKey, sessionToken string) (*ApificationBreeze, error) {
	breezeInstance := &Breeze{
		appName:      appName,
		apiKey:       apiKey,
		sessionToken: sessionToken,
		secretKey:    secretKey,
	}

	err := breezeInstance.retrieveSessionToken()
	if err != nil {
		log.Error("Error while trying to retrieve session token. Unable to proceed further.")
		return nil, errors.New("error while trying to retrieve session token")
	}

	err = breezeInstance.getStockScriptList()
	if err != nil {
		log.Error("Error while trying to stock script list. Unable to proceed further.")
		return nil, errors.New("error while trying to retrieve stock script list")
	}

	apiHandler := &ApificationBreeze{
		Breeze:   breezeInstance,
		hostname: apiBaseURLV1,
	}

	apiHandler.base64SessionToken = base64.StdEncoding.EncodeToString([]byte(breezeInstance.userID + ":" + breezeInstance.sessionToken))
	log.Debug("Setting the base64SessionToken to:", apiHandler.base64SessionToken)

	return apiHandler, nil
}

func Login(sessionToken string, useProxy bool) {
	appName := config.Config.BrokerConfig.Breeze.Name
	apiKey := config.Config.BrokerConfig.Breeze.Key
	secretKey := config.Config.BrokerConfig.Breeze.Secret

	if sessionToken == "" {
		loginURL := fmt.Sprintf("https://api.icicidirect.com/apiuser/login?api_key=%s", url.QueryEscape(apiKey))

		log.Info("Use the Url to generate the sessionToken and pass it")
		log.Info("=======================================")
		log.Info(loginURL)
		log.Info("=======================================")

		return
	}

	apiHandler, err := generateSession(appName, apiKey, secretKey, sessionToken)
	if err != nil {
		log.Error("Error creating API handler:", err)

		return
	}

	customerDetails, err := apiHandler.getCustomerDetails(sessionToken, useProxy)
	if err != nil {
		log.Error("Error getting customer details:", err)

		return
	}

	funds, err := apiHandler.getFunds(useProxy)
	if err != nil {
		log.Error("Error while trying to getFunds:", err)

		return
	}

	dematHoldings, err := apiHandler.getDematHoldings(useProxy)
	if err != nil {
		log.Error("Error while trying to getDematHoldings:", err)

		return
	}

	portfolioHoldingsRequestBody := createSamplePorfolioHoldingsRequestBody()
	porfolioHoldings, err := apiHandler.getPortfolioHoldings(useProxy, portfolioHoldingsRequestBody)
	if err != nil {
		log.Error("Error while trying to getPortfolioHoldings:", err)

		return
	}

	porfolioPositions, err := apiHandler.getPortfolioPositions(useProxy)
	if err != nil {
		log.Error("Error while trying to getPortfolioPositions:", err)

		return
	}

	quotesRequestBody := createSampleQuotesRequestBody()

	quotes, err := apiHandler.getQuotes(useProxy, quotesRequestBody)
	if err != nil {
		log.Error("Error while trying to execute getQuotes:", err)

		return
	}

	historicalChartsV2Data, err := apiHandler.getHistoricalChartsV2Data("30minute", "2023-08-01T06:00:00.000Z", "2023-08-04T06:00:00.000Z", "TATMOT", "NSE", "", "", "", "", false)
	if err != nil {
		log.Error("Error while trying to execute historicalChartsV2Data:", err)

		return
	}

	historicalChartsV2Params := HistoricalChartsV2Params{
		Interval:     "5minute",
		FromDate:     "2023-08-01T06:00:00.000Z",
		ToDate:       "2023-08-04T06:00:00.000Z",
		StockCode:    "ZOMLIM",
		ExchangeCode: "NSE",
	}

	historicalChartsV2DataStruct, err := apiHandler.getHistoricalChartsV2DataStruct(historicalChartsV2Params, false)
	if err != nil {
		log.Error("Error while trying to execute getHistoricalChartsV2DataStruct:", err)

		return
	}

	tradeListRequestBody := createTradeListRequestBody("NSE", "", "", "", "", "")
	tradeList, err := apiHandler.getTradeList(useProxy, tradeListRequestBody)
	if err != nil {
		log.Error("Error while trying to execute getTradeList:", err)

		return
	}

	tradeDetailRequestBody := createTradeDetailRequestBody("NSE", "34u203u40")
	tradeDetail, err := apiHandler.getTradeDetail(useProxy, tradeDetailRequestBody)
	if err != nil {
		log.Error("Error while trying to execute getTradeDetail:", err)

		return
	}

	brokerageChargesRequestBody := createBrokerageChargesRequestBody("ZOMLIM", "NSE", "cash", "limit", "80", "buy", "100", "", "", "", "", "")
	brokerageCharges, err := apiHandler.getBrokerageCharges(useProxy, brokerageChargesRequestBody)
	if err != nil {
		log.Error("Error while trying to execute brokerageCharges:", err)

		return
	}

	orderRequestBody := createOrderRequestBody("AXIBAN", "NSE", "futures", "buy", "limit", 0.0, 10, 1500.0, "day", "", 0, "", "", 0, "Sample remark")
	breezeOrder, err := apiHandler.placeOrder(useProxy, orderRequestBody)
	if err != nil {
		log.Error("Error while trying to execute brokerageCharges:", err)

		return
	}

	log.Println("-----000----- Stock Script Dict List -----000----- ")
	log.Println(apiHandler.stockScriptDictList)
	log.Println("-----000----- Stock Script Dict List -----000----- ")
	log.Println("-----^^^----- Token Script Dict List -----^^^----- ")
	log.Println(apiHandler.tokenScriptDictList)
	log.Println("-----^^^----- Token Script Dict List -----^^^----- ")

	log.Info("Customer details:", customerDetails)
	log.Infof("Funds: %+v", funds)

	log.Infof("DematHoldings: %v", dematHoldings)
	if dematHoldings.Error != nil {
		log.Infof("DematHoldings Error: %s", *dematHoldings.Error)
	}

	log.Infof("porfolioHoldings: %v", porfolioHoldings)
	if porfolioHoldings.Error != nil {
		log.Infof("porfolioHoldings Error: %s", *porfolioHoldings.Error)
	}

	log.Infof("porfolioPositions: %v", porfolioPositions)
	if porfolioPositions.Error != nil {
		log.Infof("porfolioPositions Error: %s", *porfolioPositions.Error)
	}

	log.Infof("quotes: %v", quotes)
	if quotes.Error != nil {
		log.Infof("quotes Error: %s", *quotes.Error)
	}

	log.Infof("historicalChartsV2Data: %v", historicalChartsV2Data)
	if historicalChartsV2Data.Error != nil {
		log.Infof("historicalChartsV2Data Error: %s", *historicalChartsV2Data.Error)
	}

	log.Infof("historicalChartsV2DataStruct: %v", historicalChartsV2DataStruct)
	if historicalChartsV2DataStruct.Error != nil {
		log.Infof("historicalChartsV2DataStruct Error: %s", *historicalChartsV2DataStruct.Error)
	}

	log.Infof("tradeList: %v", tradeList)
	if tradeList.Error != nil {
		log.Infof("tradeList Error: %s", *tradeList.Error)
	}

	log.Infof("tradeDetail: %v", tradeDetail)
	if tradeDetail.Error != nil {
		log.Infof("tradeDetail Error: %s", *tradeDetail.Error)
	}

	log.Infof("brokerageCharges: %v", brokerageCharges)
	if brokerageCharges.Error != nil {
		log.Infof("brokerageCharges Error: %s", *brokerageCharges.Error)
	}

	log.Infof("breezeOrder: %v", breezeOrder)
	if breezeOrder.Error != nil {
		log.Infof("breezeOrder Error: %s", *breezeOrder.Error)
	}
}

func createSamplePorfolioHoldingsRequestBody() map[string]string {
	quotesRequestBody := make(map[string]string)
	quotesRequestBody["exchange_code"] = "NSE"

	return quotesRequestBody
}

func createSampleQuotesRequestBody() map[string]string {
	quotesRequestBody := make(map[string]string)
	quotesRequestBody["stock_code"] = "TATMOT"
	quotesRequestBody["exchange_code"] = "NSE"
	quotesRequestBody["expiry_date"] = "2023-08-26T06:00:00.000Z"
	quotesRequestBody["product_type"] = "cash"
	quotesRequestBody["right"] = "others"
	quotesRequestBody["strike_price"] = ""

	return quotesRequestBody
}
