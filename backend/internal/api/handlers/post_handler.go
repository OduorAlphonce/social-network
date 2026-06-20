package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofrs/uuid/v5"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/api/middleware"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/services"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/utils"
)

// PostHandler handles authenticated post feed endpoints.
type PostHandler struct {
	postService services.PostService
}

// NewPostHandler creates a handler for post feed endpoints.
func NewPostHandler(ps services.PostService) *PostHandler {
	return &PostHandler{postService: ps}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		_ = utils.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		_ = utils.SendError(w, http.StatusBadRequest, "Content-Type must be multipart/form-data", nil)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		_ = utils.SendError(w, http.StatusBadRequest, "Failed to parse multipart form", nil)
		return
	}

	currentUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		_ = utils.SendError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	privacyStr := strings.TrimSpace(r.FormValue("privacy"))
	groupIDStr := strings.TrimSpace(r.FormValue("group_id"))

	var rawAudience []string
	if vals, ok := r.Form["audience_ids"]; ok {
		rawAudience = vals
	}
	if len(rawAudience) == 1 && strings.Contains(rawAudience[0], ",") {
		rawAudience = strings.Split(rawAudience[0], ",")
	}

	var audienceIDs []uuid.UUID
	for _, idStr := range rawAudience {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		parsedID, err := uuid.FromString(idStr)
		if err != nil {
			_ = utils.SendError(w, http.StatusBadRequest, "Invalid input: malformed audience identifier", map[string]string{"audience_ids": "has an invalid format"})
			return
		}
		audienceIDs = append(audienceIDs, parsedID)
	}

	var groupID *uuid.UUID
	if groupIDStr != "" {
		gID, err := uuid.FromString(groupIDStr)
		if err != nil {
			_ = utils.SendError(w, http.StatusBadRequest, "Invalid input: malformed group identifier", map[string]string{"group_id": "has an invalid format"})
			return
		}
		groupID = &gID
	}

	imageFile, _, err := r.FormFile("image")
	hasImage := err == nil
	if hasImage {
		defer imageFile.Close()
	}

	if content == "" && !hasImage {
		_ = utils.SendError(w, http.StatusBadRequest, "Either content or image is required", map[string]string{"content": "is empty and no image uploaded"})
		return
	}

	privacy := models.PostPrivacy(privacyStr)
	if privacyStr == "" {
		privacy = models.PostPrivacyPublic
	} else if privacy != models.PostPrivacyPublic &&
		privacy != models.PostPrivacyAlmostPrivate &&
		privacy != models.PostPrivacyPrivate {
		_ = utils.SendError(w, http.StatusBadRequest, "Invalid privacy value", map[string]string{"privacy": "must be public, almost_private, or private"})
		return
	}

	var savedImagePath *string
	var success bool
	defer func() {
		if !success && savedImagePath != nil {
			_ = utils.DeleteImage(*savedImagePath)
		}
	}()

	if hasImage {
		path, err := utils.SaveImage(imageFile, "/uploads/posts/")
		if err != nil {
			_ = utils.SendError(w, http.StatusBadRequest, "Failed to save image", map[string]string{"image": err.Error()})
			return
		}
		savedImagePath = &path
	}

	req := &models.CreatePostRequest{
		Content:     content,
		Privacy:     privacy,
		GroupID:     groupID,
		AudienceIDs: audienceIDs,
		ImageURL:    savedImagePath,
	}

	response, err := h.postService.CreatePost(r.Context(), req, currentUser.ID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbidden):
			_ = utils.SendError(w, http.StatusForbidden, "Forbidden: you do not have permission to post to this group or access is denied", nil)
		case errors.Is(err, services.ErrNotFollower):
			_ = utils.SendError(w, http.StatusBadRequest, "Invalid audience: all members must be accepted followers", nil)
		case errors.Is(err, services.ErrInvalidPrivacy):
			_ = utils.SendError(w, http.StatusBadRequest, "Invalid privacy value", nil)
		default:
			_ = utils.SendError(w, http.StatusInternalServerError, "Internal server error", nil)
		}
		return
	}

	success = true
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *PostHandler) GetSinglePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postID := r.PathValue("id")

	if _, err := uuid.FromString(postID); err != nil {
		h.sendError(w, http.StatusBadRequest, "shared_validation_error: malformed id")
		return
	}

	viewerID := h.extractViewerIDFromContext(r)

	payload, err := h.postService.GetSinglePost(ctx, postID, viewerID)
	if err != nil {
		if errors.Is(err, services.ErrPostNotFound) {
			h.sendError(w, http.StatusNotFound, "Post not found")
			return
		}
		if errors.Is(err, services.ErrPostForbidden) {
			h.sendError(w, http.StatusForbidden, "You do not have access to this post")
			return
		}
		h.sendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

func (h *PostHandler) sendError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (h *PostHandler) extractViewerIDFromContext(r *http.Request) *string {
	return nil
}

// Feed returns the authenticated user's home feed or a group feed.
func (h *PostHandler) Feed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		_ = utils.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	currentUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		_ = utils.SendError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	limit, offset, err := parseFeedPagination(r)
	if err != nil {
		_ = utils.SendError(w, http.StatusBadRequest, "Invalid pagination", map[string]string{"pagination": err.Error()})
		return
	}

	groupIDParam := r.URL.Query().Get("group_id")
	if groupIDParam != "" {
		groupID, err := uuid.FromString(groupIDParam)
		if err != nil {
			_ = utils.SendError(w, http.StatusBadRequest, "Invalid input", map[string]string{"group_id": "has an invalid format"})
			return
		}
		response, err := h.postService.GetGroupFeed(groupID, currentUser.ID, limit, offset)
		h.writeFeedResponse(w, response, err)
		return
	}

	response, err := h.postService.GetHomeFeed(currentUser.ID, limit, offset)
	h.writeFeedResponse(w, response, err)
}

// ProfilePosts returns posts visible on the selected user's profile.
func (h *PostHandler) ProfilePosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		_ = utils.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	currentUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		_ = utils.SendError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	profileID, err := uuid.FromString(r.PathValue("id"))
	if err != nil {
		_ = utils.SendError(w, http.StatusBadRequest, "Invalid input", map[string]string{"id": "has an invalid format"})
		return
	}

	limit, offset, err := parseFeedPagination(r)
	if err != nil {
		_ = utils.SendError(w, http.StatusBadRequest, "Invalid pagination", map[string]string{"pagination": err.Error()})
		return
	}

	response, err := h.postService.GetProfilePosts(profileID, currentUser.ID, limit, offset)
	h.writeFeedResponse(w, response, err)
}

func (h *PostHandler) writeFeedResponse(w http.ResponseWriter, response any, err error) {
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidPagination):
			_ = utils.SendError(w, http.StatusBadRequest, "Invalid pagination", nil)
		case errors.Is(err, services.ErrForbidden):
			_ = utils.SendError(w, http.StatusForbidden, "Forbidden", nil)
		case isNotFoundError(err):
			_ = utils.SendError(w, http.StatusNotFound, "Not found", nil)
		default:
			_ = utils.SendError(w, http.StatusInternalServerError, "Internal server error", nil)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func parseFeedPagination(r *http.Request) (int, int, error) {
	limit, err := parseOptionalInt(r.URL.Query().Get("limit"))
	if err != nil {
		return 0, 0, err
	}
	offset, err := parseOptionalInt(r.URL.Query().Get("offset"))
	if err != nil {
		return 0, 0, err
	}
	return limit, offset, nil
}

func parseOptionalInt(value string) (int, error) {
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

func isNotFoundError(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "not found")
}
