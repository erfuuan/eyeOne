package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"eyeOne/internal/exchange"
	"eyeOne/models"
)

var validExchanges = map[string]exchange.ExchangeType{
	"binance": exchange.Binance,
	"kucoin":  exchange.KuCoin,
}

func ExchangeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := strings.ToLower(c.Param("exchange"))
		_, ok := validExchanges[raw]
		if !ok {
			c.JSON(http.StatusBadRequest, models.ErrorPayload{
				Message:    "Invalid exchange. Allowed: binance, kucoin",
				Details:    nil,
				Timestamp:  time.Now().Unix(),
				StatusCode: http.StatusBadRequest,
			})
			c.Abort()
			return
		}
		fmt.Println(raw)
		c.Set("exchange", raw)
		c.Next()
	}
}
