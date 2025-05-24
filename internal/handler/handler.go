package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"eyeOne/internal/exchange"
)

// Handler struct holds the exchange implementation.
type Handler struct {
	exchange exchange.Exchange
}

// NewHandler creates a new Handler instance with the provided exchange.
func NewHandler(ex exchange.Exchange) *Handler {
	return &Handler{exchange: ex}
}

// CreateOrder handles the creation of a new order.
func (h *Handler) CreateOrder(c *gin.Context) {
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

	orderID, err := h.exchange.CreateOrder(c.Request.Context(), req.Symbol, req.Side, req.OrderType, req.Quantity, req.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order_id": orderID})
}

// CancelOrder handles the cancellation of an existing order.
func (h *Handler) CancelOrder(c *gin.Context) {
	symbol := c.Query("symbol")
	orderID := c.Query("order_id")

	if symbol == "" || orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing symbol or order_id"})
		return
	}

	err := h.exchange.CancelOrder(c.Request.Context(), symbol, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "order canceled"})
}

// GetBalance retrieves the balance for a specific asset.
func (h *Handler) GetBalance(c *gin.Context) {
	asset := c.Query("asset")
	if asset == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing asset parameter"})
		return
	}

	balance, err := h.exchange.GetBalance(c.Request.Context(), asset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"asset": asset, "balance": balance})
}

// GetOrderBook retrieves the order book for a specific symbol.
func (h *Handler) GetOrderBook(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing symbol parameter"})
		return
	}

	orderBook, err := h.exchange.GetOrderBook(c.Request.Context(), symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orderBook)
}
