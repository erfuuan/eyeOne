package service

import (
	"context"

	"eyeOne/internal/exchange"
)

// TradingService implements the exchange.Exchange interface.
type TradingService struct {
	exchange exchange.Exchange
}

// NewTradingService creates a new instance of TradingService.
func NewTradingService(ex exchange.Exchange) *TradingService {
	return &TradingService{exchange: ex}
}

// CreateOrder delegates the order creation to the underlying exchange.
func (ts *TradingService) CreateOrder(ctx context.Context, symbol, side, orderType string, quantity, price float64) (string, error) {
	return ts.exchange.CreateOrder(ctx, symbol, side, orderType, quantity, price)
}

// CancelOrder delegates the order cancellation to the underlying exchange.
func (ts *TradingService) CancelOrder(ctx context.Context, symbol, orderID string) error {
	return ts.exchange.CancelOrder(ctx, symbol, orderID)
}

// GetBalance retrieves the balance for a specific asset from the underlying exchange.
func (ts *TradingService) GetBalance(ctx context.Context, asset string) (float64, error) {
	return ts.exchange.GetBalance(ctx, asset)
}

// GetOrderBook retrieves the order book for a specific symbol from the underlying exchange.
func (ts *TradingService) GetOrderBook(ctx context.Context, symbol string) (exchange.OrderBook, error) {
	return ts.exchange.GetOrderBook(ctx, symbol)
}
