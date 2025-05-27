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

func (b *BitpinExchange) AuthenticateBitpin(ctx context.Context) (struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}, error, int) {
	url := fmt.Sprintf("%s/api/v1/usr/authenticate/", b.baseURL)
	data := map[string]string{
		"api_key":    b.apiKey,
		"secret_key": b.secretKey,
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	body, status, err := b.client.PostJSON(ctx, url, data, headers)
	if err != nil {
		b.logger.Error("bitpin authentication failed", zap.Error(err))
		return struct {
			Access  string `json:"access"`
			Refresh string `json:"refresh"`
		}{}, err, status
	}
	if status < 200 || status >= 300 {
		b.logger.Error("bitpin authentication failed with status", zap.Int("status", status), zap.ByteString("body", body))
		return struct {
			Access  string `json:"access"`
			Refresh string `json:"refresh"`
		}{}, fmt.Errorf("authentication failed with status %d", status), status
	}

	var tokenResp struct {
		Access  string `json:"access"`
		Refresh string `json:"refresh"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		b.logger.Error("failed to unmarshal auth response", zap.Error(err))
		return tokenResp, err, 500
	}
	return tokenResp, nil, 200
}

func (b *BitpinExchange) GetOrderBook(ctx context.Context, symbol string) (models.OrderBook, error, int) {
	tokenResp, err, status := b.AuthenticateBitpin(ctx)
	if err != nil {
		return models.OrderBook{}, err, status
	}

	url := fmt.Sprintf("%s/api/v1/mth/orderbook/%s/", b.baseURL, symbol)
	headers := map[string]string{
		"Authorization": "Bearer " + tokenResp.Access,
		"Content-Type":  "application/json",
	}
	body, status, err := b.client.Get(ctx, url, headers)
	if err != nil {
		b.logger.Error("failed to get order book", zap.Error(err))
		return models.OrderBook{}, err, status
	}
	if status < 200 || status >= 300 {
		b.logger.Error("get order book failed", zap.Int("status", status), zap.ByteString("body", body))
		return models.OrderBook{}, fmt.Errorf("order book error %d", status), status
	}

	var res models.BitpinOrderBookResponse
	if err := json.Unmarshal(body, &res); err != nil {
		b.logger.Error("failed to parse order book response", zap.Error(err))
		return models.OrderBook{}, err, 500
	}

	return models.OrderBook{
		Bids: models.ConvertToEntries(res.Bids),
		Asks: models.ConvertToEntries(res.Asks),
	}, nil, 200
}

func (b *BitpinExchange) CreateOrder(ctx context.Context, symbol, side, orderType string, quantity, price float64) (string, error, int) {
	tokenResp, err, status := b.AuthenticateBitpin(ctx)
	if err != nil {
		return "", err, status
	}

	url := fmt.Sprintf("%s/api/v1/odr/orders/", b.baseURL)
	headers := map[string]string{
		"Authorization": "Bearer " + tokenResp.Access,
		"Content-Type":  "application/json",
	}
	payload := map[string]interface{}{
		"symbol":           symbol,
		"type":             orderType,
		"side":             side,
		"base_amount":      fmt.Sprintf("%.8f", quantity),
		"price":            fmt.Sprintf("%.0f", price),
		"quote_amount":     1000,
		"stop_price":       0.1,
		"oco_target_price": 0.1,
		"identifier":       uuid.NewString(),
	}

	respBody, status, err := b.client.PostJSON(ctx, url, payload, headers)
	if err != nil || status < 200 || status >= 300 {
		b.logger.Error("order creation failed", zap.Int("status", status), zap.Error(err), zap.ByteString("body", respBody))
		return "", err, status
	}

	var orderResp struct {
		Data struct {
			ID int64 `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &orderResp); err != nil {
		b.logger.Error("unmarshal order response failed", zap.Error(err))
		return "", err, 500
	}
	return fmt.Sprintf("%d", orderResp.Data.ID), nil, 201
}

func (b *BitpinExchange) CancelOrder(ctx context.Context, symbol, orderID string) (error, int) {
	tokenResp, err, status := b.AuthenticateBitpin(ctx)
	if err != nil {
		return err, status
	}

	url := fmt.Sprintf("%s/api/v1/odr/orders/%s/", b.baseURL, orderID)
	headers := map[string]string{
		"Authorization": "Bearer " + tokenResp.Access,
		"Content-Type":  "application/json",
	}
	respBody, status, err := b.client.Delete(ctx, url, headers)
	if err != nil || status < 200 || status >= 300 {
		b.logger.Error("cancel order failed", zap.Int("status", status), zap.Error(err), zap.ByteString("body", respBody))
		return err, status
	}
	return nil, 200
}

func (b *BitpinExchange) GetBalance(ctx context.Context, asset string) (float64, error, int) {
	tokenResp, err, status := b.AuthenticateBitpin(ctx)
	if err != nil {
		return 0, err, status
	}

	url := fmt.Sprintf("%s/api/v1/wlt/wallets/", b.baseURL)
	headers := map[string]string{
		"Authorization": "Bearer " + tokenResp.Access,
		"Content-Type":  "application/json",
	}
	body, status, err := b.client.Get(ctx, url, headers)
	if err != nil {
		b.logger.Error("get wallets failed", zap.Error(err))
		return 0, err, status
	}
	if status < 200 || status >= 300 {
		b.logger.Error("get wallets failed with status", zap.Int("status", status), zap.ByteString("body", body))
		return 0, fmt.Errorf("wallets failed %d", status), status
	}

	var wallets []struct {
		Asset   string `json:"asset"`
		Balance string `json:"balance"`
	}
	if err := json.Unmarshal(body, &wallets); err != nil {
		b.logger.Error("unmarshal wallets failed", zap.Error(err))
		return 0, err, 500
	}

	for _, w := range wallets {
		if w.Asset == asset {
			balance, err := strconv.ParseFloat(w.Balance, 64)
			if err != nil {
				b.logger.Error("parse balance failed", zap.String("balance", w.Balance), zap.Error(err))
				return 0, err, status
			}
			return balance, nil, status
		}
	}
	return 0, fmt.Errorf("asset %s not found", asset), status
}
