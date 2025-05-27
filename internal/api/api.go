package api

import (
	"github.com/gin-gonic/gin"

	"eyeOne/internal/handler"
	"eyeOne/internal/middleware"
)

func SetupRouter(router *gin.Engine, h *handler.Handler) {
	api := router.Group("/api/v1")
	api.POST("/order/:exchange", middleware.ExchangeMiddleware(), h.CreateOrder)
	api.DELETE("/order/:exchange/:orderID", middleware.ExchangeMiddleware(), h.CancelOrder)
	api.GET("/balance/:exchange/:asset", middleware.ExchangeMiddleware(), h.GetBalance)
	api.GET("/order-book/:exchange/:symbol", middleware.ExchangeMiddleware(), h.GetOrderBook)
}
