package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"eyeOne/config"
	"eyeOne/internal/httpclient"
	"eyeOne/models"
)

type BitpinExchange struct {
	baseURL   string
	client    *httpclient.Client
	logger    *zap.Logger
	apiKey    string
	secretKey string
}

func NewBitpinExchange(client *httpclient.Client, logger *zap.Logger, cfg *config.Config) (*BitpinExchange, error) {
	return &BitpinExchange{
		baseURL:   "https://api.bitpin.ir",
		client:    client,
		logger:    logger,
		apiKey:    cfg.BitpinAPIKey,
		secretKey: cfg.BitpinSecretKey,
	}, nil
}

func (b *BitpinExchange) GetOrderBook(ctx context.Context, symbol string) (models.OrderBook, error, int) {
	urlToken := "https://api.bitpin.ir/api/v1/usr/authenticate/"
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	data := map[string]string{
		"api_key":    b.apiKey,
		"secret_key": b.secretKey,
	}
	bodyToken, status, err := b.client.PostJSON(ctx, urlToken, data, headers)
	if err != nil {
		b.logger.Error("failed to authenticate", zap.Error(err))
		return models.OrderBook{}, err, status
	}
	if status < 200 || status >= 300 {
		b.logger.Error("authentication failed", zap.Int("status", status), zap.ByteString("body", bodyToken))
		return models.OrderBook{}, fmt.Errorf("authentication failed with status %d", status), status
	}

	var tokenResp struct {
		Access  string `json:"access"`
		Refresh string `json:"refresh"`
	}

	if err := json.Unmarshal(bodyToken, &tokenResp); err != nil {
		b.logger.Error("failed to parse token response", zap.Error(err))
		return models.OrderBook{}, err, 500
	}
	urlGetOrderBook := fmt.Sprintf("%s/api/v1/mth/orderbook/%s/", b.baseURL, symbol)

	headers = map[string]string{
		"Authorization": "Bearer " + tokenResp.Access,
		"Content-Type":  "application/json",
	}

	bodyGetOrderBook, status, err := b.client.Get(ctx, urlGetOrderBook, headers)
	if err != nil {
		b.logger.Error("failed to get order book from bitpin", zap.Error(err))
		return models.OrderBook{}, err, status
	}
	if status < 200 || status >= 300 {
		b.logger.Error("failed to get order book", zap.Int("status", status), zap.ByteString("body", bodyGetOrderBook))
		return models.OrderBook{}, fmt.Errorf("failed to get order book with status %d", status), status
	}

	var res models.BitpinOrderBookResponse
	if err := json.Unmarshal(bodyGetOrderBook, &res); err != nil {
		b.logger.Error("failed to unmarshal bitpin order book response", zap.Error(err))
		return models.OrderBook{}, err, 500
	}

	bids := models.ConvertToEntries(res.Bids)
	asks := models.ConvertToEntries(res.Asks)

	return models.OrderBook{
		Bids: bids,
		Asks: asks,
	}, nil, 200
}

