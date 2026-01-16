package main

import (
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

	if err := seed.SeedRoles(); err != nil {
		log.Fatal("Erreur lors de l'initialisation des rôles:", err)
	}

	log.Println("Rôles insérés avec succès!")
}
