package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	internalhttp "github.com/trenddapp/backend/pkg/http"
	"github.com/trenddapp/backend/service/nft/client/etherscan"
	"github.com/trenddapp/backend/service/nft/client/nftport"
)

type Server struct {
	clientEtherscan etherscan.Client
	clientNFTPort   nftport.Client
}

func NewServer(clientEtherscan etherscan.Client, clientNFTPort nftport.Client) *Server {
	return &Server{
		clientEtherscan: clientEtherscan,
		clientNFTPort:   clientNFTPort,
	}
}

func (s *Server) ListAccountNFTs(ctx *gin.Context) {
	address := ctx.Param("address")
	if address == "" {
		internalhttp.NewError(http.StatusBadRequest, "invalid address").WriteJSON(ctx)
		return
	}

	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil {
		internalhttp.NewError(http.StatusBadRequest, "invalid page_size").WriteJSON(ctx)
		return
	}

	pageToken := ctx.Query("page_token")

	nfts, nextPageToken, err := s.clientEtherscan.ListAccountNFTs(ctx, address, pageSize, pageToken)
	if err != nil {
		internalhttp.NewError(http.StatusInternalServerError, "internal server error").WriteJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"nfts": nfts, "next_page_token": nextPageToken})
}

func (s *Server) ListContractNFTs(ctx *gin.Context) {
	address := ctx.Param("address")
	if address == "" {
		internalhttp.NewError(http.StatusBadRequest, "invalid address").WriteJSON(ctx)
		return
	}

	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil {
		internalhttp.NewError(http.StatusBadRequest, "invalid page_size")
		return
	}

	pageToken := ctx.Query("page_token")

	nfts, nextPageToken, err := s.clientNFTPort.ListContractNFTs(ctx, address, pageSize, pageToken)
	if err != nil {
		internalhttp.NewError(http.StatusInternalServerError, "internal server error").WriteJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"nfts": nfts, "next_page_token": nextPageToken})
}
