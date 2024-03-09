package handler

type TradingOrchestrator struct {
	OrderProcessor    *OrderProcessor
	OrderManager      *OrderManager
	OrderBookAnalyzer *OrderBookAnalyzer
	RiskManager       *RiskManagementSystem
}

func NewTradingOrchestrator(orderProcessor *OrderProcessor, orderManager *OrderManager, orderBookAnalyzer *OrderBookAnalyzer, riskManagement *RiskManagementSystem, orderBookDepth int) *TradingOrchestrator {
	return &TradingOrchestrator{
		OrderProcessor:    orderProcessor,
		OrderManager:      orderManager,
		OrderBookAnalyzer: orderBookAnalyzer,
		RiskManagement:    riskManagement,
	}
}
