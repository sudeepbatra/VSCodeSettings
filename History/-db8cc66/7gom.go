package controller

import (
	"errors"

	"github.com/sudeepbatra/alpha-hft/broker/smartapi"
	"github.com/sudeepbatra/alpha-hft/logger"
)

func initializeSmartApiBroker(smartAPITotp string) error {
	smartapi.Login(smartAPITotp)

	if !smartapi.IsLoggedIn() {
		return errors.New("SmartAPI login failed")
	}

	return nil
}

func InitializeBrokers(smartAPITotp string) {
	//triggers login for smart api
	if smartAPITotp != "" {
		err := initializeSmartApiBroker(smartAPITotp)
		if err != nil {
			logger.Log.Error().Err(err).Str("broker", "SmartApi").Msg("Unable to succsessfully initialize broker")
		}
	}
	// login trigger for other brokers
}

func initializeSmartApiBrokerFromState(connectWebsocket bool) error {
	err := smartapi.LoginFromState()
	if err != nil {
		return err
	}

	if !smartapi.IsLoggedIn() {
		return errors.New("SmartAPI login failed")
	}

	if connectWebsocket {

		err = smartapi.SmartApiBrokers["default"].Connect()

		if err != nil {
			return err
		}

		smartapi.SubscribeInstrumentsOnStartup()
	}

	return nil
}

func InitializeBrokersFromCurrentState(connectWebsocket bool) {
	logger.Log.Info().Msg("initialize broker from saved state")
	err := initializeSmartApiBrokerFromState(connectWebsocket)
	if err != nil {
		logger.Log.Error().Err(err).Str("broker", "SmartApi").Msg("Unable to succsessfully initialize broker")
	}
}
