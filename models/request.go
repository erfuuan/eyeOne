package models

import (
	"errors"
	"fmt"
	"regexp"

	"go.uber.org/zap"

	"eyeOne/pkg/logger"
)

// ========== Order Requests ==========

type CreateOrderRequest struct {
	Symbol    string  `json:"symbol" binding:"required"`    // e.g. BTCUSDT
	Side      string  `json:"side" binding:"required"`      // buy / sell
	OrderType string  `json:"orderType" binding:"required"` // market / limit
	Quantity  float64 `json:"quantity" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
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
	Limit  int    `json:"limit" binding:"omitempty"`
}

// ========== Validation ==========

var (
	validSymbolRegex = regexp.MustCompile(`^[A-Z0-9]{6,20}$`)
	validAssetRegex  = regexp.MustCompile(`^[A-Z0-9]{2,10}$`)
)

func ValidateSymbol(symbol string) error {
	log := logger.GetLogger()

	if symbol == "" {
		err := errors.New("symbol is required")
		log.Warn("Validation error", zap.String("field", "symbol"), zap.Error(err))
		return err
	}
	if !validSymbolRegex.MatchString(symbol) {
		err := fmt.Errorf("symbol must be uppercase letters or digits, 6 to 20 characters (got: %s)", symbol)
		log.Warn("Validation error", zap.String("field", "symbol"), zap.Error(err))
		return err
	}
	return nil
}

func ValidateAsset(asset string) error {
	log := logger.GetLogger()

	if asset == "" {
		err := errors.New("asset is required")
		log.Warn("Validation error", zap.String("field", "asset"), zap.Error(err))
		return err
	}
	if !validAssetRegex.MatchString(asset) {
		err := fmt.Errorf("asset must be uppercase letters or digits, 2 to 10 characters (got: %s)", asset)
		log.Warn("Validation error", zap.String("field", "asset"), zap.Error(err))
		return err
	}
	return nil
}
