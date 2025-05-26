package exchange

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Kucoin/kucoin-go-sdk"
	"go.uber.org/zap"

	"eyeOne/models"
	"eyeOne/pkg/logger"
)

type KucoinExchange struct {
	client *kucoin.ApiService
	log    *zap.Logger
}

func NewKucoinExchange(apiKey, apiSecret, apiPassphrase string) (*KucoinExchange, error) {
	client := kucoin.NewApiService(
		kucoin.ApiKeyOption(apiKey),
		kucoin.ApiSecretOption(apiSecret),
		kucoin.ApiPassPhraseOption(apiPassphrase),
	)

	log := logger.GetLogger()
	return &KucoinExchange{client: client, log: log}, nil
}

func (k *KucoinExchange) CreateOrder(ctx context.Context, symbol, side, orderType string, quantity, price float64) (string, error) {
	clientOid := fmt.Sprintf("%s-%d", symbol, time.Now().UnixNano())

	k.log.Info("Creating Kucoin order",
		zap.String("symbol", symbol),
		zap.String("side", side),
		zap.String("orderType", orderType),
		zap.Float64("quantity", quantity),
		zap.Float64("price", price),
		zap.String("clientOid", clientOid),
	)

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
		k.log.Error("Failed to create order", zap.Error(err))
		return "", fmt.Errorf("failed to create order: %v", err)
	}

	var orderResponse struct {
		OrderOid string `json:"orderOid"`
	}

	if err := order.ReadData(&orderResponse); err != nil {
		k.log.Error("Failed to read order response", zap.Error(err))
		return "", fmt.Errorf("failed to read order response: %v", err)
	}

	k.log.Info("Order created successfully", zap.String("orderOid", orderResponse.OrderOid))
	return orderResponse.OrderOid, nil
}

func (k *KucoinExchange) CancelOrder(ctx context.Context, symbol, orderID string) error {
	k.log.Info("Cancelling Kucoin order",
		zap.String("symbol", symbol),
		zap.String("orderID", orderID),
	)

	_, err := k.client.CancelOrder(ctx, orderID)
	if err != nil {
		k.log.Error("Failed to cancel order", zap.Error(err))
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	k.log.Info("Order cancelled successfully", zap.String("orderID", orderID))
	return nil
}

func (k *KucoinExchange) GetBalance(ctx context.Context, asset string) (float64, error) {
	k.log.Info("Fetching Kucoin balance", zap.String("asset", asset))

	rsp, err := k.client.Accounts(ctx, "", "")
	if err != nil {
		k.log.Error("Failed to fetch account balances", zap.Error(err))
		return 0, fmt.Errorf("failed to fetch account balances: %w", err)
	}

	var accounts []kucoin.AccountModel
	if err := rsp.ReadData(&accounts); err != nil {
		k.log.Error("Failed to parse account data", zap.Error(err))
		return 0, fmt.Errorf("failed to parse account data: %w", err)
	}

	for _, account := range accounts {
		if account.Currency == strings.ToUpper(asset) {
			balance, err := strconv.ParseFloat(account.Available, 64)
			if err != nil {
				k.log.Error("Failed to parse balance", zap.String("asset", asset), zap.Error(err))
				return 0, fmt.Errorf("failed to parse balance for %s: %w", asset, err)
			}
			k.log.Info("Balance fetched", zap.String("asset", asset), zap.Float64("balance", balance))
			return balance, nil
		}
	}

	k.log.Warn("Asset not found", zap.String("asset", asset))
	return 0, fmt.Errorf("asset %s not found", asset)
}

func (k *KucoinExchange) GetOrderBook(ctx context.Context, symbol string) (models.OrderBook, error) {
	k.log.Info("Fetching Kucoin order book", zap.String("symbol", symbol))

	rsp, err := k.client.AggregatedFullOrderBook(ctx, symbol)
	if err != nil {
		k.log.Error("Failed to fetch order book", zap.Error(err))
		return models.OrderBook{}, fmt.Errorf("failed to fetch order book: %w", err)
	}

	var kucoinOB kucoin.FullOrderBookModel
	if err := rsp.ReadData(&kucoinOB); err != nil {
		k.log.Error("Failed to parse order book", zap.Error(err))
		return models.OrderBook{}, fmt.Errorf("failed to parse order book: %w", err)
	}

	k.log.Info("Order book fetched successfully", zap.Int("asksCount", len(kucoinOB.Asks)), zap.Int("bidsCount", len(kucoinOB.Bids)))

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
