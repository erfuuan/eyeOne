package exchange

import (
	"context"
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"

	"eyeOne/internal/models"
)

type BinanceExchange struct {
	client *binance.Client
}

func NewBinanceExchange(apiKey, secretKey string) (Exchange, error) {
	client := binance.NewClient(apiKey, secretKey)
	return &BinanceExchange{client: client}, nil
}

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

func (b *BinanceExchange) GetOrderBook(ctx context.Context, symbol string) (models.OrderBook, error) {
	res, err := b.client.NewDepthService().Symbol(symbol).Do(ctx)
	if err != nil {
		return models.OrderBook{}, fmt.Errorf("failed to get order book: %w", err)
	}

	bids := make([]models.OrderBookEntry, 0, len(res.Bids))
	for _, bid := range res.Bids {
		price, err := strconv.ParseFloat(bid.Price, 64)
		if err != nil {
			return models.OrderBook{}, fmt.Errorf("failed to parse bid price: %w", err)
		}
		quantity, err := strconv.ParseFloat(bid.Quantity, 64)
		if err != nil {
			return models.OrderBook{}, fmt.Errorf("failed to parse bid quantity: %w", err)
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
			return models.OrderBook{}, fmt.Errorf("failed to parse ask price: %w", err)
		}
		quantity, err := strconv.ParseFloat(ask.Quantity, 64)
		if err != nil {
			return models.OrderBook{}, fmt.Errorf("failed to parse ask quantity: %w", err)
		}
		asks = append(asks, models.OrderBookEntry{
			Price:    price,
			Quantity: quantity,
		})
	}

	return models.OrderBook{
		Bids: bids,
		Asks: asks,
	}, nil
}
