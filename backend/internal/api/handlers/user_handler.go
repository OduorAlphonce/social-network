package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"social-network/internal/models"
	"social-network/internal/services"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(us services.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateUserRequest
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Parse multipart form (10 MB limit)
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
			return
		}

		req.Email = r.FormValue("email")
		req.Password = r.FormValue("password")
		req.FirstName = r.FormValue("first_name")
		req.LastName = r.FormValue("last_name")
		req.DateOfBirth = r.FormValue("date_of_birth")
		req.Nickname = r.FormValue("nickname")
		req.AboutMe = r.FormValue("about_me")
		req.IsPublic = r.FormValue("is_public") == "true"

		// Handle Avatar upload
		file, handler, err := r.FormFile("avatar")
		if err == nil {
			defer file.Close()

			// Check content type
			buff := make([]byte, 512)
			_, err = file.Read(buff)
			if err != nil {
				http.Error(w, "Failed to read avatar file", http.StatusInternalServerError)
				return
			}
			_, _ = file.Seek(0, io.SeekStart)

			fileType := http.DetectContentType(buff)
			if fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/gif" {
				http.Error(w, "Invalid file type. Only JPEG, PNG, and GIF are allowed.", http.StatusBadRequest)
				return
			}

			// Create uploads folder
			uploadsDir := "./uploads/avatars"
			err = os.MkdirAll(uploadsDir, 0755)
			if err != nil {
				http.Error(w, "Failed to create uploads directory", http.StatusInternalServerError)
				return
			}

			ext := filepath.Ext(handler.Filename)
			if ext == "" {
				if fileType == "image/jpeg" {
					ext = ".jpg"
				} else if fileType == "image/png" {
					ext = ".png"
				} else if fileType == "image/gif" {
					ext = ".gif"
				}
			}

			newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
			filePath := filepath.Join(uploadsDir, newFilename)

			dst, err := os.Create(filePath)
			if err != nil {
				http.Error(w, "Failed to save avatar", http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			_, err = io.Copy(dst, file)
			if err != nil {
				http.Error(w, "Failed to write avatar file", http.StatusInternalServerError)
				return
			}

			req.Avatar = "/uploads/avatars/" + newFilename
		}
	} else {
		// Handle JSON
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
	}

	userResponse, err := h.userService.Register(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(userResponse)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	session, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Set HttpOnly session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.ID,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"token":   session.ID,
	})
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		_ = h.userService.Logout(cookie.Value)
	}

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful",
	})
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.Authenticate(cookie.Value)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(models.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		DateOfBirth: user.DateOfBirth,
		Avatar:      user.Avatar,
		Nickname:    user.Nickname,
		AboutMe:     user.AboutMe,
		IsPublic:    user.IsPublic,
		CreatedAt:   user.CreatedAt,
	})
}
