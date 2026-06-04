package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"social-network/internal/api/middleware"
	"social-network/internal/models"
	"social-network/internal/services"
)


type FollowerHandler struct {
	followerService services.FollowerService
	userService     services.UserService
}

func NewFollowerHandler(fs services.FollowerService, us services.UserService) *FollowerHandler {
	return &FollowerHandler{
		followerService: fs,
		userService:     us,
	}
}

func (h *FollowerHandler) Follow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input models.FollowRequestInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.FollowingID == "" {
		http.Error(w, "Invalid input. following_id is required.", http.StatusBadRequest)
		return
	}

	status, err := h.followerService.Follow(currentUser.ID, input.FollowingID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Follow request processed",
		"status":  status,
	})
}

func (h *FollowerHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input models.FollowRequestInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.FollowingID == "" {
		http.Error(w, "Invalid input. following_id is required.", http.StatusBadRequest)
		return
	}

	err = h.followerService.Unfollow(currentUser.ID, input.FollowingID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Unfollowed successfully",
	})
}

func (h *FollowerHandler) AcceptFollow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input models.AcceptRejectFollowInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.FollowerID == "" {
		http.Error(w, "Invalid input. follower_id is required.", http.StatusBadRequest)
		return
	}

	err = h.followerService.AcceptFollow(input.FollowerID, currentUser.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Follow request accepted",
	})
}

func (h *FollowerHandler) RejectFollow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input models.AcceptRejectFollowInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.FollowerID == "" {
		http.Error(w, "Invalid input. follower_id is required.", http.StatusBadRequest)
		return
	}

	err = h.followerService.RejectFollow(input.FollowerID, currentUser.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Follow request rejected",
	})
}

func (h *FollowerHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Determine whose followers to list (default to self, or get from query param "user_id")
	targetUserID := r.URL.Query().Get("user_id")
	if targetUserID == "" {
		targetUserID = currentUser.ID
	}

	// Verify permission if listing another user's followers
	if targetUserID != currentUser.ID {
		err := h.verifyAccess(currentUser.ID, targetUserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
	}

	followers, err := h.followerService.GetFollowers(targetUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map to UserResponse to avoid exposing passwords
	var response []*models.UserResponse
	for _, f := range followers {
		response = append(response, &models.UserResponse{
			ID:          f.ID,
			Email:       f.Email,
			FirstName:   f.FirstName,
			LastName:    f.LastName,
			DateOfBirth: f.DateOfBirth,
			Avatar:      f.Avatar,
			Nickname:    f.Nickname,
			AboutMe:     f.AboutMe,
			IsPublic:    f.IsPublic,
			CreatedAt:   f.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func (h *FollowerHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	targetUserID := r.URL.Query().Get("user_id")
	if targetUserID == "" {
		targetUserID = currentUser.ID
	}

	// Verify permission if listing another user's following list
	if targetUserID != currentUser.ID {
		err := h.verifyAccess(currentUser.ID, targetUserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
	}

	following, err := h.followerService.GetFollowing(targetUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response []*models.UserResponse
	for _, f := range following {
		response = append(response, &models.UserResponse{
			ID:          f.ID,
			Email:       f.Email,
			FirstName:   f.FirstName,
			LastName:    f.LastName,
			DateOfBirth: f.DateOfBirth,
			Avatar:      f.Avatar,
			Nickname:    f.Nickname,
			AboutMe:     f.AboutMe,
			IsPublic:    f.IsPublic,
			CreatedAt:   f.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func (h *FollowerHandler) verifyAccess(currentUserID, targetUserID string) error {
	targetUser, err := h.userService.GetByID(targetUserID)
	if err != nil {
		return errors.New("user not found")
	}

	if targetUser.IsPublic {
		return nil
	}

	// If the profile is private, currentUserID must follow targetUserID with status = 'accepted'
	status, err := h.followerService.GetFollowStatus(currentUserID, targetUserID)
	if err != nil {
		return err
	}

	if status != "accepted" {
		return errors.New("profile is private. You must follow this user to view their activity.")
	}

	return nil
}

