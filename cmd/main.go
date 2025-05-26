package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"eyeOne/config"
	"eyeOne/internal/api"
	"eyeOne/internal/exchange"
	"eyeOne/internal/handler"
	"eyeOne/internal/service"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found, relying on environment variables")
	}

	cfg := config.LoadEnv()
	router := gin.Default()

	// Setup Exchanges
	exchanges := make(map[exchange.ExchangeType]exchange.Exchange)

	binance, err := exchange.NewBinanceExchange(cfg.BinanceAPIKey, cfg.BinanceSecretKey)
	if err != nil {
		logger.Fatal("Failed to initialize Binance", zap.Error(err))
	}
	kucoin, err := exchange.NewKucoinExchange(cfg.KucoinAPIKey, cfg.KucoinSecretKey, cfg.KucoinPassphrase)
	if err != nil {
		logger.Fatal("Failed to initialize KuCoin", zap.Error(err))
	}

	exchanges[exchange.Binance] = binance
	exchanges[exchange.KuCoin] = kucoin

	// Setup service and handlers
	tradingService := service.NewTradingService(exchanges)
	h := handler.NewHandler(tradingService)

	// Register routes
	api.SetupRouter(router, h)

	// Create HTTP server
	server := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Graceful shutdown handling
	go func() {
		logger.Info("Starting server", zap.String("port", cfg.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server startup failed", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	} else {
		logger.Info("Server exited gracefully")
	}
}
