package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"eyeOne/config"
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
		Addr:           ":" + cfg.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	exchanges := map[exchange.ExchangeType]exchange.Exchange{
		exchange.Binance: exchange.NewBinanceExchange(cfg.BinanceAPIKey, cfg.BinanceSecretKey),
		// exchange.KuCoin:  exchange.NewKuCoinExchange(cfg.KucoinAPIKey, cfg.KucoinSecretKey, cfg.KucoinPassphrase),
	}

	tradingService := service.NewTradingService(exchanges)

	h := handler.NewHandler(tradingService)

	api := router.Group("/api/v1")
	{
		api.POST("/order/:exchange", h.CreateOrder)
		api.DELETE("/order/:exchange/:orderID", h.CancelOrder)
		api.GET("/balance/:exchange/:asset", h.GetBalance)
		api.GET("/order-book/:exchange/:symbol", h.GetOrderBook)
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to run server: %v", err)
	}
}
