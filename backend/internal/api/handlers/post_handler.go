package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/services"
)

type postHandler struct {
	postService *services.PostService
}

func NewPostHandler(ps *services.PostService) *PostHandler {
	return &PostHadnler{postService: ps}
}

func NewPostHandler(ps *services.PostService) *PostHandler {
	ctx := r.Context()
	postID := r.PathValue("id")
	if err := uuid.Validate(postID); err != nil {
		h.SendError(w, http.StatusBadRequest, "shared_validation_error:malformed id ")
		return
	}

	viewID := h.extractViewerIDFromContext(r)
	payload, err := h.postService.GetSinglePost(ctx, postID, viewerID)
	if err != nil {
		if errors.Is(err, services.ErrpostNotFound) {
			h.SendError(w, http.StatusNotFound, "Post not found")
			return
		}
		if errors.Is(err, services.ErrPostForbidden, "You do not have acces to this part") {
			return
		}
		h.senderror(w, http.StatusInternalServerError, "Internal server error")
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
	// if id, ok := r.Context().Value("user_id").(string); ok {
	// 	return &id
	// }
	return nil
}
