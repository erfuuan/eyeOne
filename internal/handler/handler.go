// package handler

// import (
// 	"net/http"
// 	"strings"
// 	"time"
//
//
//
//

// 	"github.com/gin-gonic/gin"

// 	"eyeOne/internal/exchange"
// 	"eyeOne/internal/service"
// 	"eyeOne/models"
// )

// type Handler struct {
// 	service *service.TradingService
// }

// func NewHandler(s *service.TradingService) *Handler {
// 	return &Handler{service: s}
// }

// func (h *Handler) CreateOrder(c *gin.Context) {
// 	var req models.CreateOrderRequest

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, models.ErrorPayload{
// 			StatusCode: http.StatusBadRequest,
// 			Message:    "Invalid request payload",
// 			Timestamp:  time.Now().Unix(),
// 		})
// 		return
// 	}

// 	exchangeName, _ := c.Get("exchange")
// 	exNameStr := exchangeName.(string)
// 	exName := exchange.ExchangeType(exNameStr)
// 	orderID, err := h.service.CreateOrder(
// 		c.Request.Context(),
// 		exName,
// 		strings.ToUpper(req.Symbol),
// 		strings.ToLower(req.Side),
// 		req.OrderType,
// 		req.Quantity,
// 		req.Price,
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
// 			StatusCode: http.StatusInternalServerError,
// 			Message:    err.Error(),
// 			Timestamp:  time.Now().Unix(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, models.SuccessResponse{
// 		StatusCode: http.StatusCreated,
// 		Data: models.OrderDataResponse{
// 			OrderID:  orderID,
// 			Exchange: exNameStr,
// 			Symbol:   strings.ToUpper(req.Symbol),
// 			Side:     strings.ToLower(req.Side),
// 			Type:     req.OrderType,
// 			Quantity: req.Quantity,
// 			Price:    req.Price,
// 		},
// 		Message:   "order created successfully",
// 		Timestamp: time.Now().Unix(),
// 	})
// }

// func (h *Handler) CancelOrder(c *gin.Context) {
// 	exchangeName, _ := c.Get("exchange")
// 	exNameStr := exchangeName.(string)
// 	exName := exchange.ExchangeType(exNameStr)
// 	symbol := c.Param("symbol")
// 	orderID := c.Param("orderID")

// 	if orderID == "" {
// 		c.JSON(http.StatusBadRequest, models.ErrorPayload{
// 			StatusCode: http.StatusBadRequest,
// 			Message:    "missing orderID parameter",
// 			Timestamp:  time.Now().Unix(),
// 		})
// 		return
// 	}

// 	if err := models.ValidateSymbol(symbol); err != nil {
// 		c.JSON(http.StatusBadRequest, models.ErrorPayload{
// 			StatusCode: http.StatusBadRequest,
// 			Message:    err.Error(),
// 			Timestamp:  time.Now().Unix(),
// 		})
// 		return
// 	}

// 	err := h.service.CancelOrder(c.Request.Context(), exName, symbol, orderID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
// 			StatusCode: http.StatusInternalServerError,
// 			Message:    err.Error(),
// 			Timestamp:  time.Now().Unix(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, models.SuccessResponse{
// 		StatusCode: http.StatusOK,
// 		Message:    "order canceled",
// 		Timestamp:  time.Now().Unix(),
// 	})
// }

// func (h *Handler) GetBalance(c *gin.Context) {
// 	exchangeName, _ := c.Get("exchange")
// 	exNameStr := exchangeName.(string)
// 	exName := exchange.ExchangeType(exNameStr)
// 	asset := c.Param("asset")

// 	if err := models.ValidateAsset(asset); err != nil {
// 		c.JSON(http.StatusBadRequest, models.ErrorPayload{
// 			StatusCode: http.StatusBadRequest,
// 			Message:    err.Error(),
// 			Timestamp:  time.Now().Unix(),
// 		})
// 		return
// 	}

// 	balance, err := h.service.GetBalance(c.Request.Context(), exName, asset)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
// 			StatusCode: http.StatusInternalServerError,
// 			Message:    err.Error(),
// 			Timestamp:  time.Now().Unix(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, models.SuccessResponse{
// 		StatusCode: http.StatusOK,
// 		Data: models.BalanceDataResponse{
// 			Asset:   asset,
// 			Balance: balance,
// 		},
// 		Message:   "balance retrieved successfully",
// 		Timestamp: time.Now().Unix(),
// 	})

// }

// func (h *Handler) GetOrderBook(c *gin.Context) {
// 	exchangeName, _ := c.Get("exchange")
// 	exNameStr := exchangeName.(string)
// 	exName := exchange.ExchangeType(exNameStr)

// 	symbol := c.Param("symbol")
// 	// if ok, errMsg := models.ValidateSymbol(symbol); !ok {
// 	if err := models.ValidateSymbol(symbol); err != nil {
// 		c.JSON(http.StatusBadRequest, models.ErrorPayload{
// 			StatusCode: http.StatusBadRequest,
// 			Message:    err.Error(),
// 			Timestamp:  time.Now().Unix(),
// 		})
// 		return
// 	}

// 	orderBook, err := h.service.GetOrderBook(c.Request.Context(), exName, symbol)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
// 			StatusCode: http.StatusInternalServerError,
// 			Message:    err.Error(),
// 			Timestamp:  time.Now().Unix(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, models.SuccessResponse{
// 		StatusCode: http.StatusOK,
// 		Data:       orderBook,
// 		Message:    "order book retrieved successfully",
// 		Timestamp:  time.Now().Unix(),
//

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

// CreateOrder godoc
// @Summary Create new order
// @Description Place a new order on the exchange
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body models.CreateOrderRequest true "Order Payload"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorPayload
// @Failure 500 {object} models.ErrorResponse
// @Router /orders [post]
// @Security ApiKeyAuth
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

// CancelOrder godoc
// @Summary Cancel an order
// @Description Cancel an existing order on the exchange
// @Tags Orders
// @Accept json
// @Produce json
// @Param symbol path string true "Symbol"
// @Param orderID path string true "Order ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorPayload
// @Failure 500 {object} models.ErrorResponse
// @Router /orders/{symbol}/{orderID} [delete]
// @Security ApiKeyAuth
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

	if err := models.ValidateSymbol(symbol); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorPayload{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
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

// GetBalance godoc
// @Summary Get balance
// @Description Get balance for a specific asset on the exchange
// @Tags Balance
// @Accept json
// @Produce json
// @Param asset path string true "Asset symbol"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorPayload
// @Failure 500 {object} models.ErrorResponse
// @Router /balance/{asset} [get]
// @Security ApiKeyAuth
func (h *Handler) GetBalance(c *gin.Context) {
	exchangeName, _ := c.Get("exchange")
	exNameStr := exchangeName.(string)
	exName := exchange.ExchangeType(exNameStr)
	asset := c.Param("asset")

	if err := models.ValidateAsset(asset); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorPayload{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
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

// GetOrderBook godoc
// @Summary Get order book
// @Description Retrieve the order book for a specific trading pair
// @Tags Market
// @Accept json
// @Produce json
// @Param symbol path string true "Trading Symbol (e.g., BTCUSDT)"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorPayload
// @Failure 500 {object} models.ErrorResponse
// @Router /orderbook/{symbol} [get]
// @Security ApiKeyAuth
func (h *Handler) GetOrderBook(c *gin.Context) {
	exchangeName, _ := c.Get("exchange")
	exNameStr := exchangeName.(string)
	exName := exchange.ExchangeType(exNameStr)

	symbol := c.Param("symbol")

	if err := models.ValidateSymbol(symbol); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorPayload{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
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
