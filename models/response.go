package models

type SuccessResponse struct {
	StatusCode int    `json:"statusCode"`
	Data       any    `json:"data,omitempty"`
	Message    string `json:"message,omitempty"`
	Timestamp  int64  `json:"timestamp"`
}

type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Timestamp  int64  `json:"timestamp"`
}

type ErrorPayload struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Timestamp  int64  `json:"timestamp"`
}

type OrderDataResponse struct {
	OrderID  string  `json:"orderId"`
	Exchange string  `json:"exchange"`
	Symbol   string  `json:"symbol"`
	Side     string  `json:"side"`
	Type     string  `json:"type"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}

type BalanceDataResponse struct {
	Asset   string  `json:"asset"`
	Balance float64 `json:"balance"`
}
