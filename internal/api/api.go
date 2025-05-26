package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"eyeOne/internal/handler"
	"eyeOne/internal/middleware"
)

func SetupRouter(router *gin.Engine, h *handler.Handler) {
	api := router.Group("/api/v1")

	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api.POST("/order/:exchange", middleware.ExchangeMiddleware(), h.CreateOrder)
	api.DELETE("/order/:exchange/:symbol/:orderID", middleware.ExchangeMiddleware(), h.CancelOrder)
	api.GET("/balance/:exchange/:asset", middleware.ExchangeMiddleware(), h.GetBalance)
	api.GET("/order-book/:exchange/:symbol", middleware.ExchangeMiddleware(), h.GetOrderBook)
}
