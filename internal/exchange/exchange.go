// package exchange

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

// 	"eyeOne/models"
// )

// type Exchange interface {
// 	CreateOrder(ctx context.Context, symbol, side, orderType string, quantity, price float64) (string, error)
// 	CancelOrder(ctx context.Context, symbol, orderID string) error
// 	GetBalance(ctx context.Context, asset string) (float64, error)
// 	GetOrderBook(ctx context.Context, symbol string) (models.OrderBook, error)
// }

// type OrderBook struct {
// 	Asks []OrderBookEntry
// 	Bids []OrderBookEntry
// }

// type OrderBookEntry struct {
// 	Price    float64
// 	Quantity float64
// }

// type ExchangeType string

// const (
// 	Binance ExchangeType = "binance"
// 	KuCoin  ExchangeType = "kucoin"
// )

// var registry = make(map[ExchangeType]Exchange)

// func RegisterExchange(name ExchangeType, ex Exchange) {
// 	registry[name] = ex
// }

// func GetExchange(name ExchangeType) (Exchange, error) {
// 	ex, ok := registry[name]
// 	if !ok {
// 		return nil, fmt.Errorf("exchange not registered: %s", name)
// 	}
//

package exchange

import (
	"context"
	"fmt"

	"eyeOne/models"
)

type Exchange interface {
	CreateOrder(ctx context.Context, symbol, side, orderType string, quantity, price float64) (string, error)
	CancelOrder(ctx context.Context, symbol, orderID string) error
	GetBalance(ctx context.Context, asset string) (float64, error)
	GetOrderBook(ctx context.Context, symbol string) (models.OrderBook, error)
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
