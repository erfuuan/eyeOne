package handler

import (
	"fmt"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
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
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":    "success",
		"orderId":   orderID,
		"exchange":  exName,
		"symbol":    strings.ToUpper(req.Symbol),
		"side":      strings.ToLower(req.Side),
		"type":      req.OrderType,
		"quantity":  req.Quantity,
		"price":     req.Price,
		"timestamp": time.Now().Unix(),
	})
}

func (h *Handler) CancelOrder(c *gin.Context) {
	exchangeName, _ := c.Get("exchange")
	exNameStr := exchangeName.(string)
	exName := exchange.ExchangeType(exNameStr)
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
	exchangeName, _ := c.Get("exchange")
	exNameStr := exchangeName.(string)
	exName := exchange.ExchangeType(exNameStr)
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
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing symbol parameter"})
		return
	}
	orderBook, err := h.service.GetOrderBook(c.Request.Context(), "binanace", symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orderBook)
}
