package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

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

func (b *BitpinExchange) GetOrderBook(ctx context.Context, symbol string) (models.OrderBook, error) {
	// مرحله 1: گرفتن توکن
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
		return models.OrderBook{}, err
	}
	if status < 200 || status >= 300 {
		b.logger.Error("authentication failed", zap.Int("status", status), zap.ByteString("body", bodyToken))
		return models.OrderBook{}, fmt.Errorf("authentication failed with status %d", status)
	}

	// فرض می‌کنیم ساختار پاسخ این است:
	var tokenResp struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(bodyToken, &tokenResp); err != nil {
		b.logger.Error("failed to parse token response", zap.Error(err))
		return models.OrderBook{}, err
	}
	token := tokenResp.Token

	// مرحله 2: درخواست اوردر بوک با هدر Authorization
	urlGetOrderBook := fmt.Sprintf("%s/api/v1/mth/orderbook/%s/", b.baseURL, symbol)
	headersGetOrderBook := map[string]string{
		"Authorization": "Bearer " + token,
	}

	bodyGetOrderBook, status, err := b.client.Get(ctx, urlGetOrderBook, headersGetOrderBook)
	if err != nil {
		b.logger.Error("failed to get order book from bitpin", zap.Error(err))
		return models.OrderBook{}, err
	}
	if status < 200 || status >= 300 {
		b.logger.Error("failed to get order book", zap.Int("status", status), zap.ByteString("body", bodyGetOrderBook))
		return models.OrderBook{}, fmt.Errorf("failed to get order book with status %d", status)
	}

	var res models.BitpinOrderBookResponse
	if err := json.Unmarshal(bodyGetOrderBook, &res); err != nil {
		b.logger.Error("failed to unmarshal bitpin order book response", zap.Error(err))
		return models.OrderBook{}, err
	}

	bids := models.ConvertToEntries(res.Bids)
	asks := models.ConvertToEntries(res.Asks)

	return models.OrderBook{
		Bids: bids,
		Asks: asks,
	}, nil
}

func (b *BitpinExchange) CreateOrder(ctx context.Context, symbol, side, orderType string, quantity, price float64) (string, error) {
	b.logger.Info("Bitpin: creating order")
	return "mock-order-id-bitpin", nil
}

func (b *BitpinExchange) CancelOrder(ctx context.Context, symbol, orderID string) error {
	b.logger.Info("Bitpin: canceling order", zap.String("orderID", orderID))
	return nil
}

func (b *BitpinExchange) GetBalance(ctx context.Context, asset string) (float64, error) {
	b.logger.Info("Bitpin: getting balance", zap.String("asset", asset))

	urlToken := "https://api.bitpin.ir/api/v1/usr/authenticate/"
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	data := map[string]string{
		"api_key":    b.apiKey,
		"secret_key": b.secretKey,
	}
	body, status, err := b.client.PostJSON(ctx, urlToken, data, headers)
	if err != nil {
		b.logger.Error("failed to authenticate", zap.Error(err))
		return 0, err
	}
	if status < 200 || status >= 300 {
		b.logger.Error("authentication failed", zap.Int("status", status), zap.ByteString("body", body))
		return 0, fmt.Errorf("authentication failed with status %d", status)
	}

	var tokenResp struct {
		Access  string `json:"access"`
		Refresh string `json:"refresh"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		b.logger.Error("failed to unmarshal auth response", zap.Error(err))
		return 0, err
	}
	if tokenResp.Access == "" {
		return 0, fmt.Errorf("empty token received from auth")
	}

	urlWallet := "https://api.bitpin.ir/api/v1/wlt/wallets/"
	headers = map[string]string{
		"Authorization": "Bearer " + tokenResp.Access,
		"Content-Type":  "application/json",
	}
	body, status, err = b.client.Get(ctx, urlWallet, headers)
	if err != nil {
		b.logger.Error("failed to get wallets", zap.Error(err))
		return 0, err
	}
	if status < 200 || status >= 300 {
		b.logger.Error("failed to get wallets", zap.Int("status", status), zap.ByteString("body", body))
		return 0, fmt.Errorf("failed to get wallets, status %d", status)
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
		return 0, err
	}

	for _, w := range wallets {
		if w.Asset == asset {
			balance, err := strconv.ParseFloat(w.Balance, 64)
			if err != nil {
				b.logger.Error("failed to parse balance string", zap.String("balance", w.Balance), zap.Error(err))
				return 0, err
			}
			return balance, nil
		}
	}

	return 0, fmt.Errorf("balance for asset %s not found", asset)
}
