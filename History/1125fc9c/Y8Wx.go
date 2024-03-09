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
	baseURL        = "https://Openapi.5paisa.com/VendorsAPI/Service1.svc/"
	appName        = "5P50603710"
	appSource      = "10427"
	userID         = "uxuZEFys5nv"
	password       = "7elTHyW0EC3"
	userKey        = "sR12m8nkT8VEPXtfgLFlspj5BQlSqB51"
	encryptionKey  = "jTS6yEtvhXThvDTYNHQNVXmklWFEaeQj"
	mobileNumber   = "9325487506"
	totp           = "434332"
	pin            = "574812"
	totpLoginRoute = "TOTPLogin"
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

type TimeMilli time.Time

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

func main() {
	totpLoginResponse, err := totpLogin()
	if err != nil {
		fmt.Println("Error during TOTP login:", err)
		return
	}

	getAccessTokenResponse, err := getAccessToken(totpLoginResponse.Body.RequestToken)
	if err != nil {
		fmt.Println("Error getting access token:", err)
		return
	}

	fmt.Println(getAccessTokenResponse)
}

func totpLogin() (*TOTPLoginResponseBody, error) {
	totpLoginURL := fmt.Sprintf("%s%s", baseURL, totpLoginRoute)
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
	respData, err := makePOSTRequest(totpLoginURL, totpLoginJSON)
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
	getAccessTokenURL := fmt.Sprintf("%s/GetAccessToken", baseURL)
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
	respData, err := makePOSTRequest(getAccessTokenURL, getAccessTokenJSON)
	if err != nil {
		return nil, err
	}

	var accessTokenResponse AccessTokenResponse
	if err := json.Unmarshal(respData, &accessTokenResponse); err != nil {
		return nil, err
	}

	return &accessTokenResponse, nil
}

func makePOSTRequest(url string, data []byte) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(data))
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