func (b *BitpinExchange) CreateOrder(ctx context.Context, symbol, side, orderType string, quantity, price float64) (string, error, int) {
	b.logger.Info("Bitpin: creating order")

	tokenResp, err, statusCode := b.AuthenticateBitpin(ctx)
	if err != nil {
		b.logger.Error("authentication failed", zap.Error(err))
		return "", err, statusCode
	}

	urlOrder := "https://api.bitpin.ir/api/v1/odr/orders/"
	orderHeaders := map[string]string{
		"Authorization": "Bearer " + tokenResp.Access,
		"Content-Type":  "application/json",
	}

	orderPayload := map[string]interface{}{
		"symbol": symbol,
		"type":   orderType,
		"side":   side,
		// "base_amount":      fmt.Sprintf("%.8f", quantity),
		"base_amount": 100,

		"price":            fmt.Sprintf("%.0f", price),
		"quote_amount":     1000,
		"stop_price":       0.1,
		"oco_target_price": 0.1,
		"identifier":       uuid.NewString(),
	}

	respBody, status, err := b.client.PostJSON(ctx, urlOrder, orderPayload, orderHeaders)

	if status < 200 || status >= 300 {
		b.logger.Error("Bitpin: order creation failed", zap.Int("status", status), zap.ByteString("body", respBody))
		return "", err, status
	}

	if err != nil {
		b.logger.Error("Bitpin: failed to create order, with ", zap.Int("status", status), zap.Error(err))
		return "", err, status
	}

	var orderResp struct {
		Data struct {
			ID int64 `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &orderResp); err != nil {
		b.logger.Error("Bitpin: failed to unmarshal order response", zap.Error(err))
		return "", err, 500
	}

	orderID := fmt.Sprintf("%d", orderResp.Data.ID)
	b.logger.Info("Bitpin: order created successfully", zap.String("order_id", orderID))
	return orderID, nil, 201
}

func (b *BitpinExchange) CancelOrder(ctx context.Context, symbol, orderID string) (error, int) {
	b.logger.Info("Bitpin: canceling order", zap.String("orderID", orderID))

	tokenResp, err, statusCode := b.AuthenticateBitpin(ctx)
	if err != nil {
		b.logger.Error("authentication failed", zap.Error(err))
		return err, statusCode
	}

	url := fmt.Sprintf("https://api.bitpin.ir/api/v1/odr/orders/%s/", orderID)
	headers := map[string]string{
		"Authorization": "Bearer " + tokenResp.Access,
		"Content-Type":  "application/json",
	}

	respBody, status, err := b.client.Delete(ctx, url, headers)
	if status < 200 || status >= 300 {
		b.logger.Error("Bitpin: cancel order failed", zap.Int("status", status), zap.ByteString("body", respBody))
		return err, status
	}

	if err != nil {
		b.logger.Error("Bitpin: error canceling order", zap.Error(err))
		return err, status
	}

	b.logger.Info("Bitpin: order canceled successfully", zap.String("orderID", orderID))
	return nil, 200
}

func (b *BitpinExchange) GetBalance(ctx context.Context, asset string) (float64, error, int) {
	b.logger.Info("Bitpin: getting balance", zap.String("asset", asset))

	tokenResp, err, status := b.AuthenticateBitpin(ctx)
	if err != nil {
		b.logger.Error("authentication failed", zap.Error(err))
		return 0, err, status
	}

	urlWallet := "https://api.bitpin.ir/api/v1/wlt/wallets/"
	headers := map[string]string{
		"Authorization": "Bearer " + tokenResp.Access,
		"Content-Type":  "application/json",
	}
	body, status, err := b.client.Get(ctx, urlWallet, headers)
	if err != nil {
		b.logger.Error("failed to get wallets", zap.Error(err))
		return 0, err, status
	}
	if status < 200 || status >= 300 {
		b.logger.Error("failed to get wallets", zap.Int("status", status), zap.ByteString("body", body))
		return 0, fmt.Errorf("failed to get wallets, status %d", status), status
	}

	var wallets []struct {
		ID      int    `json:"id"`
		Asset   string `json:"asset"`
		Balance string `json:"balance"`
		Frozen  string `json:"frozen"`
		Service string `json:"service"`
	}
	if err := json.Unmarshal(body, &wallets); err != nil {
		b.logger.Error("failed to unmarshal wallets response", zap.Error(err))
		return 0, err, status
	}

	for _, w := range wallets {
		if w.Asset == asset {
			balance, err := strconv.ParseFloat(w.Balance, 64)
			if err != nil {
				b.logger.Error("failed to parse balance string", zap.String("balance", w.Balance), zap.Error(err))
				return 0, err, status
			}
			return balance, nil, status
		}
	}

	return 0, fmt.Errorf("balance for asset %s not found", asset), status
}
