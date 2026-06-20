package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/api/middleware"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/services"
)

func TestPostHandlerFeedReturnsTopLevelPagination(t *testing.T) {
	viewerID := uuid.Must(uuid.FromString("10000000-0000-0000-0000-000000000001"))
	service := &fakePostService{
		homeResponse: samplePostListResponse(t, 20, 0, true),
	}
	handler := NewPostHandler(service)
	request := authenticatedRequest(http.MethodGet, "/api/posts", viewerID)
	recorder := httptest.NewRecorder()

	handler.Feed(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var response map[string]any
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if _, ok := response["pagination"]; !ok {
		t.Fatalf("expected top-level pagination, got %#v", response)
	}
	if _, ok := response["data"].([]any); !ok {
		t.Fatalf("expected top-level data array, got %#v", response["data"])
	}
}

func TestPostHandlerFeedRejectsInvalidGroupID(t *testing.T) {
	handler := NewPostHandler(&fakePostService{})
	request := authenticatedRequest(http.MethodGet, "/api/posts?group_id=not-a-uuid", uuid.Must(uuid.NewV4()))
	recorder := httptest.NewRecorder()

	handler.Feed(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestPostHandlerFeedMapsForbiddenGroupFeed(t *testing.T) {
	handler := NewPostHandler(&fakePostService{groupErr: services.ErrForbidden})
	groupID := uuid.Must(uuid.FromString("20000000-0000-0000-0000-000000000001"))
	request := authenticatedRequest(http.MethodGet, "/api/posts?group_id="+groupID.String(), uuid.Must(uuid.NewV4()))
	recorder := httptest.NewRecorder()

	handler.Feed(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
	}
}

func TestPostHandlerProfilePostsParsesPathAndMapsForbidden(t *testing.T) {
	handler := NewPostHandler(&fakePostService{profileErr: services.ErrForbidden})
	profileID := uuid.Must(uuid.FromString("10000000-0000-0000-0000-000000000002"))
	request := authenticatedRequest(http.MethodGet, "/api/users/"+profileID.String()+"/posts", uuid.Must(uuid.NewV4()))
	request.SetPathValue("id", profileID.String())
	recorder := httptest.NewRecorder()

	handler.ProfilePosts(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
	}
}

func TestPostHandlerRejectsInvalidPagination(t *testing.T) {
	handler := NewPostHandler(&fakePostService{homeErr: services.ErrInvalidPagination})
	request := authenticatedRequest(http.MethodGet, "/api/posts?limit=101", uuid.Must(uuid.NewV4()))
	recorder := httptest.NewRecorder()

	handler.Feed(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func authenticatedRequest(method, target string, userID uuid.UUID) *http.Request {
	request := httptest.NewRequest(method, target, nil)
	user := &models.User{ID: userID, Email: "viewer@example.com"}
	ctx := context.WithValue(request.Context(), middleware.UserContextKey, user)
	return request.WithContext(ctx)
}

type fakePostService struct {
	homeResponse    *models.PostListResponse
	profileResponse *models.PostListResponse
	groupResponse   *models.PostListResponse
	createResponse  models.PostResponse
	homeErr         error
	profileErr      error
	groupErr        error
	createErr       error
}

func (s *fakePostService) CreatePost(ctx context.Context, req *models.CreatePostRequest, authorID uuid.UUID) (models.PostResponse, error) {
	if s.createErr != nil {
		return nil, s.createErr
	}
	return s.createResponse, nil
}

func (s *fakePostService) GetSinglePost(ctx context.Context, postID string, viewerID *string) (models.PostResponse, error) {
	return nil, nil
}

func (s *fakePostService) GetHomeFeed(viewerID uuid.UUID, limit, offset int) (*models.PostListResponse, error) {
	if s.homeErr != nil {
		return nil, s.homeErr
	}
	return s.homeResponse, nil
}

func (s *fakePostService) GetProfilePosts(profileUserID, viewerID uuid.UUID, limit, offset int) (*models.PostListResponse, error) {
	if s.profileErr != nil {
		return nil, s.profileErr
	}
	return s.profileResponse, nil
}

func (s *fakePostService) GetGroupFeed(groupID, viewerID uuid.UUID, limit, offset int) (*models.PostListResponse, error) {
	if s.groupErr != nil {
		return nil, s.groupErr
	}
	return s.groupResponse, nil
}

func samplePostListResponse(t *testing.T, limit, offset int, hasMore bool) *models.PostListResponse {
	t.Helper()
	authorID := uuid.Must(uuid.FromString("10000000-0000-0000-0000-000000000009"))
	postID := uuid.Must(uuid.FromString("aaaaaaaa-0000-0000-0000-000000000001"))
	post, err := models.MapPostResponse(&models.PostWithAuthor{
		Post: models.Post{
			ID:        postID,
			UserID:    &authorID,
			Content:   "post",
			Privacy:   models.PostPrivacyPublic,
			CreatedAt: time.Date(2026, 6, 16, 12, 0, 0, 0, time.UTC),
		},
		Author: &models.PublicUser{
			ID:        authorID,
			FirstName: "Amina",
			LastName:  "Njeri",
		},
	})
	if err != nil {
		t.Fatalf("MapPostResponse returned error: %v", err)
	}
	return &models.PostListResponse{
		Status:  "success",
		Message: "Posts returned.",
		Data:    []models.PostResponse{post},
		Errors:  nil,
		Pagination: models.Pagination{
			Limit:   limit,
			Offset:  offset,
			HasMore: hasMore,
		},
	}
}

var _ services.PostService = (*fakePostService)(nil)

func TestPostHandlerCreatePost(t *testing.T) {
	authorID := uuid.Must(uuid.FromString("10000000-0000-0000-0000-000000000001"))

	t.Run("Create post with invalid privacy rejected", func(t *testing.T) {
		handler := NewPostHandler(&fakePostService{})
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		_ = writer.WriteField("content", "hello")
		_ = writer.WriteField("privacy", "invalid-privacy")
		_ = writer.Close()

		request := httptest.NewRequest(http.MethodPost, "/api/posts", &body)
		request.Header.Set("Content-Type", writer.FormDataContentType())
		user := &models.User{ID: authorID, Email: "viewer@example.com"}
		request = request.WithContext(context.WithValue(request.Context(), middleware.UserContextKey, user))

		recorder := httptest.NewRecorder()
		handler.CreatePost(recorder, request)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
		}
	})

	t.Run("Create post with empty content and no image rejected", func(t *testing.T) {
		handler := NewPostHandler(&fakePostService{})
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		_ = writer.WriteField("content", "   ")
		_ = writer.Close()

		request := httptest.NewRequest(http.MethodPost, "/api/posts", &body)
		request.Header.Set("Content-Type", writer.FormDataContentType())
		user := &models.User{ID: authorID, Email: "viewer@example.com"}
		request = request.WithContext(context.WithValue(request.Context(), middleware.UserContextKey, user))

		recorder := httptest.NewRecorder()
		handler.CreatePost(recorder, request)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
		}
	})

	t.Run("Create post success", func(t *testing.T) {
		postResponse, err := models.MapPostResponse(&models.PostWithAuthor{
			Post: models.Post{
				ID:        uuid.Must(uuid.NewV4()),
				UserID:    &authorID,
				Content:   "hello",
				Privacy:   models.PostPrivacyPublic,
				CreatedAt: time.Now(),
			},
			Author: &models.PublicUser{
				ID:        authorID,
				FirstName: "Amina",
				LastName:  "Njeri",
			},
		})
		if err != nil {
			t.Fatalf("failed to map post: %v", err)
		}

		service := &fakePostService{
			createResponse: postResponse,
		}
		handler := NewPostHandler(service)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		_ = writer.WriteField("content", "hello")
		_ = writer.WriteField("privacy", "public")
		_ = writer.Close()

		request := httptest.NewRequest(http.MethodPost, "/api/posts", &body)
		request.Header.Set("Content-Type", writer.FormDataContentType())
		user := &models.User{ID: authorID, Email: "viewer@example.com"}
		request = request.WithContext(context.WithValue(request.Context(), middleware.UserContextKey, user))

		recorder := httptest.NewRecorder()
		handler.CreatePost(recorder, request)

		if recorder.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusCreated)
		}

		var respEnvelope map[string]any
		if err := json.NewDecoder(recorder.Body).Decode(&respEnvelope); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if respEnvelope["status"] != "success" {
			t.Fatalf("expected status success, got %v", respEnvelope["status"])
		}
		data, ok := respEnvelope["data"].(map[string]any)
		if !ok {
			t.Fatalf("expected data object inside response, got %T", respEnvelope["data"])
		}
		if data["content"] != "hello" {
			t.Fatalf("expected content 'hello', got %v", data["content"])
		}
	})
}
