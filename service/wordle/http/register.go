package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, server *Server) {
	wordles := router.Group("/users/:user_id/wordles")
	{
		wordles.GET("", server.ListWordles)
		wordles.GET("/:id", server.GetWordle)
		wordles.POST("", server.CreateWordle)
		wordles.PATCH("/:id", server.UpdateWordle)
		wordles.DELETE("/:id", server.DeleteWordle)
	}
}
