package smartapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	SmartApi "github.com/angel-one/smartapigo"
	"github.com/gorilla/websocket"
	"github.com/sudeepbatra/alpha-hft/common"
	"github.com/sudeepbatra/alpha-hft/config"
	"github.com/sudeepbatra/alpha-hft/logger"
)

const (
	maxOptionsExpiryDateFromToday = 40
	dateFormatLayout              = "02Jan2006"
)

type Client struct {
	WsConn      *websocket.Conn
	SmartClient *SmartApi.Client
}

type SmartApiApplication struct {
	*Client
	AccessToken *string
	FeedToken   string
	ClientCode  string
	AppKey      string
	Password    string
}

var SmartApiBrokers = make(map[string]*SmartApiApplication)

func (s *SmartApiApplication) initializeApplication(totp string, appName string) error {
	logger.Log.Info().
		Str("totp", totp).
		Str("app name", appName).
		Msg("Creating smartApi client for")

	smartAPIClient := SmartApi.New(s.ClientCode, s.Password, s.AppKey)

	session, err := smartAPIClient.GenerateSession(totp)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("totp", totp).
			Str("app name", appName).
			Msg("Error in generating smartapi session for")

		return err
	}

	logger.Log.Info().
		Str("app name", appName).
		Msg("Session generated for broker. Generating renewal access token.")

	session.UserSessionTokens, err = smartAPIClient.RenewAccessToken(session.RefreshToken)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("app name", appName).
			Msg("Error in generating smartapi renew access token for")

		return err
	}

	logger.Log.Debug().Any("tokens", session.UserSessionTokens).Msg("session token information")
	s.AccessToken = &session.AccessToken
	s.FeedToken = session.FeedToken
	s.SmartClient = smartAPIClient

	logger.Log.Debug().
		Str("smartapi", appName).
		Msg("saving the session to the state")

	err = updateSmartApiState(s, appName)

	if err != nil {
		logger.Log.Error().Err(err).
			Str("app name", appName).
			Msg("error in updating smart api state")

		return err
	}

	return nil
}

func (s *SmartApiApplication) generateHeaders() http.Header {
	headers := make(http.Header)
	headers["Authorization"] = []string{fmt.Sprintf("Bearer %s", *s.AccessToken)}
	headers["x-api-key"] = []string{s.AppKey}
	headers["x-client-code"] = []string{s.ClientCode}
	headers["x-feed-token"] = []string{s.FeedToken}

	return headers
}

func Login(totp string) {
	logger.Log.Info().Str("broker", "smartapi").Msg("initializing smartapi application")

	SmartApiBrokers["default"] = &SmartApiApplication{
		Client:     &Client{},
		ClientCode: config.Config.BrokerConfig.SmartApi.ClientCode,
		AppKey:     config.Config.BrokerConfig.SmartApi.Key,
		Password:   config.Config.BrokerConfig.SmartApi.Password,
	}
	logger.Log.Info().Str("broker", "smartapi").Msg("initializing smartapi historical application")

	SmartApiBrokers["historical"] = &SmartApiApplication{
		Client:     &Client{},
		ClientCode: config.Config.BrokerConfig.SmartAPIHistorical.ClientCode,
		AppKey:     config.Config.BrokerConfig.SmartAPIHistorical.APIKey,
		Password:   config.Config.BrokerConfig.SmartAPIHistorical.Password,
	}

	SmartApiBrokers["trading"] = &SmartApiApplication{
		Client:     &Client{},
		ClientCode: config.Config.BrokerConfig.SmartAPITrading.ClientCode,
		AppKey:     config.Config.BrokerConfig.SmartAPITrading.APIKey,
		Password:   config.Config.BrokerConfig.SmartAPITrading.Password,
	}

	logger.Log.Info().
		Str("smartapi", "default").
		Msg("initializing smartapi default: SmartAPI Market API Type")

	if err := SmartApiBrokers["default"].initializeApplication(totp, "default"); err != nil {
		logger.Log.Error().Err(err).Msg("error in initializing smartapi default application")
		return
	}

	logger.Log.Info().
		Str("smartapi", "historical").
		Msg("initializing smartapi default: SmartAPI Historical API Type")

	if err := SmartApiBrokers["historical"].initializeApplication(totp, "historical"); err != nil {
		logger.Log.Error().Err(err).Msg("error in initializing smartapi historical application")
		return
	}

	logger.Log.Info().
		Str("smartapi", "trading").
		Msg("initializing smartapi default: SmartAPI Trading API Type")

	if err := SmartApiBrokers["trading"].initializeApplication(totp, "trading"); err != nil {
		logger.Log.Error().Err(err).Msg("error in initializing smartapi trading application")
		return
	}

	logger.Log.Debug().
		Str("broker", "SmartApi").
		Str("default access token", *SmartApiBrokers["default"].AccessToken).
		Str("historical access token", *SmartApiBrokers["default"].AccessToken).
		Msg("Login successful: access token generated for market, historical and trading SmartAPI api's")
}

