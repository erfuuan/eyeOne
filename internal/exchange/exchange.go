package exchange

import (
	"context"
	"fmt"

	"eyeOne/models"
)

// Exchange defines the methods that any exchange implementation must provide.
type Exchange interface {
	CreateOrder(ctx context.Context, symbol, side, orderType string, quantity, price float64) (string, error)
	CancelOrder(ctx context.Context, symbol, orderID string) error
	GetBalance(ctx context.Context, asset string) (float64, error)
	// GetOrderBook(ctx context.Context, symbol string) (OrderBook, error)
	GetOrderBook(ctx context.Context, symbol string) (models.OrderBook, error)
}

// OrderBook represents the order book with bids and asks.
// type OrderBook struct {
// 	Bids [][]float64
// 	Asks [][]float64
// }

type OrderBook struct {
	Asks []OrderBookEntry
	Bids []OrderBookEntry
}

type OrderBookEntry struct {
	Price    float64
	Quantity float64
}

// ExchangeType is a type-safe alias for exchange names
type ExchangeType string

const (
	Binance ExchangeType = "binance"
	KuCoin  ExchangeType = "kucoin"
)

// registry to hold all registered exchanges
var registry = make(map[ExchangeType]Exchange)

// RegisterExchange allows each exchange to register itself
func RegisterExchange(name ExchangeType, ex Exchange) {
	registry[name] = ex
}

// GetExchange returns the exchange implementation by name
func GetExchange(name ExchangeType) (Exchange, error) {
	ex, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("exchange not registered: %s", name)
	}
	return ex, nil
}
