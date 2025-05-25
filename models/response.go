package models

type SuccessResponse struct {
	Status    string `json:"status"` // always "success"
	Data      any    `json:"data,omitempty"`
	Message   string `json:"message,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

type ErrorResponse struct {
	Status    string       `json:"status"` // always "error"
	Error     ErrorPayload `json:"error"`
	Timestamp int64        `json:"timestamp"`
}

type ErrorPayload struct {
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	Timestamp  int64       `json:"timestamp"`
	StatusCode int         `json:"statusCode"`
}
