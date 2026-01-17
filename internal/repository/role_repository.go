package repository

import (
	"context"
	"log"
	"time"

	"onestay-back/internal/database"
	"onestay-back/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RoleRepository struct {
	collection *mongo.Collection
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{
		collection: database.DB.Collection("roles"),
	}
}

func (r *RoleRepository) Create(ctx context.Context, role *models.Role) error {
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, role)
	return err
}

func (r *RoleRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *RoleRepository) FindBySlug(ctx context.Context, slug string) (*models.Role, error) {
	var raw bson.M
	err := r.collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&raw)
	if err != nil {
		return nil, err
	}
	role := decodeRoleFromRaw(raw)
	return &role, nil
}

func (r *RoleRepository) FindByID(ctx context.Context, id string) (*models.Role, error) {
	filter := bson.M{"_id": id}
	
	var raw bson.M
	err := r.collection.FindOne(ctx, filter).Decode(&raw)
	if err != nil {
		return nil, err
	}
	
	role := decodeRoleFromRaw(raw)
	return &role, nil
}

func (r *RoleRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"slug": slug})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RoleRepository) GetAll(ctx context.Context) ([]models.Role, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var roles []models.Role
	for cursor.Next(ctx) {
		var raw bson.M
		if err := cursor.Decode(&raw); err != nil {
			log.Printf("Erreur de décodage brut: %v", err)
			continue
		}
		
		role := decodeRoleFromRaw(raw)
		roles = append(roles, role)
	}
	
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	
	return roles, nil
}

func decodeRoleFromRaw(raw bson.M) models.Role {
	role := models.Role{
		Name: getString(raw, "name"),
		Slug: getString(raw, "slug"),
	}
	
	// Convertir l'ID string depuis le document brut
	if idVal, ok := raw["_id"]; ok {
		if idStr, ok := idVal.(string); ok {
			role.ID = idStr
		}
	}
	
	// Convertir les timestamps
	if createdAt, ok := raw["created_at"].(int64); ok {
		role.CreatedAt = time.Unix(createdAt/1000, (createdAt%1000)*1000000)
	} else if createdAt, ok := raw["created_at"].(time.Time); ok {
		role.CreatedAt = createdAt
	}
	if updatedAt, ok := raw["updated_at"].(int64); ok {
		role.UpdatedAt = time.Unix(updatedAt/1000, (updatedAt%1000)*1000000)
	} else if updatedAt, ok := raw["updated_at"].(time.Time); ok {
		role.UpdatedAt = updatedAt
	}
	
	return role
}

func getString(m bson.M, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
