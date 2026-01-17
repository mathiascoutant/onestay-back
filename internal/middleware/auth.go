package middleware

import (
	"net/http"
	"strings"

	"onestay-back/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		
		// Récupérer le header Authorization (format standard: "Bearer <token>")
		authHeader := c.GetHeader("Authorization")
		
		if authHeader != "" {
			// Extraire le token du header "Bearer <token>"
			authHeader = strings.TrimSpace(authHeader)
			parts := strings.SplitN(authHeader, " ", 2)
			
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				tokenString = strings.TrimSpace(parts[1])
			}
		}
		
		// Si pas trouvé dans Authorization, vérifier le header "Bearer" directement (Bruno)
		if tokenString == "" {
			bearerHeader := c.GetHeader("Bearer")
			if bearerHeader != "" {
				tokenString = strings.TrimSpace(bearerHeader)
			}
		}
		
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token manquant",
			})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token invalide ou expiré",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		// Stocker les claims dans le contexte pour les utiliser dans les handlers
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role_id", claims.RoleID)

		c.Next()
	}
}
