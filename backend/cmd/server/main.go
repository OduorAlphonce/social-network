package main

import (
	"log"
	"net/http"
	"os"

	"social-network/internal/api/handlers"
	"social-network/internal/api/routers"
	"social-network/internal/db"
	"social-network/internal/repositories"
	"social-network/internal/services"
)

func main() {
	// Configuration (using env variables or defaults)
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./social_network.db"
	}

	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "./pkg/db/migrations/sqlite"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Initializing database at %s...", dbPath)
	database, err := db.InitDB(dbPath, migrationsDir)
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
	handler := corsMiddleware(mux)

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow CORS for development
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // Vite default port
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
