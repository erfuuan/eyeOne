package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	BinanceAPIKey    string
	BinanceSecretKey string
	KucoinAPIKey     string
	KucoinSecretKey  string
	KucoinPassphrase string
}

func LoadEnv() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or could not load it, continuing with environment variables")
	}

	cfg := &Config{
		Port:             os.Getenv("PORT"),
		BinanceAPIKey:    os.Getenv("BINANCE_API_KEY"),
		BinanceSecretKey: os.Getenv("BINANCE_SECRET_KEY"),
		KucoinAPIKey:     os.Getenv("KUCOIN_API_KEY"),
		KucoinSecretKey:  os.Getenv("KUCOIN_SECRET_KEY"),
		KucoinPassphrase: os.Getenv("KUCOIN_PASSPHRASE"),
	}

	return cfg
}
