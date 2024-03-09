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

var CodeExchangeTypes = map[int]string{
	1:  "NSE",
	2:  "NSEFO",
	3:  "BSE",
	4:  "BSEFO",
	5:  "MCX",
	7:  "NCX",
	13: "CDE",
}

var NseCmTokens = []string{
	"25", "15083", "157", "236", "5900", "16669", "317", "16675", "526", "10604", "547", "694",
	"20374", "10940", "881", "910", "1232", "7229", "1333", "467", "1348", "1363", "1394", "4963", "1660",
	"5258", "1594", "11723", "1922", "17818", "11483", "2031", "10999", "11630", "17963", "2475", "14977",
	"2885", "21808", "3045", "3351", "11536", "3432", "3456", "3499", "13538", "3506", "11287", "11532", "3787",
	"99926000", "99926012", "99926009", "99926023", "99926025", "99926008", "99926004", "99926036", "99926030",
	"99926020", "99926018", "99926013", "99926029", "99926019",
}

var BseCmTokens = []string{
	"99919000",
}

var McxFoTokens = []string{
	"414", "152", "197", "114", "115",
}

var NseFoTokens = []string{
	"72898", "43500", "72889", "56589", "72896", "65623", "40887", "40874",
	"43273", "43272", "56587", "72897", "43269", "56601", "43511", "72893", "43279", "42205", "65621",
	"43514", "41708", "72894", "72888", "42195", "43257", "40890", "40883", "40893", "40870", "65625",
	"72895", "40892", "42211", "56592", "65640", "41709", "56588", "65646", "42208", "42206", "42200",
	"43512", "72884", "42201", "72887", "43521", "56591", "56596", "65624", "40881", "42207", "42210",
	"56602", "43270", "48266", "42197", "56593", "43268", "43507", "42196", "72890", "72899", "40888",
	"43255", "43506", "56595", "43259", "48270", "65622", "40869", "43256", "56590", "65626", "65627",
	"65639", "42194", "65641", "56594", "40882", "43271", "65638", "43499", "40891", "43258",
}
