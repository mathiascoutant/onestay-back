package repository

import (
	"context"
	"time"

	"onestay-back/internal/database"
	"onestay-back/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PropertyRepository struct {
	collection *mongo.Collection
}

func NewPropertyRepository() *PropertyRepository {
	return &PropertyRepository{
		collection: database.DB.Collection("properties"),
	}
}

// Create crée une nouvelle propriété
func (r *PropertyRepository) Create(ctx context.Context, property *models.Property) error {
	property.ID = primitive.NewObjectID()
	property.CreatedAt = time.Now()
	property.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, property)
	return err
}

// ExistsBySlug vérifie si un slug existe déjà
func (r *PropertyRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"slug": slug,
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByNameAndHostID vérifie si un logement avec le même nom existe déjà pour cet hôte
func (r *PropertyRepository) ExistsByNameAndHostID(ctx context.Context, name string, hostID primitive.ObjectID) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"name":   name,
		"hostId": hostID,
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindBySlug trouve une propriété par son slug
func (r *PropertyRepository) FindBySlug(ctx context.Context, slug string) (*models.Property, error) {
	var property models.Property
	err := r.collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&property)
	if err != nil {
		return nil, err
	}
	return &property, nil
}

// FindByID trouve une propriété par son ID
func (r *PropertyRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Property, error) {
	var property models.Property
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&property)
	if err != nil {
		return nil, err
	}
	return &property, nil
}

// FindByHostID trouve toutes les propriétés d'un hôte
func (r *PropertyRepository) FindByHostID(ctx context.Context, hostID primitive.ObjectID, includeDraft bool) ([]models.Property, error) {
	filter := bson.M{"hostId": hostID}
	
	// Si on ne doit pas inclure les brouillons, filtrer uniquement les publiés (status = 2)
	if !includeDraft {
		filter["status"] = 2
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var properties []models.Property
	if err := cursor.All(ctx, &properties); err != nil {
		return nil, err
	}
	return properties, nil
}

// Update met à jour une propriété
func (r *PropertyRepository) Update(ctx context.Context, id primitive.ObjectID, updates bson.M) error {
	updates["updatedAt"] = time.Now()
	
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updates},
	)
	return err
}

// Delete supprime une propriété
func (r *PropertyRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// DeleteByHostID supprime toutes les propriétés d'un hôte
func (r *PropertyRepository) DeleteByHostID(ctx context.Context, hostID primitive.ObjectID) (int64, error) {
	result, err := r.collection.DeleteMany(ctx, bson.M{"hostId": hostID})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

// FindAll trouve toutes les propriétés publiées (pour recherche publique)
func (r *PropertyRepository) FindAll(ctx context.Context, limit, skip int64) ([]models.Property, error) {
	filter := bson.M{"status": 2} // 2 = publié
	
	opts := options.Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(bson.M{"createdAt": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var properties []models.Property
	if err := cursor.All(ctx, &properties); err != nil {
		return nil, err
	}
	return properties, nil
}
