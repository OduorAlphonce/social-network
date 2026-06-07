package main

import (
	"log"
	"net/http"

	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/api/routers"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/config"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/db"
)

func main() {
	// Configuration (using env variables or defaults)
	config.Load()

	log.Printf("Initializing database at %s...", config.App.DatabasePath)
	database, err := db.InitDB(config.App.DatabasePath, config.App.MigrationsDir)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()
	log.Println("Database initialized and migrations applied successfully!")

	// Register API routes
	r := routers.RegisterRoutes(database)

	log.Printf("Server starting on port %s...", config.App.Port)
	if err := http.ListenAndServe(":"+config.App.Port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
