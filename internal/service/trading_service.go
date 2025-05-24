package service

import (
	"context"
	"fmt"

	"eyeOne/internal/exchange"
)

type TradingService struct {
	exchanges map[exchange.ExchangeType]exchange.Exchange
}

// NewTradingService creates a new instance of TradingService with support for multiple exchanges.
func NewTradingService(exchanges map[exchange.ExchangeType]exchange.Exchange) *TradingService {
	return &TradingService{exchanges: exchanges}
}

// getExchange returns the exchange instance by name.
func (ts *TradingService) getExchange(exType exchange.ExchangeType) (exchange.Exchange, error) {
	ex, ok := ts.exchanges[exType]
	if !ok {
		return nil, fmt.Errorf("exchange %s not found", exType)
	}
	return ex, nil
}

// CreateOrder places an order on the specified exchange.
func (ts *TradingService) CreateOrder(ctx context.Context, exType exchange.ExchangeType, symbol, side, orderType string, quantity, price float64) (string, error) {
	ex, err := ts.getExchange(exType)
	if err != nil {
		return "", err
	}
	return ex.CreateOrder(ctx, symbol, side, orderType, quantity, price)
}

// CancelOrder cancels an order on the specified exchange.
func (ts *TradingService) CancelOrder(ctx context.Context, exType exchange.ExchangeType, symbol, orderID string) error {
	ex, err := ts.getExchange(exType)
	if err != nil {
		return err
	}
	return ex.CancelOrder(ctx, symbol, orderID)
}

// GetBalance retrieves balance for a specific asset from the specified exchange.
func (ts *TradingService) GetBalance(ctx context.Context, exType exchange.ExchangeType, asset string) (float64, error) {
	ex, err := ts.getExchange(exType)
	if err != nil {
		return 0, err
	}
	return ex.GetBalance(ctx, asset)
}

// GetOrderBook retrieves the order book for a specific symbol from the specified exchange.
func (ts *TradingService) GetOrderBook(ctx context.Context, exType exchange.ExchangeType, symbol string) (exchange.OrderBook, error) {
	ex, err := ts.getExchange(exType)
	if err != nil {
		return exchange.OrderBook{}, err
	}
	return ex.GetOrderBook(ctx, symbol)
}
