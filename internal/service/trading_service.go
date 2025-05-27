package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"eyeOne/internal/exchange"
	"eyeOne/models"
	"eyeOne/pkg/logger"
)

type TradingService struct {
	exchanges map[exchange.ExchangeType]exchange.Exchange
	log       *zap.Logger
}

func NewTradingService(exchanges map[exchange.ExchangeType]exchange.Exchange) *TradingService {
	return &TradingService{
		exchanges: exchanges,
		log:       logger.GetLogger(),
	}
}

func (ts *TradingService) getExchange(exType exchange.ExchangeType) (exchange.Exchange, error, int) {
	ex, ok := ts.exchanges[exType]
	if !ok {
		ts.log.Error("Exchange not found",
			zap.String("exchange", string(exType)),
		)
		return nil, fmt.Errorf("exchange %s not found", exType), 500
	}
	return ex, nil, 200
}

func (ts *TradingService) CreateOrder(ctx context.Context, exType exchange.ExchangeType, symbol, side, orderType string, quantity, price float64) (string, error, int) {
	ts.log.Info("Creating order",
		zap.String("exchange", string(exType)),
		zap.String("symbol", symbol),
		zap.String("side", side),
		zap.String("orderType", orderType),
		zap.Float64("quantity", quantity),
		zap.Float64("price", price),
	)

	ex, err, statusCode := ts.getExchange(exType)
	if err != nil {
		return "", err, statusCode
	}
	orderID, err, status := ex.CreateOrder(ctx, symbol, side, orderType, quantity, price)
	if err != nil {
		ts.log.Error("Failed to create order", zap.Error(err))
		return "", err, status

	}
	return orderID, nil, status
}

func (ts *TradingService) CancelOrder(ctx context.Context, exType exchange.ExchangeType, symbol, orderID string) (error, int) {
	ts.log.Info("Canceling order",
		zap.String("exchange", string(exType)),
		zap.String("symbol", symbol),
		zap.String("orderId", orderID),
	)

	ex, err, status := ts.getExchange(exType)
	if err != nil {
		return err, status
	}
	err, status = ex.CancelOrder(ctx, symbol, orderID)
	if err != nil {
		ts.log.Error("Failed to cancel order", zap.Error(err))
	}
	return err, status
}

func (ts *TradingService) GetBalance(ctx context.Context, exType exchange.ExchangeType, asset string) (float64, error, int) {
	ts.log.Info("Getting balance",
		zap.String("exchange", string(exType)),
		zap.String("asset", asset),
	)

	ex, err, status := ts.getExchange(exType)
	if err != nil {
		return 0, err, status
	}
	balance, err, status := ex.GetBalance(ctx, asset)
	if err != nil {
		ts.log.Error("Failed to get balance", zap.Error(err))
	}
	return balance, err, status
}

func (ts *TradingService) GetOrderBook(ctx context.Context, exType exchange.ExchangeType, symbol string) (models.OrderBook, error, int) {
	ts.log.Info("Getting order book",
		zap.String("exchange", string(exType)),
		zap.String("symbol", symbol),
	)

	ex, err, status := ts.getExchange(exType)
	if err != nil {
		return models.OrderBook{}, err, status
	}
	book, err, status := ex.GetOrderBook(ctx, symbol)
	if err != nil {
		ts.log.Error("Failed to get order book", zap.Error(err))
	}
	return book, err, status
}
