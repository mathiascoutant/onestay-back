package main

import (
	"context"
	"log"

	"onestay-back/internal/config"
	"onestay-back/internal/database"
	"onestay-back/internal/seed"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatal("Erreur lors du chargement de la configuration:", err)
	}

	if err := database.Connect(); err != nil {
		log.Fatal("Erreur lors de la connexion à MongoDB:", err)
	}
	defer database.Disconnect()

	// Supprimer tous les rôles existants
	ctx := context.Background()
	collection := database.DB.Collection("roles")
	
	result, err := collection.DeleteMany(ctx, map[string]interface{}{})
	if err != nil {
		log.Fatal("Erreur lors de la suppression des rôles:", err)
	}
	log.Printf("Supprimé %d rôles", result.DeletedCount)

	// Recréer les rôles
	if err := seed.SeedRoles(); err != nil {
		log.Fatal("Erreur lors de l'initialisation des rôles:", err)
	}

	log.Println("Rôles réinitialisés avec succès!")
}
