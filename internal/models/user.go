package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Nom       string             `json:"nom" bson:"nom" binding:"required"`
	Prenom    string             `json:"prenom" bson:"prenom" binding:"required"`
	Email     string             `json:"email" bson:"email" binding:"required,email"`
	Password  string             `json:"-" bson:"password" binding:"required,min=6"`
	RoleID    string             `json:"role_id" bson:"role_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type RegisterRequest struct {
	Nom      string `json:"nom" binding:"required"`
	Prenom   string `json:"prenom" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   string `json:"role_id" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User UserProfile `json:"user"`
}

type UserProfile struct {
	ID        primitive.ObjectID `json:"id"`
	Nom       string             `json:"nom"`
	Prenom    string             `json:"prenom"`
	Email     string             `json:"email"`
	RoleID    string             `json:"role_id"`
	CreatedAt time.Time          `json:"created_at"`
}
