package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Slug      string             `json:"slug" bson:"slug"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

const (
	RoleClient     = "client"
	RoleLoueur     = "loueur"
	RoleAdmin      = "admin"
	RoleSuperAdmin = "superadmin"
)
