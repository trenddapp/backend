package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dapp-z/backend/service/currency/client/coinmarketcap"
)

type Server struct {
	clientCoinMarketCap coinmarketcap.Client
}

func NewServer(clientCoinMarketCap coinmarketcap.Client) *Server {
	return &Server{
		clientCoinMarketCap: clientCoinMarketCap,
	}
}

func (s *Server) GetConversionRate(c *gin.Context) {
	symbol := c.Param("symbol")

	conversionRate, err := s.clientCoinMarketCap.GetConversionRate(c, symbol)
	if err != nil {
		if err == coinmarketcap.ErrInvalidSymbol {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid symbol",
			})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"conversion_rate": conversionRate,
	})
}
