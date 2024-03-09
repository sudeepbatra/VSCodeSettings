package smartapi

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"time"

	"github.com/sudeepbatra/alpha-hft/logger"
)

func updateSmartApiState(client *SmartApiApplication, key string) error {
	existingStateData, err := os.ReadFile("state.json")
	var jsonData map[string]SmartApiApplication // Declare the jsonData map

	if err != nil {
		jsonData = make(map[string]SmartApiApplication)
		jsonData[key] = *client
		logger.Log.Error().Err(err).Msg("No state file found creating one")
	} else {
		if err := json.Unmarshal(existingStateData, &jsonData); err != nil {
			logger.Log.Error().Err(err).Msg("error in decoding the smartapi state file")
			return err
		}

		jsonData[key] = *client
	}

	updatedStateData, err := json.Marshal(jsonData)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in encoding the smartapi updated data")
		return err
	}

	err = os.WriteFile("state.json", updatedStateData, 0644)

	if err != nil {
		logger.Log.Error().Err(err).Msg("error in saving the smartapi state file")
		return err
	}

	return nil

}

func processResponse(respBytes []byte) *TickParsedData {
	var parsedData TickParsedData
	logger.Log.Info().Str("msg", string(respBytes)).Msg("received response from websocket")

	if bytes.Equal(respBytes, []byte("pong")) {
		logger.Log.Debug().Str("msg", string(respBytes)).Msg("received heartbeat response from websocket")
		return nil
	}

	parsedData.SubscriptionMode = respBytes[0]
	parsedData.ExchangeType = respBytes[1]
	parsedData.Token = parseTokenValue(respBytes[2:27])
	parsedData.SequenceNumber, _ = unpackData(respBytes[27:35])

	exchangeTimeStamp, err := unpackData(respBytes[35:43])
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in unpacking exchange timestamp")
	}

	exchangeTimestampInMilliseconds := int64(exchangeTimeStamp)
	exchangeTimeseconds := exchangeTimestampInMilliseconds / 1000
	exchangeTimenanoseconds := (exchangeTimestampInMilliseconds % 1000) * 1e6

	exchangeTimestampInTime := time.Unix(exchangeTimeseconds, exchangeTimenanoseconds)
	parsedData.ExchangeTimestamp = exchangeTimestampInTime

	parsedData.LastTradedPrice, _ = unpackData(respBytes[43:51])
	parsedData.LastTradedPriceFloat = float64(parsedData.LastTradedPrice) / 100

	if parsedData.SubscriptionMode == Quote || parsedData.SubscriptionMode == SnapQuote {
		parsedData.LastTradedQuantity, _ = unpackData(respBytes[51:59])
		parsedData.AverageTradedPrice, _ = unpackData(respBytes[59:67])
		parsedData.VolumeTradeForTheDay, _ = unpackData(respBytes[67:75])
		parsedData.TotalBuyQuantity, _ = unpackDoubleData(respBytes[75:83])
		parsedData.TotalSellQuantity, _ = unpackDoubleData(respBytes[83:91])
		parsedData.OpenPriceOfTheDay, _ = unpackData(respBytes[91:99])
		parsedData.HighPriceOfTheDay, _ = unpackData(respBytes[99:107])
		parsedData.LowPriceOfTheDay, _ = unpackData(respBytes[107:115])
		parsedData.ClosedPrice, _ = unpackData(respBytes[115:123])
	}

	if parsedData.SubscriptionMode == SnapQuote {
		lastTradedTimestamp, err := unpackData(respBytes[123:131])
		if err != nil {
			logger.Log.Error().Err(err).Msg("error in unpacking last traded timestamp")
		}

		parsedData.LastTradedTimestamp = time.Unix(int64(lastTradedTimestamp), 0)

		parsedData.OpenInterest, _ = unpackData(respBytes[131:139])
		parsedData.OpenInterestChangePercentage, _ = unpackData(respBytes[139:147])
		unpackBestFiveData(respBytes[147:347])
		parsedData.UpperCircuitLimit, _ = unpackData(respBytes[347:355])
		parsedData.LowerCircuitLimit, _ = unpackData(respBytes[355:363])
		parsedData.Week52HighPrice, _ = unpackData(respBytes[363:371])
		parsedData.Week52LowPrice, _ = unpackData(respBytes[371:379])
	}

	parsedDataJson, _ := json.Marshal(parsedData)

	logger.Log.Trace().Str("msg", string(parsedDataJson)).Msg("Received market data")

	return &parsedData
}

