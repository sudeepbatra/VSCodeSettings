package smartapi

import (
	"fmt"

	SmartApi "github.com/angel-one/smartapigo"
)

func smartApiSessionToken() {
	ABClient := SmartApi.New("S1632585", "DevDas123", "zlndzd6z")

	fmt.Println("Client :- ", ABClient)

	// User Login and Generate User Session
	session, err := ABClient.GenerateSession("totp here")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	session.UserSessionTokens, err = ABClient.RenewAccessToken(session.RefreshToken)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("User Session Tokens :- ", session.UserSessionTokens)

	//Get User Profile
	session.UserProfile, err = ABClient.GetUserProfile()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("User Profile :- ", session.UserProfile)
	fmt.Println("User Session Object :- ", session)

	//Place Order
	order, err := ABClient.PlaceOrder(SmartApi.OrderParams{Variety: "NORMAL", TradingSymbol: "SBIN-EQ", SymbolToken: "3045", TransactionType: "BUY", Exchange: "NSE", OrderType: "LIMIT", ProductType: "INTRADAY", Duration: "DAY", Price: "19500", SquareOff: "0", StopLoss: "0", Quantity: "1"})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Placed Order ID and Script :- ", order)
}
