package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, server *Server) {
	sessions := router.Group("/users/:id/sessions")
	{
		sessions.GET("", server.GetSession)
		sessions.POST("", server.CreateSession)
	}

	users := router.Group("/users")
	{
		users.GET("/:id", server.GetUser)
		users.POST("", server.CreateUser)
		users.PATCH("/:id", server.UpdateUser)
		users.DELETE("/:id", server.DeleteUser)
	}
}
