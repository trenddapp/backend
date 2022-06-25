package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	internalhttp "github.com/trenddapp/backend/pkg/http"
	"github.com/trenddapp/backend/service/currency/client/coinmarketcap"
)

type Server struct {
	clientCoinMarketCap coinmarketcap.Client
}

func NewServer(clientCoinMarketCap coinmarketcap.Client) *Server {
	return &Server{
		clientCoinMarketCap: clientCoinMarketCap,
	}
}

func (s *Server) GetRate(ctx *gin.Context) {
	symbol := ctx.Param("symbol")
	if symbol == "" {
		internalhttp.NewError(http.StatusBadRequest, "invalid symbol").WriteJSON(ctx)
		return
	}

	rate, err := s.clientCoinMarketCap.GetRate(ctx, symbol)
	if err != nil {
		if err == coinmarketcap.ErrInvalidSymbol {
			internalhttp.NewError(http.StatusBadRequest, "invalid symbol").WriteJSON(ctx)
			return
		}

		internalhttp.NewError(http.StatusInternalServerError, "internal server error").WriteJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"rate": rate})
}
