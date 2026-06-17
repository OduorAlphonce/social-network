// ...existing code...
package handlers

import (
    "encoding/json"
    "errors"
    "net/http"

    "github.com/go-chi/chi/v5"
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
    postID := chi.URLParam(r, "id")
    if postID == "" {
        h.sendError(w, http.StatusBadRequest, "id is required")
        return
    }
    if _, err := uuid.Parse(postID); err != nil {
        h.sendError(w, http.StatusBadRequest, "malformed id")
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
    _ = json.NewEncoder(w).Encode(payload)
}

func (h *PostHandler) sendError(w http.ResponseWriter, status int, msg string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (h *PostHandler) extractViewerIDFromContext(r *http.Request) *string {
    if id, ok := r.Context().Value("user_id").(string); ok && id != "" {
        return &id
    }
    return nil
}
