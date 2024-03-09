package cnbc

type fairValueQuoteResult struct {
	FairValueQuote []fairValueQuote `json:"FairValueQuote"`
}

type fairValueQuote struct {
	Symbol string `json:"symbol"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Last   string `json:"last"`
	// ... (other fields)
}
