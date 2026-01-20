package handlers

import (
	"fmt"
	"net/http"
	"time"

	"onestay-back/internal/models"
	"onestay-back/internal/repository"
	"onestay-back/internal/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

	role, err := h.roleRepo.FindByID(ctx, req.RoleID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Rôle introuvable",
				"role_id": req.RoleID,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erreur lors de la recherche du rôle",
				"details": err.Error(),
			})
		}
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

func (h *AuthHandler) GetRoles(c *gin.Context) {
	ctx := c.Request.Context()
	
	roles, err := h.roleRepo.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la récupération des rôles",
		})
		return
	}
	
	var rolesResponse []gin.H
	for _, role := range roles {
		rolesResponse = append(rolesResponse, gin.H{
			"id":   role.ID,
			"name": role.Name,
			"slug": role.Slug,
		})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"roles": rolesResponse,
	})
}

func (h *AuthHandler) CreateRole(c *gin.Context) {
	var req models.CreateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Données invalides",
			"details": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	// Vérifier si le slug existe déjà
	exists, err := h.roleRepo.ExistsBySlug(ctx, req.Slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la vérification du slug",
		})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Un rôle avec ce slug existe déjà",
		})
		return
	}

	// Récupérer tous les rôles pour trouver le prochain ID
	allRoles, err := h.roleRepo.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la récupération des rôles",
		})
		return
	}

	// Trouver le plus grand ID numérique
	maxID := 0
	for _, r := range allRoles {
		var idNum int
		if _, err := fmt.Sscanf(r.ID, "%d", &idNum); err == nil {
			if idNum > maxID {
				maxID = idNum
			}
		}
	}

	// Générer le prochain ID
	nextID := fmt.Sprintf("%d", maxID+1)

	role := &models.Role{
		ID:   nextID,
		Name: req.Name,
		Slug: req.Slug,
	}

	if err := h.roleRepo.Create(ctx, role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la création du rôle",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Rôle créé avec succès",
		"role": gin.H{
			"id":         role.ID,
			"name":       role.Name,
			"slug":       role.Slug,
			"created_at": role.CreatedAt.Format(time.RFC3339),
		},
	})
}

func (h *AuthHandler) DeleteRole(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID du rôle manquant",
		})
		return
	}

	ctx := c.Request.Context()

	// Empêcher la suppression des rôles système
	if roleID == "1" || roleID == "2" || roleID == "3" || roleID == "4" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Impossible de supprimer un rôle système",
		})
		return
	}

	// Vérifier que le rôle existe
	_, err := h.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Rôle introuvable",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erreur lors de la recherche du rôle",
			})
		}
		return
	}

	if err := h.roleRepo.Delete(ctx, roleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la suppression du rôle",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rôle supprimé avec succès",
		"role_id": roleID,
	})
}

func (h *AuthHandler) GetAllUsers(c *gin.Context) {
	ctx := c.Request.Context()

	users, err := h.userRepo.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la récupération des utilisateurs",
		})
		return
	}

	var usersWithRoles []models.UserWithRole
	for _, user := range users {
		role, err := h.roleRepo.FindByID(ctx, user.RoleID)
		if err != nil {
			// Si le rôle n'est pas trouvé, créer un rôle vide
			role = &models.Role{
				ID:   user.RoleID,
				Name: "Rôle introuvable",
				Slug: "unknown",
			}
		}

		usersWithRoles = append(usersWithRoles, models.UserWithRole{
			ID:        user.ID,
			Nom:       user.Nom,
			Prenom:    user.Prenom,
			Email:     user.Email,
			Role:      *role,
			CreatedAt: user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"users": usersWithRoles,
		"count": len(usersWithRoles),
	})
}

func (h *AuthHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID utilisateur manquant",
		})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Données invalides",
			"details": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	// Vérifier que l'utilisateur existe
	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Utilisateur introuvable",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erreur lors de la recherche de l'utilisateur",
			})
		}
		return
	}

	// Préparer les mises à jour
	updates := bson.M{}

	if req.Nom != "" {
		updates["nom"] = req.Nom
	}

	if req.Prenom != "" {
		updates["prenom"] = req.Prenom
	}

	if req.Email != "" {
		// Vérifier que l'email n'est pas déjà utilisé par un autre utilisateur
		if req.Email != user.Email {
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
		}
		updates["email"] = req.Email
	}

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erreur lors du hachage du mot de passe",
			})
			return
		}
		updates["password"] = hashedPassword
	}

	if req.RoleID != "" {
		// Vérifier que le rôle existe
		_, err := h.roleRepo.FindByID(ctx, req.RoleID)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Rôle introuvable",
					"role_id": req.RoleID,
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Erreur lors de la recherche du rôle",
				})
			}
			return
		}
		updates["role_id"] = req.RoleID
	}

	// Si aucune mise à jour n'est demandée
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Aucune donnée à mettre à jour",
		})
		return
	}

	// Mettre à jour l'utilisateur
	if err := h.userRepo.Update(ctx, userID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la mise à jour de l'utilisateur",
		})
		return
	}

	// Récupérer l'utilisateur mis à jour
	updatedUser, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la récupération de l'utilisateur mis à jour",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Utilisateur mis à jour avec succès",
		"user": models.UserProfile{
			ID:        updatedUser.ID,
			Nom:       updatedUser.Nom,
			Prenom:    updatedUser.Prenom,
			Email:     updatedUser.Email,
			RoleID:    updatedUser.RoleID,
			CreatedAt: updatedUser.CreatedAt,
		},
	})
}

func (h *AuthHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID utilisateur manquant",
		})
		return
	}

	ctx := c.Request.Context()

	// Vérifier que l'utilisateur existe
	_, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Utilisateur introuvable",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erreur lors de la recherche de l'utilisateur",
			})
		}
		return
	}

	// Supprimer l'utilisateur
	if err := h.userRepo.Delete(ctx, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la suppression de l'utilisateur",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Utilisateur supprimé avec succès",
		"user_id": userID,
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
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

	// Récupérer l'utilisateur
	user, err := h.userRepo.FindByID(ctx, userID.Hex())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Utilisateur introuvable",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erreur lors de la récupération de l'utilisateur",
			})
		}
		return
	}

	// Récupérer le rôle complet
	role, err := h.roleRepo.FindByID(ctx, user.RoleID)
	if err != nil {
		// Si le rôle n'est pas trouvé, créer un rôle vide
		role = &models.Role{
			ID:   user.RoleID,
			Name: "Rôle introuvable",
			Slug: "unknown",
		}
	}

	// Retourner le profil avec le rôle complet
	c.JSON(http.StatusOK, models.UserWithRole{
		ID:        user.ID,
		Nom:       user.Nom,
		Prenom:    user.Prenom,
		Email:     user.Email,
		Role:      *role,
		CreatedAt: user.CreatedAt,
	})
}
