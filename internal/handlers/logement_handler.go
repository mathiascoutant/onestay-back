package handlers

import (
	"net/http"

	"onestay-back/internal/models"
	"onestay-back/internal/repository"

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

	// Récupérer l'ID de l'utilisateur depuis le token JWT
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Utilisateur non authentifié",
		})
		return
	}

	tokenUserID, ok := userIDInterface.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la récupération de l'ID utilisateur",
		})
		return
	}

	ctx := c.Request.Context()

	// Déterminer si on doit inclure les brouillons (seulement si c'est le propriétaire)
	includeBrouillon := requestedUserID == tokenUserID

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
