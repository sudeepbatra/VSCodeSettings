package handler

import "sync"

type CandlestickDistribution struct {
	mu          sync.Mutex
	subscribers []chan *Candlestick
}
