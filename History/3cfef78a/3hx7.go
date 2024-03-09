package fivepaisa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sudeepbatra/alpha-hft/logger"
)

var log = logger.GetLogger()

const (
	baseURL          = "https://Openapi.5paisa.com/VendorsAPI/Service1.svc/"
	appName          = "5P50603710"
	appSource        = "10427"
	userID           = "uxuZEFys5nv"
	password         = "7elTHyW0EC3"
	userKey          = "sR12m8nkT8VEPXtfgLFlspj5BQlSqB51"
	encryptionKey    = "jTS6yEtvhXThvDTYNHQNVXmklWFEaeQj"
	mobileNumber     = "9325487506"
	pin              = "574812"
	clientCode       = "50603710"
	totpLoginRoute   = "TOTPLogin"
	accessTokenRoute = "GetAccessToken"
	holdingsRoute    = "V3/Holding"
)

type TOTPLoginRequestHead struct {
	Key string `json:"Key"`
}

type TOTPLoginRequestBody struct {
	EmailID string `json:"Email_ID"`
	TOTP    string `json:"TOTP"`
	PIN     string `json:"PIN"`
}

type TOTPLoginResponseBody struct {
	Body struct {
		ClientCode   string `json:"ClientCode"`
		Message      string `json:"Message"`
		RedirectURL  string `json:"RedirectURL"`
		RequestToken string `json:"RequestToken"`
		Status       int    `json:"Status"`
		UserKey      string `json:"Userkey"`
	} `json:"body"`
	Head struct {
		Status            int    `json:"Status"`
		StatusDescription string `json:"StatusDescription"`
	} `json:"head"`
}

type GetAccessTokenRequestBody struct {
	RequestToken string `json:"RequestToken"`
	EncryKey     string `json:"EncryKey"`
	UserID       string `json:"UserId"`
}

type AccessTokenResponse struct {
	Body AccessTokenResponseBodyData `json:"body"`
	Head AccessTokenResponseHeadData `json:"head"`
}

type AccessTokenResponseBodyData struct {
	AccessToken           string `json:"AccessToken"`
	AllowBseCash          string `json:"AllowBseCash"`
	AllowBseDeriv         string `json:"AllowBseDeriv"`
	AllowBseMF            string `json:"AllowBseMF"`
	AllowMCXComm          string `json:"AllowMCXComm"`
	AllowMcxSx            string `json:"AllowMcxSx"`
	AllowNSECurrency      string `json:"AllowNSECurrency"`
	AllowNSEL             string `json:"AllowNSEL"`
	AllowNseCash          string `json:"AllowNseCash"`
	AllowNseComm          string `json:"AllowNseComm"`
	AllowNseDeriv         string `json:"AllowNseDeriv"`
	AllowNseMF            string `json:"AllowNseMF"`
	BulkOrderAllowed      int    `json:"BulkOrderAllowed"`
	ClientCode            string `json:"ClientCode"`
	ClientName            string `json:"ClientName"`
	ClientType            string `json:"ClientType"`
	CommodityEnabled      string `json:"CommodityEnabled"`
	CustomerType          string `json:"CustomerType"`
	DPInfoAvailable       string `json:"DPInfoAvailable"`
	DemoTrial             string `json:"DemoTrial"`
	DirectMFCharges       int    `json:"DirectMFCharges"`
	IsIDBound             int    `json:"IsIDBound"`
	IsIDBound2            int    `json:"IsIDBound2"`
	IsOnlyMF              string `json:"IsOnlyMF"`
	IsPLM                 int    `json:"IsPLM"`
	IsPLMDefined          int    `json:"IsPLMDefined"`
	Message               string `json:"Message"`
	OTPCredentialID       string `json:"OTPCredentialID"`
	PGCharges             int    `json:"PGCharges"`
	PLMsAllowed           int    `json:"PLMsAllowed"`
	POAStatus             string `json:"POAStatus"`
	PasswordChangeFlag    int    `json:"PasswordChangeFlag"`
	PasswordChangeMessage string `json:"PasswordChangeMessage"`
	ReferralBenefits      int    `json:"ReferralBenefits"`
	RefreshToken          string `json:"RefreshToken"`
	RunningAuthorization  int    `json:"RunningAuthorization"`
	Status                int    `json:"Status"`
	VersionChanged        int    `json:"VersionChanged"`
}

type AccessTokenResponseHeadData struct {
	Status            int    `json:"Status"`
	StatusDescription string `json:"StatusDescription"`
}

type GenericRequestPayload struct {
	Head GenericHeadData `json:"head"`
	Body GenericBodyData `json:"body"`
}

type GenericHeadData struct {
	Key string `json:"key"`
}

type GenericBodyData struct {
	ClientCode string `json:"ClientCode"`
}

