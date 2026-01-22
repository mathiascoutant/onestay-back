package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"onestay-back/internal/models"
	"onestay-back/internal/repository"
	"onestay-back/internal/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PropertyHandler struct {
	propertyRepo *repository.PropertyRepository
}

func NewPropertyHandler() *PropertyHandler {
	return &PropertyHandler{
		propertyRepo: repository.NewPropertyRepository(),
	}
}

// CreateProperty crée une nouvelle propriété
func (h *PropertyHandler) CreateProperty(c *gin.Context) {
	var req models.CreatePropertyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Données invalides",
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

	hostID, ok := userIDInterface.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la récupération de l'ID utilisateur",
		})
		return
	}

	ctx := c.Request.Context()

	// Vérifier que le nom n'est pas déjà utilisé par cet hôte
	nameExists, err := h.propertyRepo.ExistsByNameAndHostID(ctx, req.Name, hostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la vérification du nom",
		})
		return
	}
	if nameExists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Vous avez déjà un logement avec ce nom",
		})
		return
	}

	// Générer le slug à partir du nom
	baseSlug := utils.GenerateSlug(req.Name)
	slug := baseSlug

	// Vérifier l'unicité du slug et ajouter un suffixe si nécessaire
	counter := 1
	for {
		exists, err := h.propertyRepo.ExistsBySlug(ctx, slug)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erreur lors de la vérification du slug",
			})
			return
		}
		if !exists {
			break
		}
		slug = baseSlug + "-" + strconv.Itoa(counter)
		counter++
	}

	// Créer les sous-documents par défaut si non fournis
	property := &models.Property{
		HostID:      hostID,
		Status:      1, // 1 = brouillon, 2 = publié
		Slug:        slug,
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		City:        req.City,
		Country:     req.Country,
		ZipCode:     req.ZipCode,
		Images:      req.Images,
	}

	// Initialiser les sous-documents avec des valeurs par défaut si non fournis
	property.CheckInOut = req.CheckInOut
	if property.CheckInOut == nil {
		property.CheckInOut = &models.CheckInOut{Enabled: false}
	}

	property.Wifi = req.Wifi
	if property.Wifi == nil {
		property.Wifi = &models.Wifi{Enabled: false}
	}

	property.Equipment = req.Equipment
	if property.Equipment == nil {
		property.Equipment = &models.Equipment{Enabled: false, Items: []models.EquipmentItem{}}
	}

	property.Instructions = req.Instructions
	if property.Instructions == nil {
		property.Instructions = &models.Instructions{Enabled: false}
	}

	property.Rules = req.Rules
	if property.Rules == nil {
		property.Rules = &models.Rules{
			Enabled:         false,
			SmokingAllowed:  false,
			PetsAllowed:     false,
			PartiesAllowed:  false,
			ChildrenAllowed: true,
		}
	}

	property.Contacts = req.Contacts
	if property.Contacts == nil {
		property.Contacts = &models.Contacts{Enabled: false, Contacts: []models.Contact{}}
	}

	property.LocalRecommendations = req.LocalRecommendations
	if property.LocalRecommendations == nil {
		property.LocalRecommendations = &models.LocalRecommendations{Enabled: false, Recommendations: []models.Recommendation{}}
	}

	property.Parking = req.Parking
	if property.Parking == nil {
		property.Parking = &models.Parking{Enabled: false}
	}

	property.Transport = req.Transport
	if property.Transport == nil {
		property.Transport = &models.Transport{Enabled: false}
	}

	property.Security = req.Security
	if property.Security == nil {
		property.Security = &models.Security{Enabled: false}
	}

	property.Services = req.Services
	if property.Services == nil {
		property.Services = &models.Services{Enabled: false}
	}

	property.BabyKids = req.BabyKids
	if property.BabyKids == nil {
		property.BabyKids = &models.BabyKids{Enabled: false}
	}

	property.Pets = req.Pets
	if property.Pets == nil {
		property.Pets = &models.Pets{Enabled: false}
	}

	property.Entertainment = req.Entertainment
	if property.Entertainment == nil {
		property.Entertainment = &models.Entertainment{Enabled: false}
	}

	property.Outdoor = req.Outdoor
	if property.Outdoor == nil {
		property.Outdoor = &models.Outdoor{Enabled: false}
	}

	property.Neighborhood = req.Neighborhood
	if property.Neighborhood == nil {
		property.Neighborhood = &models.Neighborhood{Enabled: false}
	}

	property.Emergency = req.Emergency
	if property.Emergency == nil {
		property.Emergency = &models.Emergency{Enabled: false}
	}

	if err := h.propertyRepo.Create(ctx, property); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la création de la propriété",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Propriété créée avec succès",
		"property": property,
	})
}

