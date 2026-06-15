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

	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/services"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/utils"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(us services.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		_ = utils.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req models.CreateUserRequest
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Parse multipart form (10 MB limit)
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			_ = utils.SendError(w, http.StatusBadRequest, "Failed to parse multipart form", nil)
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
				_ = utils.SendError(w, http.StatusInternalServerError, "Failed to read avatar file", nil)
				return
			}
			_, _ = file.Seek(0, io.SeekStart)

			fileType := http.DetectContentType(buff)
			if fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/gif" {
				_ = utils.SendError(w, http.StatusBadRequest, "Invalid file type. Only JPEG, PNG, and GIF are allowed.", nil)
				return
			}

			// Create uploads folder
			uploadsDir := "./uploads/avatars"
			err = os.MkdirAll(uploadsDir, 0755)
			if err != nil {
				_ = utils.SendError(w, http.StatusInternalServerError, "Failed to create uploads directory", nil)
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
				_ = utils.SendError(w, http.StatusInternalServerError, "Failed to save avatar", nil)
				return
			}
			defer dst.Close()

			_, err = io.Copy(dst, file)
			if err != nil {
				_ = utils.SendError(w, http.StatusInternalServerError, "Failed to write avatar file", nil)
				return
			}

			req.Avatar = "/uploads/avatars/" + newFilename
		}
	} else {
		// Handle JSON
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			_ = utils.SendError(w, http.StatusBadRequest, "Invalid request body", nil)
			return
		}
	}

	userResponse, err := h.userService.Register(&req)
	if err != nil {
		_ = utils.SendError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	_ = utils.SendSuccess(w, http.StatusCreated, "User registered successfully", userResponse)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		_ = utils.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		_ = utils.SendError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	session, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		_ = utils.SendError(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	// Set HttpOnly session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.ID.String(),
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	_ = utils.SendSuccess(w, http.StatusOK, "Login successful", map[string]string{
		"token": session.ID.String(),
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

	_ = utils.SendSuccess(w, http.StatusOK, "Logout successful", nil)
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		_ = utils.SendError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	user, err := h.userService.Authenticate(cookie.Value)
	if err != nil {
		_ = utils.SendError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	_ = utils.SendSuccess(w, http.StatusOK, "User retrieved successfully", models.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		DateOfBirth: user.DOB.Format("2006-01-02"),
		Avatar:      user.Avatar,
		Nickname:    user.Nickname,
		AboutMe:     user.AboutMe,
		IsPublic:    user.IsPublic,
		CreatedAt:   user.CreatedAt,
	})
}
