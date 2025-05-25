package exchange

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Kucoin/kucoin-go-sdk"

	"eyeOne/models"
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

func (k *KucoinExchange) CreateOrder(ctx context.Context, symbol, side, orderType string, quantity, price float64) (string, error) {
	clientOid := fmt.Sprintf("%s-%d", symbol, time.Now().UnixNano())

	priceStr := strconv.FormatFloat(price, 'f', -1, 64)
	quantityStr := strconv.FormatFloat(quantity, 'f', -1, 64)

	orderModel := &kucoin.CreateOrderModel{
		ClientOid:   clientOid,
		Side:        side,
		Symbol:      symbol,
		Type:        orderType,
		Price:       priceStr,
		Size:        quantityStr,
		TimeInForce: "GTC",
	}

	order, err := k.client.CreateOrder(context.Background(), orderModel)
	if err != nil {
		return "", fmt.Errorf("failed to create order: %v", err)
	}

	var orderResponse struct {
		OrderOid string `json:"orderOid"`
	}

	if err := order.ReadData(&orderResponse); err != nil {
		return "", fmt.Errorf("failed to read order response: %v", err)
	}

	return orderResponse.OrderOid, nil
}

func (k *KucoinExchange) CancelOrder(ctx context.Context, symbol, orderID string) error {
	_, err := k.client.CancelOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}
	return nil
}

func (k *KucoinExchange) GetBalance(ctx context.Context, asset string) (float64, error) {
	rsp, err := k.client.Accounts(ctx, "", "")
	if err != nil {
		return 0, fmt.Errorf("failed to fetch account balances: %w", err)
	}

	var accounts []kucoin.AccountModel
	if err := rsp.ReadData(&accounts); err != nil {
		return 0, fmt.Errorf("failed to parse account data: %w", err)
	}

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

func (k *KucoinExchange) GetOrderBook(ctx context.Context, symbol string) (models.OrderBook, error) {
	rsp, err := k.client.AggregatedFullOrderBook(ctx, symbol)
	if err != nil {
		return models.OrderBook{}, fmt.Errorf("failed to fetch order book: %w", err)
	}

	var kucoinOB kucoin.FullOrderBookModel
	if err := rsp.ReadData(&kucoinOB); err != nil {
		return models.OrderBook{}, fmt.Errorf("failed to parse order book: %w", err)
	}

	return models.OrderBook{
		Asks: convertEntries(kucoinOB.Asks),
		Bids: convertEntries(kucoinOB.Bids),
	}, nil
}

func convertEntries(entries [][]string) []models.OrderBookEntry {
	result := make([]models.OrderBookEntry, 0, len(entries))
	for _, entry := range entries {
		if len(entry) < 2 {
			continue
		}
		price, _ := strconv.ParseFloat(entry[0], 64)
		quantity, _ := strconv.ParseFloat(entry[1], 64)
		result = append(result, models.OrderBookEntry{
			Price:    price,
			Quantity: quantity,
		})
	}
	return result
}
