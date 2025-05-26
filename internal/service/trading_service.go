// package service

// import (
// 	"context"
// 	"fmt"
//
//
//
//
//
//
//
//
// 

// 	"eyeOne/internal/exchange"
// 	"eyeOne/models"
// )

// type TradingService struct {
// 	exchanges map[exchange.ExchangeType]exchange.Exchange
// }

// func NewTradingService(exchanges map[exchange.ExchangeType]exchange.Exchange) *TradingService {
// 	return &TradingService{exchanges: exchanges}
// }

// func (ts *TradingService) getExchange(exType exchange.ExchangeType) (exchange.Exchange, error) {
// 	ex, ok := ts.exchanges[exType]
// 	if !ok {
// 		return nil, fmt.Errorf("exchange %s not found", exType)
// 	}
// 	return ex, nil
// }

// func (ts *TradingService) CreateOrder(ctx context.Context, exType exchange.ExchangeType, symbol, side, orderType string, quantity, price float64) (string, error) {
// 	ex, err := ts.getExchange(exType)
// 	if err != nil {
// 		return "", err
// 	}
// 	return ex.CreateOrder(ctx, symbol, side, orderType, quantity, price)
// }

// func (ts *TradingService) CancelOrder(ctx context.Context, exType exchange.ExchangeType, symbol, orderID string) error {
// 	ex, err := ts.getExchange(exType)
// 	if err != nil {
// 		return err
// 	}
// 	return ex.CancelOrder(ctx, symbol, orderID)
// }

// func (ts *TradingService) GetBalance(ctx context.Context, exType exchange.ExchangeType, asset string) (float64, error) {
// 	ex, err := ts.getExchange(exType)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return ex.GetBalance(ctx, asset)
// }

// func (ts *TradingService) GetOrderBook(ctx context.Context, exType exchange.ExchangeType, symbol string) (models.OrderBook, error) {
// 	ex, err := ts.getExchange(exType)
// 	if err != nil {
// 		return models.OrderBook{}, err
// 	}
// 	return ex.GetOrderBook(ctx, symbol)
// }

package service

import (
	"context"
	"fmt"

	"eyeOne/internal/exchange"
	"eyeOne/models"
	"eyeOne/pkg/logger"

	"go.uber.org/zap"
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

func (ts *TradingService) getExchange(exType exchange.ExchangeType) (exchange.Exchange, error) {
	ex, ok := ts.exchanges[exType]
	if !ok {
		ts.log.Error("Exchange not found",
			zap.String("exchange", string(exType)),
		)
		return nil, fmt.Errorf("exchange %s not found", exType)
	}
	return ex, nil
}

func (ts *TradingService) CreateOrder(ctx context.Context, exType exchange.ExchangeType, symbol, side, orderType string, quantity, price float64) (string, error) {
	ts.log.Info("Creating order",
		zap.String("exchange", string(exType)),
		zap.String("symbol", symbol),
		zap.String("side", side),
		zap.String("orderType", orderType),
		zap.Float64("quantity", quantity),
		zap.Float64("price", price),
	)

	ex, err := ts.getExchange(exType)
	if err != nil {
		return "", err
	}
	orderID, err := ex.CreateOrder(ctx, symbol, side, orderType, quantity, price)
	if err != nil {
		ts.log.Error("Failed to create order", zap.Error(err))
	}
	return orderID, err
}

func (ts *TradingService) CancelOrder(ctx context.Context, exType exchange.ExchangeType, symbol, orderID string) error {
	ts.log.Info("Canceling order",
		zap.String("exchange", string(exType)),
		zap.String("symbol", symbol),
		zap.String("orderId", orderID),
	)

	ex, err := ts.getExchange(exType)
	if err != nil {
		return err
	}
	err = ex.CancelOrder(ctx, symbol, orderID)
	if err != nil {
		ts.log.Error("Failed to cancel order", zap.Error(err))
	}
	return err
}

func (ts *TradingService) GetBalance(ctx context.Context, exType exchange.ExchangeType, asset string) (float64, error) {
	ts.log.Info("Getting balance",
		zap.String("exchange", string(exType)),
		zap.String("asset", asset),
	)

	ex, err := ts.getExchange(exType)
	if err != nil {
		return 0, err
	}
	balance, err := ex.GetBalance(ctx, asset)
	if err != nil {
		ts.log.Error("Failed to get balance", zap.Error(err))
	}
	return balance, err
}

func (ts *TradingService) GetOrderBook(ctx context.Context, exType exchange.ExchangeType, symbol string) (models.OrderBook, error) {
	ts.log.Info("Getting order book",
		zap.String("exchange", string(exType)),
		zap.String("symbol", symbol),
	)

	ex, err := ts.getExchange(exType)
	if err != nil {
		return models.OrderBook{}, err
	}
	book, err := ex.GetOrderBook(ctx, symbol)
	if err != nil {
		ts.log.Error("Failed to get order book", zap.Error(err))
	}
	return book, err
}
