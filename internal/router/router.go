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
			auth.POST("/login", authHandler.Login)
			auth.GET("/roles", middleware.AuthMiddleware(), middleware.RequireAdmin(), authHandler.GetRoles)
			auth.POST("/roles", middleware.AuthMiddleware(), middleware.RequireSuperAdmin(), authHandler.CreateRole)
			auth.DELETE("/roles/:id", middleware.AuthMiddleware(), middleware.RequireSuperAdmin(), authHandler.DeleteRole)
		}

		users := api.Group("/users")
		{
			users.POST("/register", middleware.AuthMiddleware(), middleware.RequireAdmin(), authHandler.Register)
			users.GET("", middleware.AuthMiddleware(), middleware.RequireAdmin(), authHandler.GetAllUsers)
			users.PUT("/:id", middleware.AuthMiddleware(), middleware.RequireAdmin(), authHandler.UpdateUser)
			users.DELETE("/:id", middleware.AuthMiddleware(), middleware.RequireAdmin(), authHandler.DeleteUser)
		}
	}

	return r
}
