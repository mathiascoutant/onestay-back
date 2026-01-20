package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Logement struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	NomBien     string             `json:"nom_bien" bson:"nom_bien" binding:"required"`
	Description string             `json:"description" bson:"description" binding:"required"`
	Adresse     string             `json:"adresse" bson:"adresse" binding:"required"`
	Ville       string             `json:"ville" bson:"ville" binding:"required"`
	Pays        string             `json:"pays" bson:"pays" binding:"required"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id" binding:"required"`
	Status      int                `json:"status" bson:"status"` // 1 = brouillon, 2 = publié
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type CreateLogementRequest struct {
	NomBien     string `json:"nom_bien" binding:"required"`
	Description string `json:"description" binding:"required"`
	Adresse     string `json:"adresse" binding:"required"`
	Ville       string `json:"ville" binding:"required"`
	Pays        string `json:"pays" binding:"required"`
}
