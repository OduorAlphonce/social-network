package main

import (
	"log"
	"net"
	"net/http"

	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/api/routers"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/config"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/db"
)

func main() {
	// Configuration (using env variables or defaults)
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize/load database
	log.Printf("Initializing database at %s...", config.App.DatabasePath)
	database, err := db.InitDB(config.App.DatabasePath, config.App.MigrationsDir)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()
	log.Println("Database initialized and migrations applied successfully!")

	// Register Routes
	handler := routers.Router(database)

	address := net.JoinHostPort(config.App.BaseAddress, config.App.Port)
	log.Printf("Server starting on %s...", address)
	if err := http.ListenAndServe(address, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