func unpackBestFiveData(data []byte) ([]BestFivePacket, []BestFivePacket, []BestFivePacket, error) {
	var bestFiveData []BestFivePacket

	var buyData []BestFivePacket
	var sellData []BestFivePacket

	packetSize := 20 // Each packet is 20 bytes

	for i := 0; i < len(data); i += packetSize {
		packet := data[i : i+packetSize]

		// Extracting fields from the packet
		buySellFlag := int16(binary.LittleEndian.Uint16(packet[0:2]))
		quantity := int64(binary.LittleEndian.Uint64(packet[2:10]))
		price := int64(binary.LittleEndian.Uint64(packet[10:18]))
		numOfOrders := int16(binary.LittleEndian.Uint16(packet[18:20]))

		// Creating a BestFivePacket struct and appending to the result
		bestFivePacket := BestFivePacket{
			BuySellFlag: buySellFlag,
			Quantity:    quantity,
			Price:       price,
			NumOfOrders: numOfOrders,
		}

		if buySellFlag == 1 {
			buyData = append(buyData, bestFivePacket)
		} else {
			sellData = append(sellData, bestFivePacket)
		}

		bestFiveData = append(bestFiveData, bestFivePacket)
	}

	logger.Log.Trace().Str("msg", string(bestFiveData)).Msg("Received best five data")
	return bestFiveData, buyData, sellData, nil
}

func parseTokenValue(binaryPacket []byte) string {
	var token bytes.Buffer

	for _, b := range binaryPacket {
		if b == 0 {
			break
		}

		token.WriteByte(b)
	}

	return token.String()
}

func unpackData(binaryData []byte) (value int64, err error) {
	err = binary.Read(bytes.NewReader(binaryData), binary.LittleEndian, &value)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error in unpacking data")
	}

	return
}

func unpackDoubleData(binaryData []byte) (value float64, err error) {
	err = binary.Read(bytes.NewReader(binaryData), binary.LittleEndian, &value)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error in unpacking double data")
		return
	}

	return
}

func GetFromDateToDateForHistoricalData(populateOldHistoricData bool, maxDays int) (string, string) {
	// Get the current date in the Indian timezone (IST)
	now := time.Now().In(time.FixedZone("IST", 5*60*60+30*60)) // Indian Standard Time (IST)

	// Get the date 1 day before the current date by default
	fromDateTradingDay := now.AddDate(0, 0, -1)

	// If the maxDaysDataFlag is set to true, then get the date maxDays day before the current date
	if populateOldHistoricData {
		fromDateTradingDay = now.AddDate(0, 0, -maxDays)
	}

	fromDateMarketOpenTime := time.Date(fromDateTradingDay.Year(), fromDateTradingDay.Month(), fromDateTradingDay.Day(), 9, 15, 0, 0, now.Location())
	toDateMarketCloseTime := time.Date(now.Year(), now.Month(), now.Day(), 15, 30, 0, 0, now.Location())

	formattedOpenTime := fromDateMarketOpenTime.Format("2006-01-02 15:04")
	formattedCloseTime := toDateMarketCloseTime.Format("2006-01-02 15:04")

	return formattedOpenTime, formattedCloseTime
}

// type Best5Data struct {
// 	Flag       uint16
// 	Quantity   int64
// 	Price      int64
// 	NoOfOrders uint16
// }

// func ParseBest5BuyAndSellData(binaryData []byte) (map[string][]Best5Data, error) {
// 	best5BuyAndSellData := p.parseBest5BuyAndSellData(binaryData)
// 	parsedData := make(map[string][]Best5Data)
// 	parsedData["best_5_buy_data"] = best5BuyAndSellData["best_5_sell_data"]
// 	parsedData["best_5_sell_data"] = best5BuyAndSellData["best_5_buy_data"]
// 	return parsedData, nil
// }

// func parseBest5BuyAndSellData(binaryData []byte) map[string][]Best5Data {
// 	buySellPackets := p.splitPackets(binaryData)
// 	best5BuyData := make([]Best5Data, 0)
// 	best5SellData := make([]Best5Data, 0)

// 	for _, packet := range buySellPackets {
// 		eachData := Best5Data{
// 			Flag:       binary.BigEndian.Uint16(packet[0:2]),
// 			Quantity:   int64(binary.BigEndian.Uint64(packet[2:10])),
// 			Price:      int64(binary.BigEndian.Uint64(packet[10:18])),
// 			NoOfOrders: binary.BigEndian.Uint16(packet[18:20]),
// 		}

// 		if eachData.Flag == 0 {
// 			best5BuyData = append(best5BuyData, eachData)
// 		} else {
// 			best5SellData = append(best5SellData, eachData)
// 		}
// 	}

// 	return map[string][]Best5Data{
// 		"best_5_buy_data":  best5BuyData,
// 		"best_5_sell_data": best5SellData,
// 	}
// }

// func splitPackets(binaryPackets []byte) [][]byte {
// 	packetSize := 20
// 	numPackets := len(binaryPackets) / packetSize
// 	packets := make([][]byte, numPackets)

// 	for i := 0; i < numPackets; i++ {
// 		start := i * packetSize
// 		end := start + packetSize
// 		packets[i] = binaryPackets[start:end]
// 	}

// 	return packets
// }
