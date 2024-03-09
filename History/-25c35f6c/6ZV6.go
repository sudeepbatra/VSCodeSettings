package common

import (
	"crypto/tls"
	"net/http"
	"net/url"
)

func GetClient(useProxy bool) *http.Client {
	if useProxy {
		proxyURL := "http://localhost:8080"

		proxyURLParsed, err := url.Parse(proxyURL)
		if err != nil {
			panic("Failed to parse proxy URL: " + err.Error())
		}

		proxyTransport := &http.Transport{
			Proxy:               http.ProxyURL(proxyURLParsed),
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		}

		return &http.Client{
			Transport: proxyTransport,
		}
	}

	return &http.Client{}
}