func LoginFromState() error {
	file, err := os.ReadFile("state.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(file, &SmartApiBrokers); err != nil {
		return err
	}

	return nil
}

func IsLoggedIn() bool {
	return SmartApiBrokers["default"] != nil &&
		SmartApiBrokers["default"].AccessToken != nil &&
		SmartApiBrokers["historical"] != nil &&
		SmartApiBrokers["historical"].AccessToken != nil &&
		SmartApiBrokers["trading"] != nil &&
		SmartApiBrokers["trading"].AccessToken != nil
}

func SubscribeInstrumentsOnStartup() {
	nfoTokens, _ := getNFOTokens()

	request := Subscription{
		CorrelationID: "abcde12345",
		Action:        SubscribeAction, // Subscribe action
		Params: RequestParams{
			Mode: SnapQuote, // Subscription mode (LTP)
			TokenLists: []RequestTokenList{
				{
					ExchangeType: NseCm,
					Tokens:       NseCmTokens,
				},
				{
					ExchangeType: BseCm,
					Tokens:       BseCmTokens,
				},
				{
					ExchangeType: McxFo,
					Tokens:       McxFoTokens,
				},
				{
					ExchangeType: Nfo,
					Tokens:       nfoTokens,
				},
			},
		},
	}

	if err := SmartApiBrokers["default"].Subscribe(request); err != nil {
		logger.Log.Error().Err(err).Msg("Error in subscribing to smartapi websocket")
	}
}

func getNFOTokens() ([]string, error) {
	nfoTokens := []string{}
	now := time.Now()

	client := common.GetClient(false)

	nifty50LTPResponse, err := SmartApiBrokers["default"].GetLtpData(client, "NSE", "Nifty 50", "99926000")
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error in getting ltp data for nifty 50")
		return nil, err
	}

	niftybankLTPResponse, err := SmartApiBrokers["default"].GetLtpData(client, "NSE", "Nifty Bank", "99926009")
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error in getting ltp data for nifty bank")
		return nil, err
	}

	niftyIndiceLastTradePrice := nifty50LTPResponse.Data.Ltp
	logger.Log.Debug().
		Float64("niftyIndiceLastTradePrice", niftyIndiceLastTradePrice).
		Msg("Nifty Indice Last Value: ")

	maxCEStrike := niftyIndiceLastTradePrice + 500.0
	minCEStrike := niftyIndiceLastTradePrice - 100.0
	maxPEStrike := niftyIndiceLastTradePrice + 100.0
	minPEStrike := niftyIndiceLastTradePrice - 500.0

	instruments, err := GetInstrumentsScripMasterData(client)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error fetching instruments data")
		return nil, err
	}

	for _, instrument := range instruments {
		if instrument.ExchSeg == "NFO" &&
			instrument.InstrumentType == "OPTIDX" &&
			instrument.Name == "NIFTY" {
			logger.Log.Debug().
				Interface("instrument", instrument).
				Msg("NSE FO Token")

			instrumentStrike, err := strconv.ParseFloat(instrument.Strike, 64)
			if err != nil {
				logger.Log.
					Err(err).
					Interface("instrument", instrument).
					Msg("Error parsing instrument strike price")

				continue
			}

			instrumentStrike = instrumentStrike / 100.0

			parsedExpiryDate, err := time.Parse(dateFormatLayout, instrument.Expiry)
			if err != nil {
				logger.Log.Err(err).
					Interface("instrument", instrument).
					Msg("Error parsing instrument expiry date")

				continue
			}

			logger.Log.Trace().
				Time("expiry date", parsedExpiryDate).
				Msg("Index Option Expiry Date: ")

			differenceInDays := int(parsedExpiryDate.Sub(now).Hours() / 24)

			if instrument.OptionType == "CE" &&
				instrumentStrike <= maxCEStrike &&
				instrumentStrike >= minCEStrike &&
				differenceInDays < maxOptionsExpiryDateFromToday {
				nfoTokens = append(nfoTokens, instrument.Token)
			}

			if instrument.OptionType == "PE" &&
				instrumentStrike <= maxPEStrike &&
				instrumentStrike >= minPEStrike &&
				differenceInDays < maxOptionsExpiryDateFromToday {
				nfoTokens = append(nfoTokens, instrument.Token)
			}
		}
	}

	return nfoTokens, nil
}
