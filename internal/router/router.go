package router

import (
	"onestay-back/internal/handlers"
	"onestay-back/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Configuration CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
	}))

	authHandler := handlers.NewAuthHandler()
	propertyHandler := handlers.NewPropertyHandler()

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
			users.GET("/profile", middleware.AuthMiddleware(), authHandler.GetProfile)
			users.PUT("/profile", middleware.AuthMiddleware(), authHandler.UpdateProfile)
			users.DELETE("/profile", middleware.AuthMiddleware(), authHandler.DeleteAccount)
			users.GET("", middleware.AuthMiddleware(), middleware.RequireAdmin(), authHandler.GetAllUsers)
			users.PUT("/:id", middleware.AuthMiddleware(), middleware.RequireAdmin(), authHandler.UpdateUser)
			users.DELETE("/:id", middleware.AuthMiddleware(), middleware.RequireAdmin(), authHandler.DeleteUser)
		}

		properties := api.Group("/properties")
		{
			properties.POST("", middleware.AuthMiddleware(), propertyHandler.CreateProperty)
			properties.GET("/user/:id", propertyHandler.GetUserProperties)
			properties.GET("/:id", propertyHandler.GetProperty)
			properties.PUT("/:id", middleware.AuthMiddleware(), propertyHandler.UpdateProperty)
			properties.POST("/:id/publish", middleware.AuthMiddleware(), propertyHandler.PublishProperty)
			properties.DELETE("/:id", middleware.AuthMiddleware(), propertyHandler.DeleteProperty)
		}
	}

	return r
}