// GetUserProperties récupère les propriétés d'un utilisateur
func (h *PropertyHandler) GetUserProperties(c *gin.Context) {
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
	includeDraft := isOwner

	// Récupérer les propriétés
	properties, err := h.propertyRepo.FindByHostID(ctx, requestedUserID, includeDraft)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erreur lors de la récupération des propriétés",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"properties": properties,
		"count":      len(properties),
	})
}

// GetProperty récupère une propriété par son slug ou ID
func (h *PropertyHandler) GetProperty(c *gin.Context) {
	identifier := c.Param("id")
	if identifier == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Identifiant manquant",
		})
		return
	}

	ctx := c.Request.Context()

	var property *models.Property
	var err error

	// Essayer d'abord comme ObjectID
	if id, err := primitive.ObjectIDFromHex(identifier); err == nil {
		property, err = h.propertyRepo.FindByID(ctx, id)
	} else {
		// Sinon, traiter comme un slug
		property, err = h.propertyRepo.FindBySlug(ctx, identifier)
	}

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Propriété introuvable",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erreur lors de la récupération de la propriété",
			})
		}
		return
	}

	// Si la propriété est en brouillon (status = 1), vérifier que l'utilisateur est le propriétaire
	if property.Status == 1 {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Propriété introuvable",
			})
			return
		}

		userID, ok := userIDInterface.(primitive.ObjectID)
		if !ok || userID != property.HostID {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Propriété introuvable",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"property": property,
	})
}

// UpdateProperty met à jour une propriété
func (h *PropertyHandler) UpdateProperty(c *gin.Context) {
	identifier := c.Param("id")
	if identifier == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Identifiant manquant",
		})
		return
	}

	ctx := c.Request.Context()

	// Trouver la propriété
	var property *models.Property
	var err error

	if id, err := primitive.ObjectIDFromHex(identifier); err == nil {
		property, err = h.propertyRepo.FindByID(ctx, id)
	} else {
		property, err = h.propertyRepo.FindBySlug(ctx, identifier)
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Propriété introuvable",
		})
		return
	}

	// Vérifier que l'utilisateur est le propriétaire
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Utilisateur non authentifié",
		})
		return
	}

	userID, ok := userIDInterface.(primitive.ObjectID)
	if !ok || userID != property.HostID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Vous n'êtes pas autorisé à modifier cette propriété",
		})
		return
	}

	// Lire les données de mise à jour
	var req models.UpdatePropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Données invalides",
			"details": err.Error(),
		})
		return
	}

	// Construire les mises à jour
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
		// Régénérer le slug si le nom change
		baseSlug := utils.GenerateSlug(req.Name)
		slug := baseSlug
		counter := 1
		for {
			exists, err := h.propertyRepo.ExistsBySlug(ctx, slug)
			if err == nil && !exists {
				break
			}
			if slug == property.Slug {
				break // Garder le slug actuel s'il est toujours valide
			}
			slug = baseSlug + "-" + strconv.Itoa(counter)
			counter++
		}
		updates["slug"] = slug
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.City != "" {
		updates["city"] = req.City
	}
	if req.Country != "" {
		updates["country"] = req.Country
	}
	if req.ZipCode != "" {
		updates["zipCode"] = req.ZipCode
	}
	if req.Images != nil {
		updates["images"] = req.Images
	}

	// Mettre à jour les sous-documents si fournis
	if req.CheckInOut != nil {
		updates["checkInOut"] = req.CheckInOut
	}
	if req.Wifi != nil {
		updates["wifi"] = req.Wifi
	}
	if req.Equipment != nil {
		updates["equipment"] = req.Equipment
	}
	if req.Instructions != nil {
		updates["instructions"] = req.Instructions
	}
	if req.Rules != nil {
		updates["rules"] = req.Rules
	}
	if req.Contacts != nil {
		updates["contacts"] = req.Contacts
	}
	if req.LocalRecommendations != nil {
		updates["localRecommendations"] = req.LocalRecommendations
	}
	if req.Parking != nil {
		updates["parking"] = req.Parking
	}
	if req.Transport != nil {
		updates["transport"] = req.Transport
	}
	if req.Security != nil {
		updates["security"] = req.Security
	}
	if req.Services != nil {
		updates["services"] = req.Services
	}
	if req.BabyKids != nil {
		updates["babyKids"] = req.BabyKids
	}
	if req.Pets != nil {
		updates["pets"] = req.Pets
	}
	if req.Entertainment != nil {
		updates["entertainment"] = req.Entertainment
	}
	if req.Outdoor != nil {
		updates["outdoor"] = req.Outdoor
	}
	if req.Neighborhood != nil {
		updates["neighborhood"] = req.Neighborhood
	}
	if req.Emergency != nil {
		updates["emergency"] = req.Emergency
	}

	if len(updates) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Aucune modification à appliquer",
		})
		return
	}

	// Appliquer les mises à jour
	if err := h.propertyRepo.Update(ctx, property.ID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la mise à jour de la propriété",
		})
		return
	}

	// Récupérer la propriété mise à jour
	updatedProperty, err := h.propertyRepo.FindByID(ctx, property.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la récupération de la propriété mise à jour",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Propriété mise à jour avec succès",
		"property": updatedProperty,
	})
}

