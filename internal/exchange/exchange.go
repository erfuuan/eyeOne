package exchange

import (
	"context"
	"fmt"

	"eyeOne/models"
)

type Exchange interface {
	CreateOrder(ctx context.Context, symbol, side, orderType string, quantity, price float64) (string, error, int)
	CancelOrder(ctx context.Context, symbol, orderID string) (error, int)
	GetBalance(ctx context.Context, asset string) (float64, error, int)
	GetOrderBook(ctx context.Context, symbol string) (models.OrderBook, error, int)
}

type OrderBook struct {
	Asks []OrderBookEntry
	Bids []OrderBookEntry
}

type OrderBookEntry struct {
	Price    float64
	Quantity float64
}

type ExchangeType string

const (
	Binance ExchangeType = "binance"
	KuCoin  ExchangeType = "kucoin"
	Bitpin  ExchangeType = "bitpin"
)

var registry = make(map[ExchangeType]Exchange)

func RegisterExchange(name ExchangeType, ex Exchange) {
	registry[name] = ex
}

func GetExchange(name ExchangeType) (Exchange, error) {
	ex, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("exchange not registered: %s", name)
	}
	return ex, nil
}
