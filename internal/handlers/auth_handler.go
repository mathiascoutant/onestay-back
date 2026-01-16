package handlers

import (
	"net/http"
	"time"

	"onestay-back/internal/models"
	"onestay-back/internal/repository"
	"onestay-back/internal/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		userRepo: repository.NewUserRepository(),
		roleRepo: repository.NewRoleRepository(),
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Données invalides",
			"details": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	exists, err := h.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la vérification de l'email",
		})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Cet email est déjà utilisé",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors du hachage du mot de passe",
		})
		return
	}

	roleID, err := primitive.ObjectIDFromHex(req.RoleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Format de role_id invalide",
		})
		return
	}

	role, err := h.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Rôle introuvable",
		})
		return
	}

	user := &models.User{
		Nom:      req.Nom,
		Prenom:   req.Prenom,
		Email:    req.Email,
		Password: hashedPassword,
		RoleID:   role.ID,
	}

	if err := h.userRepo.Create(ctx, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la création du compte",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Compte créé avec succès",
		"user": gin.H{
			"id":         user.ID,
			"nom":        user.Nom,
			"prenom":     user.Prenom,
			"email":      user.Email,
			"role_id":    user.RoleID,
			"created_at": user.CreatedAt.Format(time.RFC3339),
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Données invalides",
			"details": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	user, err := h.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email ou mot de passe incorrect",
		})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email ou mot de passe incorrect",
		})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.RoleID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la génération du token",
		})
		return
	}

	c.Header("Authorization", "Bearer "+token)

	response := models.LoginResponse{
		User: models.UserProfile{
			ID:        user.ID,
			Nom:       user.Nom,
			Prenom:    user.Prenom,
			Email:     user.Email,
			RoleID:    user.RoleID,
			CreatedAt: user.CreatedAt,
		},
	}

	c.JSON(http.StatusOK, response)
}
