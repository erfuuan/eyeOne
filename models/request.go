package models

type CreateOrderRequest struct {
	// Symbol    string  `json:"symbol" validate:"required,alphanum,min=6,max=10"`
	// Side      string  `json:"side" validate:"required,oneof=buy sell"`
	// Quantity  float64 `json:"quantity" validate:"required,gt=0"`
	// Price     float64 `json:"price" validate:"required_if=Side buy,gt=0"`
	// OrderType string  `json:"orderType" binding:"required,oneof=limit market"`

	Symbol    string  `json:"symbol" binding:"required"`
	Side      string  `json:"side" binding:"required"`
	OrderType string  `json:"orderType" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
}

type OrderResponse struct {
	OrderID   string `json:"orderId"`
	Status    string `json:"status"`
	Symbol    string `json:"symbol"`
	Exchange  string `json:"exchange"`
	Timestamp int64  `json:"timestamp"`
}

type CancelOrderRequest struct {
	Symbol  string `json:"symbol" binding:"required"`
	OrderID string `json:"orderId" binding:"required"`
}

type GetBalanceRequest struct {
	Asset string `json:"asset" binding:"required"`
}

type GetOrderBookRequest struct {
	Symbol string `json:"symbol" binding:"required"`
	Limit  int    `json:"limit" binding:"omitempty"` // optional
}
