package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("role_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Rôle non trouvé",
			})
			c.Abort()
			return
		}

		roleIDStr, ok := roleID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erreur lors de la vérification du rôle",
			})
			c.Abort()
			return
		}

		// Vérifier que le rôle est admin (ID: "3") ou superadmin (ID: "4")
		if roleIDStr != "3" && roleIDStr != "4" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Accès refusé. Seuls les administrateurs peuvent accéder à cette ressource",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("role_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Rôle non trouvé",
			})
			c.Abort()
			return
		}

		roleIDStr, ok := roleID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erreur lors de la vérification du rôle",
			})
			c.Abort()
			return
		}

		// Vérifier que le rôle est superadmin (ID: "4")
		if roleIDStr != "4" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Accès refusé. Seuls les super administrateurs peuvent accéder à cette ressource",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
