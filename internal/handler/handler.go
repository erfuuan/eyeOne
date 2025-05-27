package handler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"eyeOne/internal/exchange"
	"eyeOne/internal/service"
	"eyeOne/models"
)

type Handler struct {
	service *service.TradingService
}

func NewHandler(s *service.TradingService) *Handler {
	return &Handler{service: s}
}

func getExchange(c *gin.Context) (exchange.ExchangeType, string, bool) {
	exchangeName, exists := c.Get("exchange")
	if !exists {
		return "", "", false
	}

	exNameStr, ok := exchangeName.(string)
	if !ok {
		return "", "", false
	}

	return exchange.ExchangeType(exNameStr), exNameStr, true
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorPayload{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
			Timestamp:  time.Now().Unix(),
		})
		return
	}

	exName, exNameStr, ok := getExchange(c)
	if !ok {
		c.JSON(http.StatusBadRequest, models.ErrorPayload{
			StatusCode: http.StatusBadRequest,
			Message:    "Missing or invalid exchange name",
			Timestamp:  time.Now().Unix(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	orderID, err, status := h.service.CreateOrder(
		ctx,
		exName,
		strings.ToUpper(req.Symbol),
		strings.ToLower(req.Side),
		req.OrderType,
		req.Quantity,
		req.Price,
	)
	if err != nil {
		c.JSON(status, models.ErrorResponse{
			StatusCode: status,
			Message:    err.Error(),
			Timestamp:  time.Now().Unix(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		StatusCode: http.StatusCreated,
		Data: models.OrderDataResponse{
			OrderID:  orderID,
			Exchange: exNameStr,
			Symbol:   strings.ToUpper(req.Symbol),
			Side:     strings.ToLower(req.Side),
			Type:     req.OrderType,
			Quantity: req.Quantity,
			Price:    req.Price,
		},
		Message:   "order created successfully",
		Timestamp: time.Now().Unix(),
	})
}

func (h *Handler) CancelOrder(c *gin.Context) {
	exName, _, ok := getExchange(c)
	if !ok {
		c.JSON(http.StatusBadRequest, models.ErrorPayload{
			StatusCode: http.StatusBadRequest,
			Message:    "Missing or invalid exchange name",
			Timestamp:  time.Now().Unix(),
		})
		return
	}

	symbol := c.Param("symbol")
	orderID := c.Param("orderID")

	if orderID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorPayload{
			StatusCode: http.StatusBadRequest,
			Message:    "Missing orderID parameter",
			Timestamp:  time.Now().Unix(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err, status := h.service.CancelOrder(ctx, exName, symbol, orderID)
	if err != nil {
		c.JSON(status, models.ErrorResponse{
			StatusCode: status,
			Message:    err.Error(),
			Timestamp:  time.Now().Unix(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		StatusCode: http.StatusOK,
		Message:    "order canceled",
		Timestamp:  time.Now().Unix(),
	})
}

func (h *Handler) GetBalance(c *gin.Context) {
	exName, _, ok := getExchange(c)
	if !ok {
		c.JSON(http.StatusBadRequest, models.ErrorPayload{
			StatusCode: http.StatusBadRequest,
			Message:    "Missing or invalid exchange name",
			Timestamp:  time.Now().Unix(),
		})
		return
	}

	asset := c.Param("asset")
	if err := models.ValidateAsset(asset); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorPayload{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Timestamp:  time.Now().Unix(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	balance, err, status := h.service.GetBalance(ctx, exName, asset)
	if err != nil {
		c.JSON(status, models.ErrorResponse{
			StatusCode: status,
			Message:    err.Error(),
			Timestamp:  time.Now().Unix(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		StatusCode: http.StatusOK,
		Data: models.BalanceDataResponse{
			Asset:   asset,
			Balance: balance,
		},
		Message:   "balance retrieved successfully",
		Timestamp: time.Now().Unix(),
	})
}

func (h *Handler) GetOrderBook(c *gin.Context) {
	exName, _, ok := getExchange(c)
	if !ok {
		c.JSON(http.StatusBadRequest, models.ErrorPayload{
			StatusCode: http.StatusBadRequest,
			Message:    "Missing or invalid exchange name",
			Timestamp:  time.Now().Unix(),
		})
		return
	}

	symbol := c.Param("symbol")
	if err := models.ValidateSymbol(symbol, string(exName)); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorPayload{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Timestamp:  time.Now().Unix(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	orderBook, err, status := h.service.GetOrderBook(ctx, exName, symbol)
	if err != nil {
		c.JSON(status, models.ErrorResponse{
			StatusCode: status,
			Message:    err.Error(),
			Timestamp:  time.Now().Unix(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		StatusCode: http.StatusOK,
		Data:       orderBook,
		Message:    "order book retrieved successfully",
		Timestamp:  time.Now().Unix(),
	})
}
