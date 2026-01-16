package repository

import (
	"context"
	"time"

	"onestay-back/internal/database"
	"onestay-back/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	role.ID = primitive.NewObjectID()
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, role)
	return err
}

func (r *RoleRepository) FindBySlug(ctx context.Context, slug string) (*models.Role, error) {
	var role models.Role
	err := r.collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Role, error) {
	var role models.Role
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&role)
	if err != nil {
		return nil, err
	}
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
	if err := cursor.All(ctx, &roles); err != nil {
		return nil, err
	}
	return roles, nil
}
