package routers

import (
	"net/http"

	"social-network/internal/api/handlers"
	"social-network/internal/api/middleware"
	"social-network/internal/services"
)

func RegisterRoutes(mux *http.ServeMux, userHandler *handlers.UserHandler, followerHandler *handlers.FollowerHandler, userService services.UserService) {
	// Public routes
	mux.HandleFunc("/api/users/register", userHandler.Register)
	mux.HandleFunc("/api/users/login", userHandler.Login)
	mux.HandleFunc("/api/users/logout", userHandler.Logout)

	// Auth middleware
	auth := middleware.Auth(userService)

	// Authenticated routes
	mux.Handle("/api/users/me", auth(http.HandlerFunc(userHandler.Me)))

	mux.Handle("/api/followers/follow", auth(http.HandlerFunc(followerHandler.Follow)))
	mux.Handle("/api/followers/unfollow", auth(http.HandlerFunc(followerHandler.Unfollow)))
	mux.Handle("/api/followers/accept", auth(http.HandlerFunc(followerHandler.AcceptFollow)))
	mux.Handle("/api/followers/reject", auth(http.HandlerFunc(followerHandler.RejectFollow)))
	mux.Handle("/api/followers/followers", auth(http.HandlerFunc(followerHandler.GetFollowers)))
	mux.Handle("/api/followers/following", auth(http.HandlerFunc(followerHandler.GetFollowing)))
}