type HoldingResponseData struct {
	Head interface{} `json:"head"`
	Body Holdings    `json:"body"`
}

type Holdings struct {
	Data []Holding `json:"Data"`
}

type Holding struct {
	Exchange     string  `json:"Exch"`
	ExchangeType string  `json:"ExchType"`
	NseCode      int     `json:"NseCode"`
	BseCode      int     `json:"BseCode"`
	Symbol       string  `json:"Symbol"`
	FullName     string  `json:"FullName"`
	Quantity     int     `json:"Quantity"`
	CurrentPrice float32 `json:"CurrentPrice"`
	PoolQty      int     `json:"PoolQty"`
	DPQty        int     `json:"DPQty"`
	POASigned    string  `json:"POASigned"`
}

func Login(totp string) {
	totpLoginResponse, err := totpLogin(totp)
	if err != nil || totpLoginResponse.Body.Status != "" {
		log.Error("Error during TOTP login:", err)
		return
	}

	accessTokenResponse, err := getAccessToken(totpLoginResponse.Body.RequestToken)
	if err != nil {
		log.Error("Error getting access token:", err)
		return
	}

	log.Info("accessTokenResponse: ", accessTokenResponse)

	holdingsResponseData, err := getHoldings(accessTokenResponse.Body.AccessToken)
	if err != nil {
		log.Error("Error while trying to get holdings", err)
		return
	}

	log.Info("holdingsResponseData: ", holdingsResponseData)

	log.Info("***************************************************")
	log.Info("accessToken", accessTokenResponse.Body.AccessToken)
	log.Info("***************************************************")
}

func totpLogin(totp string) (*TOTPLoginResponseBody, error) {
	totpLoginURL := baseURL + totpLoginRoute
	totpLoginData := struct {
		Head TOTPLoginRequestHead `json:"head"`
		Body TOTPLoginRequestBody `json:"body"`
	}{
		Head: TOTPLoginRequestHead{Key: userKey},
		Body: TOTPLoginRequestBody{
			EmailID: mobileNumber,
			TOTP:    totp,
			PIN:     pin,
		},
	}

	totpLoginJSON, err := json.Marshal(totpLoginData)
	if err != nil {
		log.Error("Error while trying to totp login to 5Paisa", err)
	}

	respData, err := makePOSTRequest(totpLoginURL, totpLoginJSON, "")
	if err != nil {
		return nil, err
	}

	var totpLoginResponseBody TOTPLoginResponseBody
	if err := json.Unmarshal(respData, &totpLoginResponseBody); err != nil {
		return nil, err
	}

	return &totpLoginResponseBody, nil
}

func getAccessToken(requestToken string) (*AccessTokenResponse, error) {
	getAccessTokenURL := baseURL + accessTokenRoute
	getAccessTokenData := struct {
		Head TOTPLoginRequestHead      `json:"head"`
		Body GetAccessTokenRequestBody `json:"body"`
	}{
		Head: TOTPLoginRequestHead{Key: userKey},
		Body: GetAccessTokenRequestBody{
			RequestToken: requestToken,
			EncryKey:     encryptionKey,
			UserID:       userID,
		},
	}

	getAccessTokenJSON, err := json.Marshal(getAccessTokenData)
	if err != nil {
		log.Error("Error while trying to get access token json", err)
	}

	respData, err := makePOSTRequest(getAccessTokenURL, getAccessTokenJSON, "")
	if err != nil {
		return nil, err
	}

	var accessTokenResponse AccessTokenResponse
	if err := json.Unmarshal(respData, &accessTokenResponse); err != nil {
		return nil, err
	}

	return &accessTokenResponse, nil
}

func getHoldings(accessToken string) (*HoldingResponseData, error) {
	getAccessTokenURL := baseURL + holdingsRoute

	payloadJSON, err := createGenericRequestPayload()
	if err != nil {
		return nil, err
	}

	respData, err := makePOSTRequest(getAccessTokenURL, payloadJSON, accessToken)
	if err != nil {
		return nil, err
	}

	var holdingResponseData HoldingResponseData
	if err := json.Unmarshal(respData, &holdingResponseData); err != nil {
		return nil, err
	}

	return &holdingResponseData, nil
}

func makePOSTRequest(url string, data []byte, accessToken string) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if accessToken != "" {
		req.Header.Set("Authorization", "bearer "+accessToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return respData, nil
}

func createGenericRequestPayload() ([]byte, error) {
	payload := GenericRequestPayload{
		Head: GenericHeadData{
			Key: userKey,
		},
		Body: GenericBodyData{
			ClientCode: clientCode,
		},
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		log.Error("Error while marshal generic request payload:", err)
		return payloadJSON, err
	}

	return payloadJSON, nil
}
