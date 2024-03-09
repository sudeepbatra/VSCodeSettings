package smartapiwslogin

import (
	"fmt"

	SmartApi "github.com/angel-one/smartapigo"
	"github.com/labstack/gommon/log"
)

func SmartAPIWebSocketLogin(totp string) {
	ABClient := SmartApi.New("S1632585", "4321", "iGKWS2zU")
	log.Info("Client: ", ABClient)

	session, err := ABClient.GenerateSession(totp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Info("Session: ", session)

	session.UserSessionTokens, err = ABClient.RenewAccessToken(session.RefreshToken)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Info("User Session Tokens: ", session.UserSessionTokens)

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
