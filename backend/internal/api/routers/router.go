package routers

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/api/handlers"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/api/middleware"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/repositories"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/services"
)

func RegisterRoutes(database *sql.DB) http.Handler {
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

	// Serve static uploads (for avatar files)
	uploadsDir := "./uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadsDir))))

	
	// Public routes
	mux.HandleFunc("/api/users/register", userHandler.Register)
	mux.HandleFunc("/api/users/login", userHandler.Login)

	// Auth middleware
	auth := middleware.Auth(userService)

	// Authenticated routes
	mux.Handle("/api/users/me", auth(http.HandlerFunc(userHandler.Me)))
	mux.HandleFunc("/api/users/logout", userHandler.Logout)

	mux.Handle("/api/followers/follow", auth(http.HandlerFunc(followerHandler.Follow)))
	mux.Handle("/api/followers/unfollow", auth(http.HandlerFunc(followerHandler.Unfollow)))
	mux.Handle("/api/followers/accept", auth(http.HandlerFunc(followerHandler.AcceptFollow)))
	mux.Handle("/api/followers/reject", auth(http.HandlerFunc(followerHandler.RejectFollow)))
	mux.Handle("/api/followers/followers", auth(http.HandlerFunc(followerHandler.GetFollowers)))
	mux.Handle("/api/followers/following", auth(http.HandlerFunc(followerHandler.GetFollowing)))

	return middleware.CorsMiddleware(mux)
}
