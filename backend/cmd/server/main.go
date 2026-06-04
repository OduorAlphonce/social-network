package main

import (
	"log"
	"net/http"
	"os"

	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/api/handlers"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/api/middleware"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/api/routers"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/config"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/db"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/repositories"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/services"
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

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database)
	sessionRepo := repositories.NewSessionRepository(database)
	followerRepo := repositories.NewFollowerRepository(database)

	// Initialize services
	userService := services.NewUserService(userRepo, sessionRepo)
	followerService := services.NewFollowerService(followerRepo, userRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	followerHandler := handlers.NewFollowerHandler(followerService, userService)

	// Initialize ServeMux
	mux := http.NewServeMux()

	// Register API routes
	routers.RegisterRoutes(mux, userHandler, followerHandler, userService)

	// Serve static uploads (for avatar files)
	uploadsDir := "./uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadsDir))))

	// Setup simple logging and CORS middleware
	handler := middleware.corsMiddleware(mux)

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
