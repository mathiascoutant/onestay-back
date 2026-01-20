package repository

import (
	"context"
	"time"

	"onestay-back/internal/database"
	"onestay-back/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type LogementRepository struct {
	collection *mongo.Collection
}

func NewLogementRepository() *LogementRepository {
	return &LogementRepository{
		collection: database.DB.Collection("logements"),
	}
}

func (r *LogementRepository) Create(ctx context.Context, logement *models.Logement) error {
	logement.ID = primitive.NewObjectID()
	logement.CreatedAt = time.Now()
	logement.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, logement)
	return err
}

func (r *LogementRepository) ExistsByNomBienAndUserID(ctx context.Context, nomBien string, userID primitive.ObjectID) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"nom_bien": nomBien,
		"user_id":  userID,
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *LogementRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID, includeBrouillon bool) ([]models.Logement, error) {
	filter := bson.M{"user_id": userID}
	
	// Si on ne doit pas inclure les brouillons, filtrer uniquement les publiés (status = 2)
	if !includeBrouillon {
		filter["status"] = 2
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logements []models.Logement
	if err := cursor.All(ctx, &logements); err != nil {
		return nil, err
	}
	return logements, nil
}
