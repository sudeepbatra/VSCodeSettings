package breeze

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/sudeepbatra/alpha-hft/logger"
)

var log = logger.GetLogger()

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

// retrieveSessionToken retrieves Session Token using the customer details endpoint.
func (b *Breeze) retrieveSessionToken() error {
	log.Debug("Starting retrieving Session Info...")

	body := map[string]string{
		"SessionToken": b.sessionToken,
		"AppKey":       b.apiKey,
	}

	requestBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("json.Marshal failed in retrieveSessionToken: %v", err)
	}

	url := apiURL + customerDetailsEndpoint

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("http get new request failed in retrieveSessionToken: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	log.Debug("Executing the request: ", req)
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do error in retrieveSessionToken: %v", err)
	}
	defer response.Body.Close()

	var responseData map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return fmt.Errorf("json.NewDecoder failed for retrieveSessionToken: %v", err)
	}
	log.Debug("Received responseData:", responseData)

	if responseData["Success"] != nil {
		base64SessionToken := responseData["Success"].(map[string]interface{})["session_token"].(string)
		log.Debug("Received Base64 Session Token", base64SessionToken)
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
			return fmt.Errorf("retrieveSessionToken sessionParts less than 2")
		}
	} else {
		return fmt.Errorf("apiUtil responseData does not have success")
	}

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

func (a *ApificationBreeze) makeRequest(method, endpoint, body string, headers http.Header, useProxy bool) ([]byte, error) {
	client := getClient(useProxy)

	url := a.hostname + endpoint
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		a.errorException(fmt.Sprintf("makeRequest(%s %s)", method, url), err)
	}

	req.Header = headers

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

func getClient(useProxy bool) *http.Client {
	if useProxy {
		proxyURL := "http://localhost:8080" // Replace with your proxy URL
		proxyURLParsed, err := url.Parse(proxyURL)
		if err != nil {
			panic("Failed to parse proxy URL: " + err.Error())
		}

		proxyTransport := &http.Transport{
			Proxy:           http.ProxyURL(proxyURLParsed),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Use this if you trust the proxy, but it's not recommended for production.
		}

		return &http.Client{
			Transport: proxyTransport,
		}
	}
	
	return &http.Client{}
}
