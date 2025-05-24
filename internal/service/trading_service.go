package service

import "eyeOne/internal/exchange"

// TradingService handles trading operations.
type TradingService struct {
	exchange exchange.Exchange
}

// NewTradingService creates a new TradingService instance.
func NewTradingService(exch exchange.Exchange) *TradingService {
	return &TradingService{exchange: exch}
}

// Implement methods for TradingService here...
