package router

import (
	"onestay-back/internal/handlers"
	"onestay-back/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	authHandler := handlers.NewAuthHandler()

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", middleware.AuthMiddleware(), middleware.RequireAdmin(), authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/roles", middleware.AuthMiddleware(), middleware.RequireAdmin(), authHandler.GetRoles)
			auth.POST("/roles", middleware.AuthMiddleware(), middleware.RequireSuperAdmin(), authHandler.CreateRole)
			auth.DELETE("/roles/:id", middleware.AuthMiddleware(), middleware.RequireSuperAdmin(), authHandler.DeleteRole)
		}
	}

	return r
}
