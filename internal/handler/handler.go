package handler

import (
	"eyeOne/internal/service"
)

// Handler manages HTTP requests.
type Handler struct {
	service *service.TradingService
}

// NewHandler creates a new Handler instance.
func NewHandler(svc *service.TradingService) *Handler {
	return &Handler{service: svc}
}

// Implement HTTP handler methods here...
