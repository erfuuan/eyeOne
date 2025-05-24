package exchange

import "context"

// Exchange defines the methods that any exchange implementation must provide.
type Exchange interface {
	CreateOrder(ctx context.Context, symbol, side, orderType string, quantity, price float64) (string, error)
	CancelOrder(ctx context.Context, symbol, orderID string) error
	GetBalance(ctx context.Context, asset string) (float64, error)
	GetOrderBook(ctx context.Context, symbol string) (OrderBook, error)
}

// OrderBook represents the order book with bids and asks.
type OrderBook struct {
	Bids [][]float64
	Asks [][]float64
}
