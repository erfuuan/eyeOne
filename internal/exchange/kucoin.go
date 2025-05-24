package exchange

import (
	"context"
	"strconv"

	"github.com/Kucoin/kucoin-go-sdk"
)

type KuCoinExchange struct {
	client *kucoin.ApiClient
}

func NewKuCoinExchange(apiKey, secretKey string) Exchange {
	config := kucoin.NewConfig().
		WithApiKey(apiKey).
		WithApiSecret(secretKey).
		WithApiPassphrase("your-passphrase") // باید پاس‌فریز خود را وارد کنید

	client := kucoin.NewApiClient(config)
	return &KuCoinExchange{client: client}
}

func (k *KuCoinExchange) CreateOrder(ctx context.Context, symbol string, quantity float64, price float64) error {
	orderReq := kucoin.OrderRequest{
		Symbol:      symbol,
		Side:        kucoin.SideTypeBuy,
		Type:        kucoin.OrderTypeLimit,
		Price:       strconv.FormatFloat(price, 'f', -1, 64),
		Size:        strconv.FormatFloat(quantity, 'f', -1, 64),
		TimeInForce: kucoin.TimeInForceGTC,
	}

	_, err := k.client.OrderApi.CreateOrder(ctx).OrderRequest(orderReq).Execute()
	return err
}

func (k *KuCoinExchange) CancelOrder(ctx context.Context, orderID string) error {
	_, err := k.client.OrderApi.CancelOrder(ctx, orderID).Execute()
	return err
}

func (k *KuCoinExchange) GetBalance(ctx context.Context) (map[string]float64, error) {
	resp, err := k.client.AccountApi.GetAccounts(ctx).Execute()
	if err != nil {
		return nil, err
	}

	balances := make(map[string]float64)
	for _, a := range resp.Data {
		bal, _ := strconv.ParseFloat(a.Available, 64)
		balances[a.Currency] = bal
	}
	return balances, nil
}

func (k *KuCoinExchange) GetOrderBook(ctx context.Context, symbol string) (OrderBook, error) {
	resp, err := k.client.MarketApi.GetLevel2(ctx).Symbol(symbol).Execute()
	if err != nil {
		return OrderBook{}, err
	}

	bids := make([][]float64, len(resp.Data.Bids))
	for i, b := range resp.Data.Bids {
		price, _ := strconv.ParseFloat(b[0], 64)
		size, _ := strconv.ParseFloat(b[1], 64)
		bids[i] = []float64{price, size}
	}

	asks := make([][]float64, len(resp.Data.Asks))
	for i, a := range resp.Data.Asks {
		price, _ := strconv.ParseFloat(a[0], 64)
		size, _ := strconv.ParseFloat(a[1], 64)
		asks[i] = []float64{price, size}
	}

	return OrderBook{Bids: bids, Asks: asks}, nil
}
