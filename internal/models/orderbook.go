package models

type OrderBook struct {
	Asks []OrderBookEntry
	Bids []OrderBookEntry
}

type OrderBookEntry struct {
	Price    float64
	Quantity float64
}
