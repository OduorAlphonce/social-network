package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/services"
)

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler(ps *services.PostService) *PostHandler {
	return &PostHandler{postService: ps}
}

func (h *PostHandler) GetSinglePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	postID := r.PathValue("id")

	if err := uuid.Validate(postID); err != nil {
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