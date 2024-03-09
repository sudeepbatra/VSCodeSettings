package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const stockScriptCSVURL = "https://traderweb.icicidirect.com/Content/File/txtFile/ScripFile/StockScriptNew.csv"
const API_URL = "https://api.icicidirect.com/breezeapi/api/v1/"
const customerDetailsEndpoint = "customerdetails"

type ApificationBreeze struct {
	*Breeze
	hostname           string
	base64SessionToken string
}

type Breeze struct {
	appName             string
	apiKey              string
	secretKey           string
	sessionToken        string
	userID              string
	stockScriptDictList map[string]map[string]string
	tokenScriptDictList map[string]map[string][]string
}

func generateSession(appName, apiKey, secretKey, sessionToken string) (*ApificationBreeze, error) {
	breezeInstance := &Breeze{
		appName:      appName,
		apiKey:       apiKey,
		sessionToken: sessionToken,
		secretKey:    secretKey,
	}
	breezeInstance.retrieveSessionInfo()
	breezeInstance.getStockScriptList()
	apiHandler := &ApificationBreeze{
		Breeze:   breezeInstance,
		hostname: API_URL,
	}

	apiHandler.base64SessionToken = base64.StdEncoding.EncodeToString([]byte(breezeInstance.userID + ":" + breezeInstance.sessionToken))
	return apiHandler, nil
}

func (b *Breeze) retrieveSessionInfo() error {
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	log.Println(headers)
	body := map[string]string{
		"SessionToken": b.sessionToken,
		"AppKey":       b.apiKey,
	}

	requestBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	url := API_URL + customerDetailsEndpoint

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var responseData map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return err
	}
	if responseData["Success"] != nil {
		base64SessionToken := responseData["Success"].(map[string]interface{})["session_token"].(string)
		decoded, err := base64.StdEncoding.DecodeString(base64SessionToken)
		if err != nil {
			return err
		}
		sessionInfo := string(decoded)
		sessionParts := strings.Split(sessionInfo, ":")
		if len(sessionParts) >= 2 {
			b.userID = sessionParts[0]
			b.sessionToken = sessionParts[1]
		} else {
			return fmt.Errorf("apiUtil sessionParts less than 2")
		}
	} else {
		return fmt.Errorf("apiUtil responseData does not have success")
	}
	return nil
}

func (b *Breeze) getStockScriptList() error {
	// Fetch CSV content
	response, err := http.Get(stockScriptCSVURL)

	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Decode CSV content and process rows
	decodedContent, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	cr := csv.NewReader(strings.NewReader(string(decodedContent)))
	myList, err := cr.ReadAll()
	if err != nil {
		return err
	}

	// Initialize the stock script and token script maps
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

	// Update the stock script and token script maps in the Breeze struct
	b.stockScriptDictList = stockScriptDictList
	b.tokenScriptDictList = tokenScriptDictList

	return nil
}

func (a *ApificationBreeze) errorException(funcName string, err error) {
	message := fmt.Sprintf("%s() Error: %s", funcName, err)
	panic(message)
}

func (a *ApificationBreeze) validationErrorResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"Success": "",
		"Status":  500,
		"Error":   message,
	}
}

func (a *ApificationBreeze) generateHeaders(body string) map[string]string {
	currentDate := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	checksum := sha256HexDigest(currentDate + body + a.Breeze.secretKey)
	headers := map[string]string{
		"Content-Type":   "application/json",
		"X-Checksum":     "token " + checksum,
		"X-Timestamp":    currentDate,
		"X-AppKey":       a.Breeze.apiKey,
		"X-SessionToken": a.base64SessionToken,
	}
	return headers
}

func sha256HexDigest(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func (a *ApificationBreeze) makeRequest(method, endpoint, body string, headers map[string]string) ([]byte, error) {
	url := a.hostname + endpoint
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		a.errorException(fmt.Sprintf("makeRequest(%s %s)", method, url), err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	res, err := client.Do(req)
	if err != nil {
		a.errorException(fmt.Sprintf("makeRequest(%s %s)", method, url), err)
	}
	defer res.Body.Close()
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		a.errorException(fmt.Sprintf("makeRequest(%s %s)", method, url), err)
	}
	return responseBody, nil
}

func (a *ApificationBreeze) getCustomerDetails(apiSession string) (map[string]interface{}, error) {
	if apiSession == "" {
		return a.validationErrorResponse("Empty apiSession value received in getCustomerDetails"), nil
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	body := map[string]interface{}{
		"SessionToken": apiSession,
		"AppKey":       a.Breeze.apiKey,
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		a.errorException("getCustomerDetails", err)
	}
	responseBody, err := a.makeRequest(http.MethodGet, customerDetailsEndpoint, string(bodyJSON), headers)
	if err != nil {
		a.errorException("getCustomerDetails", err)
	}
	var response map[string]interface{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		a.errorException("getCustomerDetails", err)
	}
	if success, ok := response["Success"].(map[string]interface{}); ok && success != nil {
		delete(success, "session_token")
		response["Success"] = success
	}
	return response, nil
}

func main() {
	appName := "StoxMarket"
	apiKey := "s76162#+U35414Y*S413=099_FA6P567"
	secretKey := "I04M0Y9!5vP7G3ct53e395+41F27=621"

	loginURL := fmt.Sprintf("https://api.icicidirect.com/apiuser/login?api_key=%s", url.QueryEscape(apiKey))
	log.Println("=======================================")
	log.Println("Use the url to generate the sessionToken and set it below")
	log.Println(loginURL)
	log.Println("=======================================")

	sessionToken := "16402994"

	apiHandler, err := generateSession(appName, apiKey, secretKey, sessionToken)
	if err != nil {
		fmt.Println("Error creating API handler:", err)
		return
	}

	customerDetails, err := apiHandler.getCustomerDetails(sessionToken)
	if err != nil {
		fmt.Println("Error getting customer details:", err)
		return
	}

	log.Println("-----000----- Stock Script Dict List -----000----- ")
	log.Println(apiHandler.stockScriptDictList)
	log.Println("-----000----- Stock Script Dict List -----000----- ")
	log.Println("-----^^^----- Token Script Dict List -----^^^----- ")
	log.Println(apiHandler.tokenScriptDictList)
	log.Println("-----^^^----- Token Script Dict List -----^^^----- ")

	fmt.Println("Customer details:", customerDetails)
}
