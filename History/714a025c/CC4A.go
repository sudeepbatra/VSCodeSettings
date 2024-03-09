package handler

import (
	"fmt"

	"github.com/sudeepbatra/alpha-hft/common"
)

type TradingSystem struct {
	OrderProcessor    *OrderProcessor
	OrderManager      *OrderManager
	OrderBookAnalyzer *OrderBookAnalyzer
	RiskManager       *RiskManagementSystem
}

func NewTradingOrchestrator(orderProcessor *OrderProcessor, orderManager *OrderManager, orderBookAnalyzer *OrderBookAnalyzer, riskManager *RiskManager) *TradingSystem {
	return &TradingSystem{
		OrderProcessor:    orderProcessor,
		OrderManager:      orderManager,
		OrderBookAnalyzer: orderBookAnalyzer,
		RiskManager:       riskManager,
	}
}

func (to *TradingSystem) ProcessTradeCycle(alphaSignal common.AlphaSignal) {
	// 1. Check existing orders for the instrument
	existingOrders := to.OrderManager.GetOrders(alphaSignal.Token, alphaSignal.Exchange)

	// 2. Analyze order book depth for the instrument
	orderBookDepth := to.OrderBookAnalyzer.AnalyzeDepth(alphaSignal.Token, to.OrderBookDepth)

	// 3. Process orders based on the alpha signal, existing orders, and order book depth
	to.OrderProcessor.ProcessOrders(alphaSignal, existingOrders)
}

type OrderProcessor struct {
	RiskManager *RiskManagementSystem
}

func NewOrderProcessor(riskManager *RiskManagementSystem) *OrderProcessor {
	return &OrderProcessor{RiskManager: riskManager}
}

func (op *OrderProcessor) ProcessOrders(alphaSignal common.AlphaSignal, existingOrders []common.Order, orderBookDepth int) {
	risk := op.RiskManager.EvaluateRisk(alphaSignal, existingOrders, orderBookDepth)
	if risk > op.RiskManager.MaxAllowableRisk {
		// Handle excessive risk, log, or take appropriate action
		fmt.Println("Excessive risk detected. No orders will be placed.")
		return
	}

	if alphaSignal.Signal == "BUY" {
		// Example: Place a buy order
		buyOrder := Order{
			Instrument: alphaSignal.Instrument,
			Side:       "BUY",
			// Add other relevant fields
		}
		err := op.ExchangeClient.PlaceOrder(buyOrder)
		if err != nil {
			// Handle order placement error
			fmt.Println("Error placing buy order:", err)
		}
	} else if alphaSignal.Signal == "SELL" {
		// Example: Place a sell order
		sellOrder := Order{
			Instrument: alphaSignal.Instrument,
			Side:       "SELL",
			// Add other relevant fields
		}
		err := op.ExchangeClient.PlaceOrder(sellOrder)
		if err != nil {
			// Handle order placement error
			fmt.Println("Error placing sell order:", err)
		}
	} else {
		// Handle other signals or do nothing
		fmt.Println("No action taken for signal:", alphaSignal.Signal)
	}
}
