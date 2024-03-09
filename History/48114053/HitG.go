package cnbc

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/sudeepbatra/alpha-hft/common"
	"github.com/sudeepbatra/alpha-hft/data"
	"github.com/sudeepbatra/alpha-hft/logger"
)

const (
	cnbcQuotesURL             = "https://quote.cnbc.com/quote-html-webservice/fvquote.htm?requestMethod=quick&noform=1&realtime=0&client=fairValue&output=json&symbols=DJ%7CSP%7CND%7CTF"
	cnbcOtherSymbolsQuotesURL = "https://quote.cnbc.com/quote-html-webservice/restQuote/symbolType/symbol?partnerId=2&requestMethod=quick&exthrs=1&noform=1&fund=1&output=json&symbols=.HSI%7C.N225%7C.STI%7C.AXJO%7C.SSEC%7C.FTSE%7C.GDAXI%7C.FCHI%7C.AEX%7C.STOXX50%7C@CL.1%7C@RB.1%7C@NG.1%7C@GC.1%7C@SI.1%7CEUR=%7CJPY=%7CGBP=%7CCHF=%7CCAD=%7C.VIX%7C.VXN%7C.OVX%7C.OOI%7C.FTFCNBCG%7CUS3M%7CUS2Y%7CUS5Y%7CUS10Y%7CUS30Y%7C.RUT%7C.DJT%7C.DJU%7C.NDX%7C.NYA%7C.GSPT%7C.GSPS%7C.GSPD%7C.GSPE%7C.GSPF"
)

func FetchQuoteDataAndSaveInDB() {
	logger.Log.Info().
		Msg("Fetching CNBC Quotes and Saving in DB...")

	client := common.GetClient(false)

	req, err := http.NewRequest("GET", cnbcQuotesURL, nil)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error creating request")
		return
	}

	setRequestHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error making request")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error reading response body")
		return
	}

	logger.Log.Trace().
		Interface("response", string(body)).
		Msg("response from CNBC quote api")

	var quoteData QuoteData
	err = json.Unmarshal(body, &quoteData)

	if err != nil {
		logger.Log.Error().
			Err(err).
			Msg("error in unmarshalling response from cnbc quotes api")
	}

	convertLastTimeToLocalTime(quoteData)

	logger.Log.Trace().
		Interface("quoteData", quoteData).
		Msg("quoteData")

	for _, fairValueQuote := range quoteData.FairValueQuoteResult.FairValueQuote {
		err := data.InsertRecord("fair_value_quotes", fairValueQuote)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error storing FairValueQuote to the DB:")
		}
	}

	cnbcOtherSymbolsQuotesRequest, err := http.NewRequest("GET", cnbcOtherSymbolsQuotesURL, nil)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error creating cnbc other symbols quotes url")
		return
	}

	cnbcOtherSymbolsQuotesResponse, err := client.Do(cnbcOtherSymbolsQuotesRequest)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error making cnbc other symbols quotes request")
		return
	}
	defer resp.Body.Close()

	cnbcOtherSymbolsQuotesBody, err := io.ReadAll(cnbcOtherSymbolsQuotesResponse.Body)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error reading cnbc other symbols quotes body")
		return
	}

	logger.Log.Trace().
		Interface("response", string(cnbcOtherSymbolsQuotesBody)).
		Msg("response from CNBC quote api")

	var formattedQuoteData FormattedQuoteData
	err = json.Unmarshal(body, &formattedQuoteData)

	if err != nil {
		logger.Log.Error().
			Err(err).
			Msg("error in unmarshalling response from cnbc quotes api")
	}

	logger.Log.Info().
		Interface("formattedQuoteData", formattedQuoteData).
		Msg("formattedQuoteData")

	logger.Log.Info().
		Msg("Finished Fetching CNBC Quotes and Saving in DB!")
}

func convertLastTimeToLocalTime(quoteData QuoteData) {
	for i, quote := range quoteData.FairValueQuoteResult.FairValueQuote {
		parsedTime, err := time.Parse("2006-01-02T15:04:05.000-0700", quote.LastTimeString)
		if err != nil {
			logger.Log.Error().Err(err).Msgf("error parsing time for quote %d", i)
			continue
		}

		quoteData.FairValueQuoteResult.FairValueQuote[i].LastTime = parsedTime
	}

	localTimeZone := time.Local
	for i := range quoteData.FairValueQuoteResult.FairValueQuote {
		quoteData.FairValueQuoteResult.FairValueQuote[i].LastTime = quoteData.FairValueQuoteResult.FairValueQuote[i].LastTime.In(localTimeZone)
	}
}

func setRequestHeaders(req *http.Request) {
	req.Header.Set("authority", "quote.cnbc.com")
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("origin", "https://www.cnbc.com")
	req.Header.Set("referer", "https://www.cnbc.com/")
	req.Header.Set("sec-ch-ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
}
