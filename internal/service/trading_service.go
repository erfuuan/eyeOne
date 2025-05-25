package service

import (
	"context"
	"fmt"

	"eyeOne/internal/exchange"
	"eyeOne/models"
)

type TradingService struct {
	exchanges map[exchange.ExchangeType]exchange.Exchange
}

func NewTradingService(exchanges map[exchange.ExchangeType]exchange.Exchange) *TradingService {
	return &TradingService{exchanges: exchanges}
}

func (ts *TradingService) getExchange(exType exchange.ExchangeType) (exchange.Exchange, error) {
	ex, ok := ts.exchanges[exType]
	if !ok {
		return nil, fmt.Errorf("exchange %s not found", exType)
	}
	return ex, nil
}

func (ts *TradingService) CreateOrder(ctx context.Context, exType exchange.ExchangeType, symbol, side, orderType string, quantity, price float64) (string, error) {
	ex, err := ts.getExchange(exType)
	if err != nil {
		return "", err
	}
	return ex.CreateOrder(ctx, symbol, side, orderType, quantity, price)
}

func (ts *TradingService) CancelOrder(ctx context.Context, exType exchange.ExchangeType, symbol, orderID string) error {
	ex, err := ts.getExchange(exType)
	if err != nil {
		return err
	}
	return ex.CancelOrder(ctx, symbol, orderID)
}

func (ts *TradingService) GetBalance(ctx context.Context, exType exchange.ExchangeType, asset string) (float64, error) {
	ex, err := ts.getExchange(exType)
	if err != nil {
		return 0, err
	}
	return ex.GetBalance(ctx, asset)
}

func (ts *TradingService) GetOrderBook(ctx context.Context, exType exchange.ExchangeType, symbol string) (models.OrderBook, error) {
	ex, err := ts.getExchange(exType)
	if err != nil {
		return models.OrderBook{}, err
	}
	return ex.GetOrderBook(ctx, symbol)
}
