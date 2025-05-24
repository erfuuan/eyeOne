package exchange

import (
	"context"
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"
)

// BinanceExchange implements the Exchange interface for Binance.
type BinanceExchange struct {
	client *binance.Client
}

// NewBinanceExchange creates a new instance of BinanceExchange.
func NewBinanceExchange(apiKey, secretKey string) Exchange {
	client := binance.NewClient(apiKey, secretKey)
	return &BinanceExchange{client: client}
}

// CreateOrder places a new order on Binance.
func (b *BinanceExchange) CreateOrder(ctx context.Context, symbol, side, orderType string, quantity, price float64) (string, error) {
	order, err := b.client.NewCreateOrderService().
		Symbol(symbol).
		Side(binance.SideType(side)).
		Type(binance.OrderType(orderType)).
		TimeInForce("GTC").
		Quantity(fmt.Sprintf("%f", quantity)).
		Price(fmt.Sprintf("%f", price)).
		Do(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", order.OrderID), nil
}

// CancelOrder cancels an existing order on Binance.
func (b *BinanceExchange) CancelOrder(ctx context.Context, symbol, orderID string) error {
	id, err := strconv.ParseInt(orderID, 10, 64)
	if err != nil {
		return err
	}
	_, err = b.client.NewCancelOrderService().
		Symbol(symbol).
		OrderID(id).
		Do(ctx)
	return err
}

// GetBalance retrieves the balance for a specific asset.
func (b *BinanceExchange) GetBalance(ctx context.Context, asset string) (float64, error) {
	account, err := b.client.NewGetAccountService().Do(ctx)
	if err != nil {
		return 0, err
	}
	for _, balance := range account.Balances {
		if balance.Asset == asset {
			free, err := strconv.ParseFloat(balance.Free, 64)
			if err != nil {
				return 0, err
			}
			return free, nil
		}
	}
	return 0, fmt.Errorf("asset %s not found", asset)
}

// GetOrderBook retrieves the order book for a specific symbol.
func (b *BinanceExchange) GetOrderBook(ctx context.Context, symbol string) (OrderBook, error) {
	ob, err := b.client.NewDepthService().Symbol(symbol).Do(ctx)
	if err != nil {
		return OrderBook{}, err
	}
	bids := convertOrderBookEntries(ob.Bids)
	asks := convertOrderBookEntries(ob.Asks)
	return OrderBook{Bids: bids, Asks: asks}, nil
}

// convertOrderBookEntries converts Binance order book entries to a slice of float64 pairs.
func convertOrderBookEntries(entries []binance.Bid) [][]float64 {
	result := make([][]float64, len(entries))
	for i, entry := range entries {
		price, _ := strconv.ParseFloat(entry.Price, 64)
		qty, _ := strconv.ParseFloat(entry.Quantity, 64)
		result[i] = []float64{price, qty}
	}
	return result
}
