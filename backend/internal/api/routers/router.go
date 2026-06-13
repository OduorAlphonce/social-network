package routers

import (
	"log"
	"net/http"
	"os"

	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/api/handlers"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/api/middleware"
)

// RegisterRoutes configures the application's HTTP ServeMux.
// It maps URL paths to their corresponding handler functions, sets up static file serving for uploads,
// and wraps the router with necessary middleware like CORS and authentication.
// Note that all registered application endpoints fall under the "/api/" path.
func RegisterRoutes(userHandler *handlers.UserHandler, followerHandler *handlers.FollowerHandler, authMiddleware func(http.Handler) http.Handler) http.Handler {
	// Initialize ServeMux
	mux := http.NewServeMux()

	// Serve static uploads (for avatar files)
	uploadsDir := "./uploads"
	if err := os.MkdirAll(uploadsDir, 0o755); err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadsDir))))

	// Public routes
	mux.HandleFunc("/api/users/register", userHandler.Register)
	mux.HandleFunc("/api/users/login", userHandler.Login)
	mux.HandleFunc("/api/users/logout", userHandler.Logout)

	// Authenticated routes
	mux.Handle("/api/users/me", authMiddleware(http.HandlerFunc(userHandler.Me)))

	mux.Handle("/api/followers/follow", authMiddleware(http.HandlerFunc(followerHandler.Follow)))
	mux.Handle("/api/followers/unfollow", authMiddleware(http.HandlerFunc(followerHandler.Unfollow)))
	mux.Handle("/api/followers/accept", authMiddleware(http.HandlerFunc(followerHandler.AcceptFollow)))
	mux.Handle("/api/followers/reject", authMiddleware(http.HandlerFunc(followerHandler.RejectFollow)))
	mux.Handle("/api/followers/followers", authMiddleware(http.HandlerFunc(followerHandler.GetFollowers)))
	mux.Handle("/api/followers/following", authMiddleware(http.HandlerFunc(followerHandler.GetFollowing)))

	// Setup simple logging and CORS middleware
	return middleware.CorsMiddleware(mux)
}
