package cnbc

type FairValueQuoteResult struct {
	FairValueQuote []fairValueQuote `json:"FairValueQuote"`
}

type fairValueQuote struct {
	Symbol         string `json:"symbol"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	Last           string `json:"last"`
	LastTime       string `json:"last_time"`
	LastTimeMsec   string `json:"last_time_msec"`
	Exchange       string `json:"exchange"`
	Provider       string `json:"provider"`
	TodaysClosing  string `json:"todays_closing"`
	ProviderSymbol string `json:"providerSymbol"`
	IndexClose     string `json:"index_close"`
	FVClose        string `json:"fv_close"`
	FVChange       string `json:"fv_change"`
	FVSpread       string `json:"fv_spread"`
	FVRaw          string `json:"fv_raw"`
	LastTimedate   string `json:"last_timedate"`
	Realtime       string `json:"realtime"`
	Shortname      string `json:"shortname"`
	Altsymbol      string `json:"altsymbol"`
	IssueID        string `json:"issueid"`
	FmtLast        string `json:"fmt_last"`
	FmtChange      string `json:"fmt_change"`
	FVChangePct    string `json:"fv_change_pct"`
	ChangePct      string `json:"change_pct"`
	Change         string `json:"change"`
}
