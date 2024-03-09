package cnbc

import "time"

type FairValueQuote struct {
	Symbol         string    `json:"symbol"`
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	Last           string    `json:"last"`
	LastTime       time.Time `json:"-"`
	LastTimeString string    `json:"last_time"`
	LastTimeMsec   string    `json:"last_time_msec"`
	Exchange       string    `json:"exchange"`
	Provider       string    `json:"provider"`
	TodaysClosing  string    `json:"todays_closing"`
	ProviderSymbol string    `json:"providerSymbol"`
	IndexClose     string    `json:"index_close"`
	FVClose        string    `json:"fv_close"`
	FVChange       string    `json:"fv_change"`
	FVSpread       string    `json:"fv_spread"`
	FVRaw          string    `json:"fv_raw"`
	LastTimeDate   string    `json:"last_timedate"`
	Realtime       string    `json:"realtime"`
	Shortname      string    `json:"shortname"`
	AltSymbol      string    `json:"altsymbol"`
	IssueID        string    `json:"issueid"`
	FmtLast        string    `json:"fmt_last"`
	FmtChange      string    `json:"fmt_change"`
	FVChangePct    string    `json:"fv_change_pct"`
	ChangePct      string    `json:"change_pct"`
	Change         string    `json:"change"`
}

type FairValueQuoteResult struct {
	FairValueQuote []FairValueQuote `json:"FairValueQuote"`
}

type QuoteData struct {
	FairValueQuoteResult FairValueQuoteResult `json:"FairValueQuoteResult"`
}

type OtherSymbolsQuote struct {
	Symbol             string `json:"symbol"`
	SymbolType         string `json:"symbolType"`
	Code               int    `json:"code"`
	Name               string `json:"name"`
	ShortName          string `json:"shortName"`
	OnAirName          string `json:"onAirName"`
	AltName            string `json:"altName"`
	Last               string `json:"last"`
	LastTimeDate       string `json:"last_timedate"`
	LastTime           string `json:"last_time"`
	ChangeType         string `json:"changetype"`
	Type               string `json:"type"`
	SubType            string `json:"subType"`
	Exchange           string `json:"exchange"`
	Source             string `json:"source"`
	Open               string `json:"open"`
	High               string `json:"high"`
	Low                string `json:"low"`
	Change             string `json:"change"`
	ChangePercent      string `json:"change_pct"`
	CurrencyCode       string `json:"currencyCode"`
	Volume             string `json:"volume"`
	VolumeAlt          string `json:"volume_alt"`
	Provider           string `json:"provider"`
	PreviousDayClosing string `json:"previous_day_closing"`
	AltSymbol          string `json:"altSymbol"`
	RealTime           string `json:"realTime"`
	CurMktStatus       string `json:"curmktstatus"`
	YrHiPrice          string `json:"yrhiprice"`
	YrHiDate           string `json:"yrhidate"`
	YrLoPrice          string `json:"yrloprice"`
	YrLoDate           string `json:"yrlodate"`
	Streamable         string `json:"streamable"`
	IssueID            string `json:"issue_id"`
	CountryCode        string `json:"countryCode"`
	TimeZone           string `json:"timeZone"`
	FeedSymbol         string `json:"feedSymbol"`
	PortfolioIndicator string `json:"portfolioindicator"`
}

type OtherSymbolsQuotes struct {
	FormattedQuote []OtherSymbolsQuote `json:"FormattedQuote"`
}

type OtherSymbolsData struct {
	FairValueQuoteResult OtherSymbolsQuotes `json:"FormattedQuoteResult"`
}
