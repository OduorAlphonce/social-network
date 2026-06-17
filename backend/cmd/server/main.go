package main

import (
	"log"
	"net"
	"net/http"

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
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Initializing database at %s...", config.App.DatabasePath)
	database, err := db.InitDB(config.App.DatabasePath, config.App.MigrationsDir)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	log.Println("Database initialized and migrations applied successfully!")

	userRepo := repositories.NewUserRepository(database)
	sessionRepo := repositories.NewSessionRepository(database)
	followerRepo := repositories.NewFollowerRepository(database)
	postRepo := repositories.NewPostRepository(database)
	groupMembershipRepo := repositories.NewGroupMembershipRepository(database)

	userService := services.NewUserService(userRepo, sessionRepo)
	followerService := services.NewFollowerService(followerRepo, userRepo)
	postService := services.NewPostService(postRepo, userRepo, followerRepo, groupMembershipRepo)

	userHandler := handlers.NewUserHandler(userService)
	followerHandler := handlers.NewFollowerHandler(followerService, userService)
	postHandler := handlers.NewPostHandler(postService)
	authMiddleware := middleware.Auth(userService)

	// Register Routes
	handler := routers.RegisterRoutes(userHandler, followerHandler, postHandler, authMiddleware, config.App.AllowedOrigin)

	address := net.JoinHostPort(config.App.BaseAddress, config.App.Port)
	log.Printf("Server starting on %s...", address)
	if err := http.ListenAndServe(address, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
