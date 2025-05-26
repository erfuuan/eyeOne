// internal/exchange/wallex.go
package exchange

import (
	"context"
	"eyeOne/models"

	"go.uber.org/zap"
)

type WallexExchange struct {
	logger *zap.Logger
	// add more configs later: apiKey, apiSecret etc.
}

func NewWallexExchange(logger *zap.Logger) *WallexExchange {
	return &WallexExchange{logger: logger}
}

func (w *WallexExchange) CreateOrder(ctx context.Context, symbol, side, orderType string, quantity, price float64) (string, error) {
	w.logger.Info("Wallex: creating order", zap.String("symbol", symbol))
	// TODO: implement actual call
	return "mock-order-id-wallex", nil
}

func (w *WallexExchange) CancelOrder(ctx context.Context, symbol, orderID string) error {
	w.logger.Info("Wallex: canceling order", zap.String("orderID", orderID))
	// TODO: implement actual call
	return nil
}

func (w *WallexExchange) GetBalance(ctx context.Context, asset string) (float64, error) {
	w.logger.Info("Wallex: getting balance", zap.String("asset", asset))
	// TODO: implement actual call
	return 100.0, nil
}

func (w *WallexExchange) GetOrderBook(ctx context.Context, symbol string) (models.OrderBook, error) {
	w.logger.Info("Wallex: getting order book", zap.String("symbol", symbol))
	// TODO: implement actual call
	return models.OrderBook{}, nil
}
