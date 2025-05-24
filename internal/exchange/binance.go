package exchange

import (
	"context"
	"strconv"

	binance "github.com/adshao/go-binance/v2"
)

type BinanceExchange struct {
	client *binance.Client
}

func NewBinanceExchange(apiKey, secretKey string) Exchange {
	client := binance.NewClient(apiKey, secretKey)
	return &BinanceExchange{client: client}
}

func (b *BinanceExchange) CreateOrder(ctx context.Context, symbol string, quantity float64, price float64) error {
	_, err := b.client.NewCreateOrderService().
		Symbol(symbol).
		Side(binance.SideTypeBuy). // Just example: hardcoded buy side
		Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).
		Quantity(strconv.FormatFloat(quantity, 'f', -1, 64)).
		Price(strconv.FormatFloat(price, 'f', -1, 64)).
		Do(ctx)
	return err
}

func (b *BinanceExchange) CancelOrder(ctx context.Context, orderID string) error {
	id, err := strconv.ParseInt(orderID, 10, 64)
	if err != nil {
		return err
	}
	_, err = b.client.NewCancelOrderService().
		OrderID(id).
		Do(ctx)
	return err
}

func (b *BinanceExchange) GetBalance(ctx context.Context) (map[string]float64, error) {
	account, err := b.client.NewGetAccountService().Do(ctx)
	if err != nil {
		return nil, err
	}
	balances := make(map[string]float64)
	for _, b := range account.Balances {
		free, _ := strconv.ParseFloat(b.Free, 64)
		balances[b.Asset] = free
	}
	return balances, nil
}

func (b *BinanceExchange) GetOrderBook(ctx context.Context, symbol string) (OrderBook, error) {
	ob, err := b.client.NewDepthService().Symbol(symbol).Do(ctx)
	if err != nil {
		return OrderBook{}, err
	}
	return OrderBook{
		Bids: parseDepthEntries(ob.Bids),
		Asks: parseDepthEntries(ob.Asks),
	}, nil
}

// Helper function to parse bids/asks into [][]float64
func parseDepthEntries(entries []binance.Bid) [][]float64 {
	var result [][]float64
	for _, e := range entries {
		price, _ := strconv.ParseFloat(e.Price, 64)
		quantity, _ := strconv.ParseFloat(e.Quantity, 64)
		result = append(result, []float64{price, quantity})
	}
	return result
}
