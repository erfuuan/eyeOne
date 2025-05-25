package exchange

import (
	"context"
	"fmt"
	"time"

	"github.com/Kucoin/kucoin-go-sdk"
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

func (k *KucoinExchange) CreateOrder(symbol, side, orderType string, quantity, price float64) (string, error) {
	clientOid := fmt.Sprintf("%s-%d", symbol, time.Now().UnixNano())
	order, err := k.client.CreateOrder(
		context.Background(),
		symbol,
		side,
		orderType,
		quantity,
		price,
		"GTC",
		clientOid,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create order: %v", err)
	}
	return order.OrderOid, nil
}

func (k *KucoinExchange) CancelOrder(ctx context.Context, orderID string) error {
	_, err := k.client.CancelOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}
	return nil
}

func GetBalance(asset string) (float64, error) {
	account, err := client.GetAccount(context.Background())
	if err != nil {
		return 0, fmt.Errorf("failed to get account: %v", err)
	}
	for _, balance := range account.Balances {
		if balance.Currency == asset {
			return balance.Available, nil
		}
	}
	return 0, fmt.Errorf("asset %s not found", asset)
}

func GetOrderBook(symbol string) (OrderBook, error) {
	orderBook, err := client.GetOrderBook(context.Background(), symbol)
	if err != nil {
		return OrderBook{}, fmt.Errorf("failed to get order book: %v", err)
	}
	return orderBook, nil
}
