package handler

import (
	"fmt"
	"sync"
	"time"

	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/sudeepbatra/alpha-hft/broker/smartapi"
	"github.com/sudeepbatra/alpha-hft/data"
	"github.com/sudeepbatra/alpha-hft/logger"
)

type OrderBookTree struct {
	*redblacktree.Tree
}

func (obt *OrderBookTree) UpdateOrderTree(newBestData []smartapi.Best5Data) {
	for _, item := range newBestData {
		priceKey := item.Price
		obt.Put(int(priceKey), item)
	}
}

func (obt *OrderBookTree) RemoveBookedOrders(newBestData []smartapi.Best5Data, orderType string) {
	currentOrders := obt.Values()

	if currentOrders == nil {
		return
	}

	bestPrice := newBestData[0].Price
	var ordersToRemove []smartapi.Best5Data

	for i := 0; i < len(newBestData)-1; i++ {
		startRange := newBestData[i].Price
		endRange := newBestData[i+1].Price

		for index := len(currentOrders) - 1; index >= 0; index-- {
			if item, ok := currentOrders[index].(smartapi.Best5Data); ok {
				if orderType == "BUY" {
					if item.Price > bestPrice || (item.Price > endRange && item.Price < startRange) {
						ordersToRemove = append(ordersToRemove, item)
					}
				} else if orderType == "SELL" {
					if item.Price < bestPrice || (item.Price > startRange && item.Price < endRange) {
						ordersToRemove = append(ordersToRemove, item)
					}
				}

			} else {
				return
			}
		}
	}

	if ordersToRemove == nil {
		return
	}

	for _, order := range ordersToRemove {
		obt.Remove(int(order.Price))
	}
}

type OrderBook struct {
	BuyOrders       *OrderBookTree
	SellOrders      *OrderBookTree
	BuyMarketDepth  map[int64]int64
	SellMarketDepth map[int64]int64
	BestBid         int64
	BestAsk         int64
}

func (ob *OrderBook) UpdateMarketDepth() {
	buyOrders := ob.BuyOrders.Values()
	sellOrders := ob.SellOrders.Values()

	var buyMarketDepth int64
	var sellMarketDepth int64

	for index := len(buyOrders) - 1; index >= 0; index-- {
		if item, ok := buyOrders[index].(smartapi.Best5Data); ok {
			buyMarketDepth += int64(item.NumOfOrders) * item.Quantity
			ob.BuyMarketDepth[item.Price] = buyMarketDepth
		}
	}

	for index := 0; index < len(sellOrders); index++ {
		if item, ok := sellOrders[index].(smartapi.Best5Data); ok {
			sellMarketDepth += int64(item.NumOfOrders) * item.Quantity
			ob.SellMarketDepth[item.Price] = sellMarketDepth
		}
	}

}

func (ob *OrderBook) UpdateOrderBook(newBestData []smartapi.Best5Data, orderType string) {

	if len(newBestData) == 0 {
		return
	}

	bestOrder := newBestData[0]

	if orderType == "BUY" {
		ob.BestBid = bestOrder.Price
		ob.BuyOrders.RemoveBookedOrders(newBestData, orderType)
		ob.BuyOrders.UpdateOrderTree(newBestData)
	} else if orderType == "SELL" {
		ob.BestAsk = bestOrder.Price
		ob.SellOrders.RemoveBookedOrders(newBestData, orderType)
		ob.SellOrders.UpdateOrderTree(newBestData)
	}
}

type OrderBookManager struct {
	OrderBooks           map[string]map[string]*OrderBook
	OpenOrders           map[string]map[string]*smartapi.Order // to be migrated in a limit order book with similar structure above
	OpenOrdersUpdateLock sync.Mutex
}

func NewOrderBookManager() *OrderBookManager {
	return &OrderBookManager{
		OrderBooks: make(map[string]map[string]*OrderBook),
		OpenOrders: make(map[string]map[string]*smartapi.Order),
	}
}

func (obm *OrderBookManager) ProcessOrderBookUpdates() {
	parsedTickData := smartapi.SmartApiDataManager.Subscribe()

	for tickData := range parsedTickData {
		if tickData == nil {
			logger.Log.Error().
				Str("parser", "order book from tick").
				Msg("channel is closed. returning from ProcessTicks!")

			break
		}
		exchange := smartapi.CodeExchangeTypes[int(tickData.ExchangeType)]

		_, tokenExists := obm.OrderBooks[tickData.Token]

		if !tokenExists {
			obm.OrderBooks[tickData.Token] = make(map[string]*OrderBook)
		}

		_, exchangeExist := obm.OrderBooks[tickData.Token][exchange]

		if !exchangeExist {
			obm.OrderBooks[tickData.Token][exchange] = &OrderBook{
				BuyOrders:       &OrderBookTree{redblacktree.NewWithIntComparator()},
				SellOrders:      &OrderBookTree{redblacktree.NewWithIntComparator()},
				BuyMarketDepth:  make(map[int64]int64),
				SellMarketDepth: make(map[int64]int64),
			}
		}

		obm.OrderBooks[tickData.Token][exchange].UpdateOrderBook(tickData.Best5BuyData, "BUY")
		obm.OrderBooks[tickData.Token][exchange].UpdateOrderBook(tickData.Best5SellData, "SELL")
		obm.OrderBooks[tickData.Token][exchange].UpdateMarketDepth()

		logger.Log.Debug().Str("token", tickData.Token).
			Str("exchange", exchange).
			Interface("BestBuyData", tickData.Best5BuyData).
			Interface("BestSellData", tickData.Best5SellData).
			Interface("OrderBook", obm.OrderBooks[tickData.Token][exchange].SellOrders).
			Msg("order book updated for the given token")
	}

}

