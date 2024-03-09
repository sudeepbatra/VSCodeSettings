package handler

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

func (to *TradingOrchestrator) ProcessTradeCycle(alphaSignal AlphaSignal) {
	// 1. Check existing orders for the instrument
	existingOrders := to.OrderManager.GetOrders(alphaSignal.Instrument)

	// 2. Analyze order book depth for the instrument
	orderBookDepth := to.OrderBookAnalyzer.AnalyzeDepth(alphaSignal.Instrument, to.OrderBookDepth)

	// 3. Process orders based on the alpha signal, existing orders, and order book depth
	to.OrderProcessor.ProcessOrders(alphaSignal, existingOrders, orderBookDepth)
}
