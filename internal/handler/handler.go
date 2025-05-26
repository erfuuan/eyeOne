package handler

import (
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

	exchangeName, _ := c.Get("exchange")
	exNameStr := exchangeName.(string)
	exName := exchange.ExchangeType(exNameStr)
	orderID, err := h.service.CreateOrder(
		c.Request.Context(),
		exName,
		strings.ToUpper(req.Symbol),
		strings.ToLower(req.Side),
		req.OrderType,
		req.Quantity,
		req.Price,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Timestamp:  time.Now().Unix(),
		})
		return
	}

	type OrderData struct {
		OrderID  string  `json:"orderId"`
		Exchange string  `json:"exchange"`
		Symbol   string  `json:"symbol"`
		Side     string  `json:"side"`
		Type     string  `json:"type"`
		Quantity float64 `json:"quantity"`
		Price    float64 `json:"price"`
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		StatusCode: http.StatusCreated,
		Data: OrderData{
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
	exchangeName, _ := c.Get("exchange")
	exNameStr := exchangeName.(string)
	exName := exchange.ExchangeType(exNameStr)
	symbol := c.Param("symbol")
	orderID := c.Param("orderID")

	if orderID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorPayload{
			StatusCode: http.StatusBadRequest,
			Message:    "missing orderID parameter",
			Timestamp:  time.Now().Unix(),
		})
		return
	}

	if ok, errMsg := models.ValidateSymbol(symbol); !ok {
		c.JSON(http.StatusBadRequest, models.ErrorPayload{
			StatusCode: http.StatusBadRequest,
			Message:    errMsg,
			Timestamp:  time.Now().Unix(),
		})
		return
	}

	err := h.service.CancelOrder(c.Request.Context(), exName, symbol, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
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
	exchangeName, _ := c.Get("exchange")
	exNameStr := exchangeName.(string)
	exName := exchange.ExchangeType(exNameStr)
	asset := c.Param("asset")

	if ok, errMsg := models.ValidateAsset(asset); !ok {
		c.JSON(http.StatusBadRequest, models.ErrorPayload{
			StatusCode: http.StatusBadRequest,
			Message:    errMsg,
			Timestamp:  time.Now().Unix(),
		})
		return
	}

	balance, err := h.service.GetBalance(c.Request.Context(), exName, asset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Timestamp:  time.Now().Unix(),
		})
		return
	}

	type BalanceData struct {
		Asset   string  `json:"asset"`
		Balance float64 `json:"balance"`
	}
	c.JSON(http.StatusOK, models.SuccessResponse{
		StatusCode: http.StatusOK,
		Data: BalanceData{
			Asset:   asset,
			Balance: balance,
		},
		Message:   "balance retrieved successfully",
		Timestamp: time.Now().Unix(),
	})

}

func (h *Handler) GetOrderBook(c *gin.Context) {
	exchangeName, _ := c.Get("exchange")
	exNameStr := exchangeName.(string)
	exName := exchange.ExchangeType(exNameStr)

	symbol := c.Param("symbol")
	if ok, errMsg := models.ValidateSymbol(symbol); !ok {
		c.JSON(http.StatusBadRequest, models.ErrorPayload{
			StatusCode: http.StatusBadRequest,
			Message:    errMsg,
			Timestamp:  time.Now().Unix(),
		})
		return
	}

	orderBook, err := h.service.GetOrderBook(c.Request.Context(), exName, symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
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
