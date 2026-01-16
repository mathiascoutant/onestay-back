package main

import (
	"log"
	"net/http"

	"onestay-back/internal/config"
	"onestay-back/internal/database"
	"onestay-back/internal/router"
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

	r := router.SetupRouter()

	log.Printf("Serveur démarré sur le port %s", config.AppConfig.Port)
	if err := http.ListenAndServe(":"+config.AppConfig.Port, r); err != nil {
		log.Fatal("Erreur lors du démarrage du serveur:", err)
	}
}
