package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"eyeOne/internal/exchange"
	"eyeOne/models"
	"eyeOne/pkg/logger"
)

var validExchanges = map[string]exchange.ExchangeType{
	"binance": exchange.Binance,
	"kucoin":  exchange.KuCoin,
}

func ExchangeMiddleware() gin.HandlerFunc {
	log := logger.GetLogger()

	return func(c *gin.Context) {
		raw := strings.ToLower(c.Param("exchange"))
		_, ok := validExchanges[raw]
		if !ok {
			log.Warn("Invalid exchange parameter",
				zap.String("exchange", raw),
				zap.String("path", c.Request.URL.Path),
			)

			c.JSON(http.StatusBadRequest, models.ErrorPayload{
				Message:    "Invalid exchange. Allowed: binance, kucoin",
				Timestamp:  time.Now().Unix(),
				StatusCode: http.StatusBadRequest,
			})
			c.Abort()
			return
		}

		log.Info("Exchange validated",
			zap.String("exchange", raw),
			zap.String("path", c.Request.URL.Path),
		)

		c.Set("exchange", raw)
		c.Next()
	}
}
