package smartapi

var ExchangeCodeTypes = map[string]int{
	"NSE":   1,
	"NSEFO": 2,
	"BSE":   3,
	"BSEFO": 4,
	"MCX":   5,
	"NCX":   7,
	"CDE":   13,
}

var ReverseExchangeTypes = map[int]string{
	1:  "NSE",
	2:  "NSEFO",
	3:  "BSE",
	4:  "BSEFO",
	5:  "MCX",
	7:  "NCX",
	13: "CDE",
}
