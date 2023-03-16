package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, server *Server) {
	router.GET("/accounts/:address/nfts", server.ListAccountNFTs)
	router.GET("/contracts/:address/nfts", server.ListContractNFTs)
}
