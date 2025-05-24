package exchange

import "context"

// OrderBook struct for order book data
type OrderBook struct {
	Bids [][]float64
	Asks [][]float64
}

// Exchange interface: all exchanges must implement these
type Exchange interface {
	CreateOrder(ctx context.Context, symbol string, quantity float64, price float64) error
	CancelOrder(ctx context.Context, orderID string) error
	GetBalance(ctx context.Context) (map[string]float64, error)
	GetOrderBook(ctx context.Context, symbol string) (OrderBook, error)
}
