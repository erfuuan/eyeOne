package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"eyeOne/internal/config"
	"eyeOne/internal/exchange"
	"eyeOne/internal/handler"
	"eyeOne/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	cfg := config.LoadEnv()

	router := gin.Default()

	server := &http.Server{
		Addr:           ":3000",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	binanceExchange := exchange.NewBinanceExchange(cfg.BinanceAPIKey, cfg.BinanceSecretKey)
	// kucoinExchange := exchange.NewKuCoinExchange(cfg.KucoinAPIKey, cfg.KucoinSecretKey, cfg.KucoinPassphrase)

	// tradingService := service.NewTradingService(binanceExchange, kucoinExchange)
	tradingService := service.NewTradingService(binanceExchange)

	h := handler.NewHandler(tradingService)

	api := router.Group("/api/v1")
	{
		api.POST("/order", h.CreateOrder)
		api.DELETE("/order-book/:id", h.CancelOrder)
		api.GET("/balance/:asset", h.GetBalance)
		api.GET("/order-book/:symbol", h.GetOrderBook)
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to run server: %v", err)
	}
}
