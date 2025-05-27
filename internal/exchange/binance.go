package exchange

import (
	"context"
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"go.uber.org/zap"

	"eyeOne/models"
	"eyeOne/pkg/logger"
)

type BinanceExchange struct {
	client *binance.Client
	log    *zap.Logger
}

func NewBinanceExchange(apiKey, secretKey string) (Exchange, error) {
	client := binance.NewClient(apiKey, secretKey)
	log := logger.GetLogger()
	log.Info("Initialized Binance client")

	return &BinanceExchange{client: client, log: log}, nil
}

func (b *BinanceExchange) CreateOrder(ctx context.Context, symbol, side, orderType string, quantity, price float64) (string, error, int) {
	order, err := b.client.NewCreateOrderService().
		Symbol(symbol).
		Side(binance.SideType(side)).
		Type(binance.OrderType(orderType)).
		TimeInForce("GTC").
		Quantity(fmt.Sprintf("%f", quantity)).
		Price(fmt.Sprintf("%f", price)).
		Do(ctx)
	if err != nil {
		b.log.Error("Failed to create order", zap.String("symbol", symbol), zap.Error(err))
		return "", err, 500
	}
	b.log.Info("Order created", zap.String("symbol", symbol), zap.Int64("orderId", order.OrderID))
	return fmt.Sprintf("%d", order.OrderID), nil, 201
}

func (b *BinanceExchange) CancelOrder(ctx context.Context, symbol, orderID string) (error, int) {
	id, err := strconv.ParseInt(orderID, 10, 64)
	if err != nil {
		b.log.Warn("Invalid order ID format", zap.String("orderId", orderID), zap.Error(err))
		return err, 500
	}
	_, err = b.client.NewCancelOrderService().
		Symbol(symbol).
		OrderID(id).
		Do(ctx)
	if err != nil {
		b.log.Error("Failed to cancel order", zap.String("symbol", symbol), zap.Int64("orderId", id), zap.Error(err))
	}
	return err, 500
}

func (b *BinanceExchange) GetBalance(ctx context.Context, asset string) (float64, error, int) {
	account, err := b.client.NewGetAccountService().Do(ctx)
	if err != nil {
		b.log.Error("Failed to get account info", zap.Error(err))
		return 0, err, 500
	}
	for _, balance := range account.Balances {
		if balance.Asset == asset {
			free, err := strconv.ParseFloat(balance.Free, 64)
			if err != nil {
				b.log.Error("Failed to parse balance", zap.String("asset", asset), zap.String("value", balance.Free), zap.Error(err))
				return 0, err, 500
			}
			return free, nil, 200
		}
	}
	b.log.Warn("Asset not found", zap.String("asset", asset))
	return 0, fmt.Errorf("asset %s not found", asset), 404
}

func (b *BinanceExchange) GetOrderBook(ctx context.Context, symbol string) (models.OrderBook, error, int) {
	res, err := b.client.NewDepthService().Symbol(symbol).Do(ctx)
	if err != nil {
		b.log.Error("Failed to get order book", zap.String("symbol", symbol), zap.Error(err))
		return models.OrderBook{}, fmt.Errorf("failed to get order book: %w", err), 500
	}

	bids := make([]models.OrderBookEntry, 0, len(res.Bids))
	for _, bid := range res.Bids {
		price, err := strconv.ParseFloat(bid.Price, 64)
		if err != nil {
			b.log.Error("Failed to parse bid price", zap.String("price", bid.Price), zap.Error(err))
			return models.OrderBook{}, fmt.Errorf("failed to parse bid price: %w", err), 500
		}
		quantity, err := strconv.ParseFloat(bid.Quantity, 64)
		if err != nil {
			b.log.Error("Failed to parse bid quantity", zap.String("quantity", bid.Quantity), zap.Error(err))
			return models.OrderBook{}, fmt.Errorf("failed to parse bid quantity: %w", err), 500
		}
		bids = append(bids, models.OrderBookEntry{
			Price:    price,
			Quantity: quantity,
		})
	}

	asks := make([]models.OrderBookEntry, 0, len(res.Asks))
	for _, ask := range res.Asks {
		price, err := strconv.ParseFloat(ask.Price, 64)
		if err != nil {
			b.log.Error("Failed to parse ask price", zap.String("price", ask.Price), zap.Error(err))
			return models.OrderBook{}, fmt.Errorf("failed to parse ask price: %w", err), 500
		}
		quantity, err := strconv.ParseFloat(ask.Quantity, 64)
		if err != nil {
			b.log.Error("Failed to parse ask quantity", zap.String("quantity", ask.Quantity), zap.Error(err))
			return models.OrderBook{}, fmt.Errorf("failed to parse ask quantity: %w", err), 500
		}
		asks = append(asks, models.OrderBookEntry{
			Price:    price,
			Quantity: quantity,
		})
	}

	b.log.Info("Fetched order book", zap.String("symbol", symbol), zap.Int("bids", len(bids)), zap.Int("asks", len(asks)))
	return models.OrderBook{
		Bids: bids,
		Asks: asks,
	}, nil, 200
}