func (obm *OrderBookManager) getOrderBookForToken(token string, exchange string) *OrderBook {
	_, tokenExists := obm.OrderBooks[token]

	if !tokenExists {
		return nil
	}

	_, exchangeExist := obm.OrderBooks[token][exchange]

	if !exchangeExist {
		return nil
	}

	return obm.OrderBooks[token][exchange]

}

func (obm *OrderBookManager) getBestPriceAvailiable(order *smartapi.Order, orderBook *OrderBook) string {
	var bestPrice float64
	if order.OrderType == "BUY" {
		bestPrice = float64(orderBook.BestAsk) / 100
	} else {
		bestPrice = float64(orderBook.BestBid) / 100
	}

	return fmt.Sprintf("%.2f", bestPrice)
}

func (obm *OrderBookManager) ProcessAlphaSignals() {
	alphaSignalChannel := AlphaSignalManager.Subscribe()
	for alphaSignal := range alphaSignalChannel {

		//temporary more dynamic solution to be put later for order filtering

		_, present := smartapi.Nifty50Instruments[alphaSignal.Symbol]

		if alphaSignal.Exchange != "NSE" || !present {
			logger.Log.Debug().Str("symbol", alphaSignal.Symbol).
				Str("exchange", alphaSignal.Exchange).
				Msg("Ignoring order as for given exchange and symbol is not active.")
		}

		transactionType := "BUY"
		if alphaSignal.Signal == "SHORT" {
			transactionType = "SELL"
		}

		newOrder := &smartapi.Order{
			Variety:          "NORMAL",
			SymbolToken:      alphaSignal.Token,
			TradingSymbol:    alphaSignal.Symbol,
			Exchange:         alphaSignal.Exchange,
			TransactionType:  transactionType,
			OrderType:        "LIMIT",
			ProductType:      "INTRADAY",
			Duration:         "DAY",
			Price:            fmt.Sprintf("%.2f", alphaSignal.Price),
			SquareOff:        "0",
			StopLoss:         "0",
			Quantity:         "1",
			OrderInitiatedAt: time.Now(),
		}
		obm.OpenOrdersUpdateLock.Lock()
		_, ok := obm.OpenOrders[alphaSignal.Exchange]
		if !ok {
			obm.OpenOrders[alphaSignal.Exchange] = make(map[string]*smartapi.Order)

		} else {
			_, tokenExists := obm.OpenOrders[alphaSignal.Exchange][alphaSignal.Token]
			if tokenExists {
				logger.Log.Info().
					Str("token", alphaSignal.Token).
					Str("exchange", alphaSignal.Exchange).
					Msg("order already exists for the given token and exchange ignoring")
				obm.OpenOrdersUpdateLock.Unlock()
				continue
			}
		}
		obm.OpenOrders[alphaSignal.Exchange][alphaSignal.Token] = newOrder
		logger.Log.Info().Interface("order", newOrder).Msg("received alpha signal converting to order")
		obm.OpenOrdersUpdateLock.Unlock()
	}
}

func (obm *OrderBookManager) PlaceOpenOrder(order *smartapi.Order) {
	placeOrderResponse, err := smartapi.SmartApiBrokers["trading"].PlaceOrderV2(order)

	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("token", order.SymbolToken).
			Str("exchange", order.Exchange).
			Str("transactionType", order.TransactionType).
			Str("orderPrice", order.Price).
			Msg("error in placing the order")
	} else {
		if placeOrderResponse.Status {
			order.OrderID = placeOrderResponse.Data.OrderID
			order.UniqueOrderID = placeOrderResponse.Data.UniqueOrderID
			order.Status = "pending"
		} else {
			order.Status = "rejected"
			order.Text = placeOrderResponse.Message
		}
		err = data.InsertRecord(data.OrderTable, *order)
		if err != nil {
			logger.Log.Error().
				Err(err).
				Msg("error in inserting the order")
		}

	}
}

func (obm *OrderBookManager) ProcessOrders() {
	for {
		obm.OpenOrdersUpdateLock.Lock()
		for exchange, orders := range obm.OpenOrders {
			for token, order := range orders {
				orderBook := obm.getOrderBookForToken(token, exchange)
				if orderBook == nil {
					logger.Log.Warn().Str("token", token).
						Str("exchange", exchange).
						Str("transactionType", order.TransactionType).
						Msg("No data found in the order book! removing the order.")

					delete(obm.OpenOrders[exchange], token)
					continue
				}

				if time.Since(order.OrderInitiatedAt) > 1*time.Minute {
					logger.Log.Info().Str("token", order.SymbolToken).
						Str("exchange", order.Exchange).
						Str("symbol", order.TradingSymbol).
						Str("transactionType", order.TransactionType).
						Str("orderPrice", order.Price).
						Msg("price hunting time limit reached placing the order with the current best price.")

					obm.PlaceOpenOrder(order)
					delete(obm.OpenOrders[exchange], token)
					continue
				}

				updatedPrice := obm.getBestPriceAvailiable(order, orderBook)
				if updatedPrice != order.Price {
					order.Price = updatedPrice
					logger.Log.Info().Str("token", order.SymbolToken).
						Str("exchange", order.Exchange).
						Str("transactionType", order.TransactionType).
						Str("orderPrice", order.Price).
						Str("orderInitiatedAt", order.OrderInitiatedAt.String()).
						Msg("updating the bestPrice availiable")

				}
			}
		}

		obm.OpenOrdersUpdateLock.Unlock()
		time.Sleep(500 * time.Millisecond)
	}
}

func (obm *OrderBookManager) StartOrderManager() {
	go obm.ProcessOrderBookUpdates()
	go obm.ProcessAlphaSignals()
	go obm.ProcessOrders()
}
