package smartapi_ws_login

import "fmt"

func SmartAPIWebSocketLogin(totp string) {

	// Create New Angel Broking Client
	ABClient := SmartApi.New("S1632585", "4321", "iGKWS2zU")

	fmt.Println("Client :- ", ABClient)

	// User Login and Generate User Session
	session, err := ABClient.GenerateSession("595973")

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
	fmt.Println("Feed Token: ", session.FeedToken)
	fmt.Println("Access Token: ", session.AccessToken)
	fmt.Println("Refresh Token: ", session.RefreshToken)

}
