package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofrs/uuid/v5"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/api/middleware"
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
