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

type OrderResponse struct {
	OrderID   string `json:"orderId"`
	Status    string `json:"status"`
	Symbol    string `json:"symbol"`
	Exchange  string `json:"exchange"`
	Timestamp int64  `json:"timestamp"`
}