// PublishProperty publie une propriété (change le status de 1 à 2)
func (h *PropertyHandler) PublishProperty(c *gin.Context) {
	identifier := c.Param("id")
	if identifier == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Identifiant manquant",
		})
		return
	}

	ctx := c.Request.Context()

	// Trouver la propriété
	var property *models.Property
	var err error

	if id, err := primitive.ObjectIDFromHex(identifier); err == nil {
		property, err = h.propertyRepo.FindByID(ctx, id)
	} else {
		property, err = h.propertyRepo.FindBySlug(ctx, identifier)
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Propriété introuvable",
		})
		return
	}

	// Vérifier que l'utilisateur est le propriétaire
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Utilisateur non authentifié",
		})
		return
	}

	userID, ok := userIDInterface.(primitive.ObjectID)
	if !ok || userID != property.HostID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Vous n'êtes pas autorisé à publier cette propriété",
		})
		return
	}

	// Mettre à jour le status (2 = publié)
	now := time.Now()
	updates := map[string]interface{}{
		"status":      2,
		"publishedAt": now,
	}

	if err := h.propertyRepo.Update(ctx, property.ID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la publication de la propriété",
		})
		return
	}

	// Récupérer la propriété mise à jour
	updatedProperty, err := h.propertyRepo.FindByID(ctx, property.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la récupération de la propriété mise à jour",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Propriété publiée avec succès",
		"property": updatedProperty,
	})
}

// DeleteProperty supprime une propriété
func (h *PropertyHandler) DeleteProperty(c *gin.Context) {
	identifier := c.Param("id")
	if identifier == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Identifiant manquant",
		})
		return
	}

	ctx := c.Request.Context()

	// Trouver la propriété
	var property *models.Property
	var err error

	if id, err := primitive.ObjectIDFromHex(identifier); err == nil {
		property, err = h.propertyRepo.FindByID(ctx, id)
	} else {
		property, err = h.propertyRepo.FindBySlug(ctx, identifier)
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Propriété introuvable",
		})
		return
	}

	// Vérifier que l'utilisateur est le propriétaire
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Utilisateur non authentifié",
		})
		return
	}

	userID, ok := userIDInterface.(primitive.ObjectID)
	if !ok || userID != property.HostID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Vous n'êtes pas autorisé à supprimer cette propriété",
		})
		return
	}

	// Supprimer la propriété
	if err := h.propertyRepo.Delete(ctx, property.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la suppression de la propriété",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Propriété supprimée avec succès",
	})
}
