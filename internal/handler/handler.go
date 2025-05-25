package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"eyeOne/internal/exchange"
	"eyeOne/internal/service"
)

type Handler struct {
	service *service.TradingService
}

func NewHandler(s *service.TradingService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) parseExchange(c *gin.Context) (exchange.ExchangeType, bool) {
	exName := exchange.ExchangeType(c.Param("exchange"))
	if exName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing exchange parameter"})
		return "", false
	}
	return exName, true
}

func (h *Handler) CreateOrder(c *gin.Context) {
	exName, ok := h.parseExchange(c)
	if !ok {
		return
	}

	var req struct {
		Symbol    string  `json:"symbol"`
		Side      string  `json:"side"`
		OrderType string  `json:"order_type"`
		Quantity  float64 `json:"quantity"`
		Price     float64 `json:"price"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderID, err := h.service.CreateOrder(c.Request.Context(), exName, req.Symbol, req.Side, req.OrderType, req.Quantity, req.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order_id": orderID})
}

func (h *Handler) CancelOrder(c *gin.Context) {
	exName, ok := h.parseExchange(c)
	if !ok {
		return
	}

	symbol := c.Query("symbol")
	orderID := c.Param("orderID")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing orderID parameter"})
		return
	}

	err := h.service.CancelOrder(c.Request.Context(), exName, symbol, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "order canceled"})
}

func (h *Handler) GetBalance(c *gin.Context) {
	exName, ok := h.parseExchange(c)
	if !ok {
		return
	}

	asset := c.Param("asset")
	if asset == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing asset parameter"})
		return
	}

	balance, err := h.service.GetBalance(c.Request.Context(), exName, asset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"asset": asset, "balance": balance})
}

func (h *Handler) GetOrderBook(c *gin.Context) {
	exName, ok := h.parseExchange(c)
	if !ok {
		return
	}

	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing symbol parameter"})
		return
	}

	orderBook, err := h.service.GetOrderBook(c.Request.Context(), exName, symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orderBook)
}
