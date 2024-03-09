package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	totp             = "311718"
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

type FivePaisaRequestPayload struct {
	Head FivePaisaHeadData `json:"head"`
	Body FivePaisaBodyData `json:"body"`
}

type FivePaisaHeadData struct {
	Key string `json:"key"`
}

type FivePaisaBodyData struct {
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

func fivePaisaLogin(totp string) {
	totpLoginResponse, err := totpLogin()
	if err != nil {
		fmt.Println("Error during TOTP login:", err)
		return
	}

	accessTokenResponse, err := getAccessToken(totpLoginResponse.Body.RequestToken)
	if err != nil {
		fmt.Println("Error getting access token:", err)
		return
	}

	fmt.Println(accessTokenResponse)
	getHoldings(accessTokenResponse.Body.AccessToken)
}

func totpLogin() (*TOTPLoginResponseBody, error) {
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

	totpLoginJSON, _ := json.Marshal(totpLoginData)
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

	getAccessTokenJSON, _ := json.Marshal(getAccessTokenData)
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
	payloadJson, err := createFivePaisaPayload()
	if err != nil {
		return nil, err
	}

	respData, err := makePOSTRequest(getAccessTokenURL, payloadJson, accessToken)
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

func createFivePaisaPayload() ([]byte, error) {
	payload := FivePaisaRequestPayload{
		Head: FivePaisaHeadData{
			Key: userKey,
		},
		Body: FivePaisaBodyData{
			ClientCode: clientCode,
		},
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error:", err)
		return payloadJson, err
	}

	return payloadJson, nil
}
