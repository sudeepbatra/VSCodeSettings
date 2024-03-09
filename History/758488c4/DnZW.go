package smartapi

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sudeepbatra/alpha-hft/logger"
)

const (
	// Angel One Smart API Websocket Streaming 2.0
	rootURI               = "ws://smartapisocket.angelone.in/smart-stream"
	heartBeatMessage      = "ping"
	heartBeatInterval     = 10 * time.Second
	littleEndianByteOrder = "<"

	// Available Actions
	SubscribeAction   = 1
	UnsubscribeAction = 0

	// Possible Subscription Modes
	LTPMode   = 1
	Quote     = 2
	SnapQuote = 3

	// Exchange Types
	NseCm = 1
	Nfo   = 2
	BseCm = 3
	BseFo = 4
	McxFo = 5
	NcxFo = 7
	CdeFo = 13
)

func (s *SmartApiApplication) Connect() error {
	logger.Log.Info().Msg("Connecting to websocket for smartApi")
	logger.WebSocketLog.Info().Msg("Connecting to websocket for smartApi")

	headers := s.generateHeaders()

	conn, _, err := websocket.DefaultDialer.Dial(rootURI, headers)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error in connecting to smart api websocket")
		logger.WebSocketLog.Error().Err(err).Msg("Error in connecting to smart api websocket")

		return err
	}

	s.WsConn = conn
	// Heartbeat message
	logger.Log.Info().Msg("Creating corouting for sending heartbeats to the smart api websocket")
	logger.WebSocketLog.Info().Msg("Creating corouting for sending heartbeats to the smart api websocket")

	go s.sendHeartbeats()
	go s.consumeMessage()

	return nil
}

func (s *SmartApiApplication) consumeMessage() {
	for {
		_, respBytes, err := s.WsConn.ReadMessage()
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error in consuming message from smart api websocket")
			logger.WebSocketLog.Error().Err(err).Msg("Error in consuming message from smart api websocket")

			s.WsConn.Close()

			for {
				logger.Log.Info().Msg("Attempting to reconnect to smart api websocket...")
				logger.WebSocketLog.Info().Msg("Attempting to reconnect to smart api websocket...")

				err := s.Connect()
				if err == nil {
					logger.Log.Info().Msg("Reconnected to smart api websocket successfully. Some messages may have been lost.")
					logger.WebSocketLog.Info().Msg("Reconnected to smart api websocket successfully. Some messages may have been lost.")

					SubscribeInstrumentsOnStartup()

					break
				}

				logger.Log.Error().Err(err).Msg("Failed to reconnect. Retrying in 5 seconds...")
				logger.WebSocketLog.Error().Err(err).Msg("Failed to reconnect. Retrying in 5 seconds...")

				time.Sleep(5 * time.Second)
			}

			continue
		}

		// Process the binary response
		parsedData := processResponse(respBytes)

		if parsedData != nil && SmartApiDataManager != nil {
			go SmartApiDataManager.PushMessageForDistribution(parsedData)
		}
	}
}

func (s *SmartApiApplication) sendHeartbeats() {
	for {
		logger.Log.Debug().Msg("Sending heartbeats to the smart api websocket...")
		logger.WebSocketLog.Debug().Msg("Sending heartbeats to the smart api websocket...")

		time.Sleep(heartBeatInterval)
		logger.Log.Debug().Msg("Sent heartbeats to the smart api websocket!")

		if err := s.WsConn.WriteMessage(websocket.TextMessage, []byte(heartBeatMessage)); err != nil {
			logger.Log.Error().Err(err).Msg("error in sending heartbeats to the smart api websocket")
			return
		}
	}
}

func (s *SmartApiApplication) Subscribe(subRequest Subscription) error {
	reqBytes, err := json.Marshal(subRequest)
	if err != nil {
		logger.Log.Error().Str("request", string(reqBytes)).Msg("Error in decoding subscription request")
		return err
	}

	logger.Log.Info().Msg("Sending Subscribe to the smart api websocket")

	err = s.WsConn.WriteMessage(websocket.TextMessage, reqBytes)
	if err != nil {
		logger.Log.Error().Str("request", string(reqBytes)).Msg("Error in writing subscription request to smart api websocket")
		return err
	}

	return nil
}

func (s *SmartApiApplication) Unsubscribe(subRequest Subscription) error {
	reqBytes, err := json.Marshal(subRequest)
	if err != nil {
		logger.Log.Error().Str("request", string(reqBytes)).Msg("Error in decoding subscription request")
		return err
	}

	logger.Log.Info().Msg("Sending Unsubscribe to the smart api websocket")
	err = s.WsConn.WriteMessage(websocket.TextMessage, reqBytes)
	if err != nil {
		logger.Log.Error().Str("request", string(reqBytes)).Msg("Error in writing subscription request to smart api websocket")
		return err
	}

	return nil
}
