package handlers

import (
	"net/http"
	"strings"

	"onestay-back/internal/models"
	"onestay-back/internal/repository"
	"onestay-back/internal/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogementHandler struct {
	logementRepo *repository.LogementRepository
}

func NewLogementHandler() *LogementHandler {
	return &LogementHandler{
		logementRepo: repository.NewLogementRepository(),
	}
}

func (h *LogementHandler) CreateLogement(c *gin.Context) {
	var req models.CreateLogementRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Données invalides",
			"details": err.Error(),
		})
		return
	}

	// Récupérer l'ID de l'utilisateur depuis le contexte (défini par le middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Utilisateur non authentifié",
		})
		return
	}

	userID, ok := userIDInterface.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la récupération de l'ID utilisateur",
		})
		return
	}

	ctx := c.Request.Context()

	// Vérifier que le nom du bien n'est pas déjà utilisé par cet utilisateur
	exists, err := h.logementRepo.ExistsByNomBienAndUserID(ctx, req.NomBien, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la vérification du nom du bien",
		})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Un logement avec ce nom existe déjà pour votre compte",
		})
		return
	}

	// Créer le logement
	logement := &models.Logement{
		NomBien:     req.NomBien,
		Description: req.Description,
		Adresse:     req.Adresse,
		Ville:       req.Ville,
		Pays:        req.Pays,
		UserID:      userID,
		Status:      1, // 1 = brouillon, 2 = publié
	}

	if err := h.logementRepo.Create(ctx, logement); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la création du logement",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Logement créé avec succès",
		"logement": gin.H{
			"id":          logement.ID,
			"nom_bien":    logement.NomBien,
			"description": logement.Description,
			"adresse":     logement.Adresse,
			"ville":       logement.Ville,
			"pays":        logement.Pays,
			"user_id":     logement.UserID,
			"status":      logement.Status,
			"created_at":  logement.CreatedAt,
		},
	})
}

func (h *LogementHandler) GetUserLogements(c *gin.Context) {
	// Récupérer l'ID de l'utilisateur depuis l'URL
	userIDParam := c.Param("id")
	if userIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID utilisateur manquant",
		})
		return
	}

	// Convertir l'ID de l'URL en ObjectID
	requestedUserID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID utilisateur invalide",
		})
		return
	}

	ctx := c.Request.Context()

	// Essayer de récupérer l'ID de l'utilisateur depuis le token JWT (optionnel)
	var tokenUserID primitive.ObjectID
	var isOwner bool

	// Vérifier si un token est présent dans le contexte (défini par le middleware)
	userIDInterface, exists := c.Get("user_id")
	if exists {
		if id, ok := userIDInterface.(primitive.ObjectID); ok {
			tokenUserID = id
			// Vérifier si c'est le propriétaire
			isOwner = requestedUserID == tokenUserID
		}
	} else {
		// Si pas dans le contexte, essayer d'extraire le token manuellement
		var tokenString string
		
		// Récupérer le header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			authHeader = strings.TrimSpace(authHeader)
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				tokenString = strings.TrimSpace(parts[1])
			}
		}
		
		// Si pas trouvé dans Authorization, vérifier le header "Bearer" directement
		if tokenString == "" {
			bearerHeader := c.GetHeader("Bearer")
			if bearerHeader != "" {
				tokenString = strings.TrimSpace(bearerHeader)
			}
		}
		
		// Si un token est présent, essayer de le valider
		if tokenString != "" {
			claims, err := utils.ValidateToken(tokenString)
			if err == nil {
				tokenUserID = claims.UserID
				isOwner = requestedUserID == tokenUserID
			}
		}
	}

	// Déterminer si on doit inclure les brouillons (seulement si c'est le propriétaire)
	includeBrouillon := isOwner

	// Récupérer les logements
	logements, err := h.logementRepo.FindByUserID(ctx, requestedUserID, includeBrouillon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la récupération des logements",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logements": logements,
		"count":     len(logements),
	})
}
