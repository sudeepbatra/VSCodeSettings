package handler

import (
	"fmt"

	"github.com/sudeepbatra/alpha-hft/common"
)

type TradingOrchestrator struct {
	OrderProcessor    *OrderProcessor
	OrderManager      *OrderManager
	OrderBookAnalyzer *OrderBookAnalyzer
	RiskManager       *RiskManagementSystem
}

func NewTradingOrchestrator(orderProcessor *OrderProcessor, orderManager *OrderManager, orderBookAnalyzer *OrderBookAnalyzer, riskManager *RiskManager) *TradingOrchestrator {
	return &TradingOrchestrator{
		OrderProcessor:    orderProcessor,
		OrderManager:      orderManager,
		OrderBookAnalyzer: orderBookAnalyzer,
		RiskManager:       riskManager,
	}
}

func (to *TradingOrchestrator) ProcessTradeCycle(alphaSignal common.AlphaSignal) {
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
}
