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
func RegisterRoutes(
	userHandler *handlers.UserHandler,
	followerHandler *handlers.FollowerHandler,
	postHandler *handlers.PostHandler,
	groupHandler *handlers.GroupHandler,
	eventHandler *handlers.EventHandler,
	chatHandler *handlers.ChatHandler,
	notificationHandler *handlers.NotificationHandler,
	authMiddleware func(http.Handler) http.Handler,
	allowedOrigin string,
) http.Handler {
	// Initialize ServeMux
	mux := http.NewServeMux()

	// Serve static uploads (for avatar files)
	uploadsDir := "./uploads"
	if err := os.MkdirAll(uploadsDir, 0o755); err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadsDir))))

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Public routes
	mux.HandleFunc("/api/users/register", userHandler.Register)
	mux.HandleFunc("/api/users/login", userHandler.Login)

	// Authenticated routes
	mux.Handle("/api/users/me", http.HandlerFunc(userHandler.Me))
	mux.Handle("/api/users/update", http.HandlerFunc(userHandler.Update))
	mux.HandleFunc("/api/users/logout", userHandler.Logout)

	mux.Handle("/api/followers/follow", authMiddleware(http.HandlerFunc(followerHandler.Follow)))
	mux.Handle("/api/followers/unfollow", authMiddleware(http.HandlerFunc(followerHandler.Unfollow)))
	mux.Handle("/api/followers/accept", authMiddleware(http.HandlerFunc(followerHandler.AcceptFollow)))
	mux.Handle("/api/followers/reject", authMiddleware(http.HandlerFunc(followerHandler.RejectFollow)))
	mux.Handle("/api/followers/followers", authMiddleware(http.HandlerFunc(followerHandler.GetFollowers)))
	mux.Handle("/api/followers/following", authMiddleware(http.HandlerFunc(followerHandler.GetFollowing)))
	mux.Handle("/api/posts", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			postHandler.CreatePost(w, r)
		} else {
			postHandler.Feed(w, r)
		}
	})))
	mux.Handle("/api/users/{id}/posts", authMiddleware(http.HandlerFunc(postHandler.ProfilePosts)))

	mux.Handle("/api/posts/{id}", authMiddleware(http.HandlerFunc(postHandler.GetSinglePost)))

	// Groups routes
	mux.Handle("/api/groups", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			groupHandler.CreateGroup(w, r)
		} else {
			groupHandler.ListGroups(w, r)
		}
	})))
	mux.Handle("/api/groups/{id}/join", authMiddleware(http.HandlerFunc(groupHandler.RequestJoin)))
	mux.Handle("/api/groups/{id}/invite", authMiddleware(http.HandlerFunc(groupHandler.InviteUser)))
	mux.Handle("/api/groups/{id}/respond", authMiddleware(http.HandlerFunc(groupHandler.RespondMembership)))
	mux.Handle("/api/groups/{id}/members", authMiddleware(http.HandlerFunc(groupHandler.ListMembers)))
	mux.Handle("/api/groups/{id}/requests", authMiddleware(http.HandlerFunc(groupHandler.ListPendingRequests)))

	// Events routes
	mux.Handle("/api/groups/{id}/events", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			eventHandler.CreateEvent(w, r)
		} else {
			eventHandler.ListEvents(w, r)
		}
	})))
	mux.Handle("/api/events/{id}/rsvp", authMiddleware(http.HandlerFunc(eventHandler.RespondEvent)))

	// Chat/Messages routes
	mux.Handle("/api/conversations", authMiddleware(http.HandlerFunc(chatHandler.GetConversations)))
	mux.Handle("/api/messages", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			chatHandler.SendMessage(w, r)
		} else {
			chatHandler.GetMessages(w, r)
		}
	})))
	mux.Handle("/api/ws", authMiddleware(http.HandlerFunc(chatHandler.HandleWS)))

	// Notifications routes
	mux.Handle("/api/notifications", authMiddleware(http.HandlerFunc(notificationHandler.ListNotifications)))
	mux.Handle("/api/notifications/{id}/read", authMiddleware(http.HandlerFunc(notificationHandler.MarkAsRead)))
	mux.Handle("/api/notifications/read/all", authMiddleware(http.HandlerFunc(notificationHandler.MarkAllAsRead)))

	return middleware.CorsMiddleware(mux)
}

