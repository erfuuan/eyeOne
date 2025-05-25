package exchange

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Kucoin/kucoin-go-sdk"

	"eyeOne/internal/exchange"
)

type KucoinExchange struct {
	client *kucoin.ApiService
}

func NewKucoinExchange(apiKey, apiSecret, apiPassphrase string) (*KucoinExchange, error) {

	client := kucoin.NewApiService(
		kucoin.ApiKeyOption(apiKey),
		kucoin.ApiSecretOption(apiSecret),
		kucoin.ApiPassPhraseOption(apiPassphrase),
	)
	return &KucoinExchange{client: client}, nil
}

// func (k *KucoinExchange) CreateOrder(symbol, side, orderType string, quantity, price float64) (string, error) {
func (k *KucoinExchange) CreateOrder(ctx context.Context, symbol, side, orderType string, quantity, price float64) (string, error) {
	clientOid := fmt.Sprintf("%s-%d", symbol, time.Now().UnixNano())

	// Convert price and quantity to strings
	priceStr := strconv.FormatFloat(price, 'f', -1, 64)
	quantityStr := strconv.FormatFloat(quantity, 'f', -1, 64)

	// Create order model
	orderModel := &kucoin.CreateOrderModel{
		ClientOid:   clientOid,
		Side:        side,
		Symbol:      symbol,
		Type:        orderType,
		Price:       priceStr,
		Size:        quantityStr,
		TimeInForce: "GTC",
	}

	// Create order
	order, err := k.client.CreateOrder(context.Background(), orderModel)
	if err != nil {
		return "", fmt.Errorf("failed to create order: %v", err)
	}

	// Assuming the order response contains an OrderOid field
	var orderResponse struct {
		OrderOid string `json:"orderOid"`
	}

	if err := order.ReadData(&orderResponse); err != nil {
		return "", fmt.Errorf("failed to read order response: %v", err)
	}

	return orderResponse.OrderOid, nil
}

func (k *KucoinExchange) CancelOrder(ctx context.Context, symbol, orderID string) error {
	// func (k *KucoinExchange) CancelOrder(ctx context.Context, orderID string) error {
	_, err := k.client.CancelOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}
	return nil
}

// func (k *KucoinExchange) GetBalances(ctx context.Context) ([]kucoin.AccountModel, error) {
// 	// Retrieve all account balances
// 	rsp, err := k.client.Accounts(ctx, "", "")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch account balances: %w", err)
// 	}

// 	// Parse the response into a slice of AccountModel
// 	var accounts []kucoin.AccountModel
// 	if err := rsp.ReadData(&accounts); err != nil {
// 		return nil, fmt.Errorf("failed to parse account data: %w", err)
// 	}

// 	return accounts, nil
// }

func (k *KucoinExchange) GetBalance(ctx context.Context, asset string) (float64, error) {
	// Retrieve all account balances
	rsp, err := k.client.Accounts(ctx, "", "")
	if err != nil {
		return 0, fmt.Errorf("failed to fetch account balances: %w", err)
	}

	// Parse the response into a slice of AccountModel
	var accounts []kucoin.AccountModel
	if err := rsp.ReadData(&accounts); err != nil {
		return 0, fmt.Errorf("failed to parse account data: %w", err)
	}

	// Find the balance for the requested asset
	for _, account := range accounts {
		if account.Currency == strings.ToUpper(asset) {
			balance, err := strconv.ParseFloat(account.Available, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse balance for %s: %w", asset, err)
			}
			return balance, nil
		}
	}

	return 0, fmt.Errorf("asset %s not found", asset)
}

// func (k *KucoinExchange) GetOrderBook(ctx context.Context, symbol string) (*kucoin.FullOrderBookModel, error) {
// 	// Retrieve the full aggregated order book for the specified symbol
// 	rsp, err := k.client.AggregatedFullOrderBook(ctx, symbol)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch order book for %s: %w", symbol, err)
// 	}

// 	// Parse the response into a FullOrderBookModel
// 	var orderBook kucoin.FullOrderBookModel
// 	if err := rsp.ReadData(&orderBook); err != nil {
// 		return nil, fmt.Errorf("failed to parse order book data: %w", err)
// 	}

// 	return &orderBook, nil
// }

func (k *KucoinExchange) GetOrderBook(ctx context.Context, symbol string) (exchange.OrderBook, error) {
	// Retrieve the full aggregated order book
	rsp, err := k.client.AggregatedFullOrderBook(ctx, symbol)
	if err != nil {
		return exchange.OrderBook{}, fmt.Errorf("failed to fetch order book for %s: %w", symbol, err)
	}

	// Parse the response
	var kucoinOrderBook kucoin.FullOrderBookModel
	if err := rsp.ReadData(&kucoinOrderBook); err != nil {
		return exchange.OrderBook{}, fmt.Errorf("failed to parse order book data: %w", err)
	}

	// Convert KuCoin order book to generic OrderBook
	return convertKucoinOrderBook(&kucoinOrderBook), nil
}

func convertKucoinOrderBook(kb *kucoin.FullOrderBookModel) exchange.OrderBook {
	ob := exchange.OrderBook{
		Asks: make([]exchange.OrderBookEntry, 0, len(kb.Asks)),
		Bids: make([]exchange.OrderBookEntry, 0, len(kb.Bids)),
	}

	// Convert asks
	for _, ask := range kb.Asks {
		price, _ := strconv.ParseFloat(ask[0], 64)
		quantity, _ := strconv.ParseFloat(ask[1], 64)
		ob.Asks = append(ob.Asks, exchange.OrderBookEntry{
			Price:    price,
			Quantity: quantity,
		})
	}

	// Convert bids
	for _, bid := range kb.Bids {
		price, _ := strconv.ParseFloat(bid[0], 64)
		quantity, _ := strconv.ParseFloat(bid[1], 64)
		ob.Bids = append(ob.Bids, exchange.OrderBookEntry{
			Price:    price,
			Quantity: quantity,
		})
	}

	return ob
}
