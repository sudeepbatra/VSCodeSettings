package smartapiwslogin

import (
	"fmt"

	SmartApi "github.com/angel-one/smartapigo"
	"github.com/sudeepbatra/alpha-hft/logger"
)

var log = logger.GetLogger()

func SmartAPIWebSocketLogin(totp string) {
	ABClient := SmartApi.New("S1632585", "4321", "iGKWS2zU")
	log.Info("Client: ", ABClient)

	session, err := ABClient.GenerateSession(totp)
	if err != nil {
		log.Error(err.Error())
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
	log.Info("Session User Profile: ", session.UserProfile)

	log.Info("Session: ", session)
	log.Info("Feed Token: ", session.FeedToken)
	log.Info("Access Token: ", session.AccessToken)
	log.Info("Refresh Token: ", session.RefreshToken)
}
