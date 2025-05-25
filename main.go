package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"eyeOne/config"
	"eyeOne/internal/api"
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

	exchanges := make(map[exchange.ExchangeType]exchange.Exchange)

	binance, err := exchange.NewBinanceExchange(cfg.BinanceAPIKey, cfg.BinanceSecretKey)
	if err != nil {
		log.Fatalf("Error initializing Binance: %v", err)
	}
	kucoin, err := exchange.NewKucoinExchange(cfg.KucoinAPIKey, cfg.KucoinSecretKey, cfg.KucoinPassphrase)
	if err != nil {
		log.Fatalf("Error initializing KuCoin: %v", err)
	}

	exchanges[exchange.Binance] = binance
	exchanges[exchange.KuCoin] = kucoin

	tradingService := service.NewTradingService(exchanges)

	h := handler.NewHandler(tradingService)

	api.SetupRouter(router, h)

	server := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to run server: %v", err)
	}
}
