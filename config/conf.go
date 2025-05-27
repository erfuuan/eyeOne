package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"eyeOne/pkg/logger"
)

type Config struct {
	Port             string
	BinanceAPIKey    string
	BinanceSecretKey string
	KucoinAPIKey     string
	KucoinSecretKey  string
	KucoinPassphrase string
	BitpinAPIKey     string
	BitpinSecretKey  string
}

func LoadEnv() *Config {
	logger := logger.GetLogger()

	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found. Using system environment variables")
	}

	cfg := &Config{
		Port:             getEnv("PORT", "8080"),
		BinanceAPIKey:    mustGetEnv("BINANCE_API_KEY", logger),
		BinanceSecretKey: mustGetEnv("BINANCE_SECRET_KEY", logger),
		KucoinAPIKey:     mustGetEnv("KUCOIN_API_KEY", logger),
		KucoinSecretKey:  mustGetEnv("KUCOIN_SECRET_KEY", logger),
		KucoinPassphrase: mustGetEnv("KUCOIN_PASSPHRASE", logger),
		BitpinAPIKey:     mustGetEnv("BITPIN_API_KEY", logger),
		BitpinSecretKey:  mustGetEnv("BITPIN_SECRET_KEY", logger),
	}

	return cfg
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func mustGetEnv(key string, logger *zap.Logger) string {
	val := os.Getenv(key)
	if val == "" {
		logger.Fatal("Required environment variable not set", zap.String("key", key))
	}
	return val
}
