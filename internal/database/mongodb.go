package database

import (
	"context"
	"log"
	"time"

	"onestay-back/internal/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client
var DB *mongo.Database

func Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.AppConfig.MongoURI)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	Client = client
	DB = client.Database(config.AppConfig.DBName)

	log.Printf("Connected to MongoDB successfully - Database: %s", config.AppConfig.DBName)
	return nil
}

func Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return Client.Disconnect(ctx)
}
