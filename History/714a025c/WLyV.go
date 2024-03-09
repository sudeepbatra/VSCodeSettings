package handler

type TradingOrchestrator struct {
	OrderProcessor    *OrderProcessor
	OrderManager      *OrderManager
	OrderBookAnalyzer *OrderBookAnalyzer
	RiskManager       *RiskManagementSystem
}
