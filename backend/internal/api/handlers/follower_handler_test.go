package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
)

func TestFollowerHandlerFollowValidatesInputAndCallsService(t *testing.T) {
	currentUserID := uuid.Must(uuid.FromString("10000000-0000-0000-0000-000000000501"))
	targetID := uuid.Must(uuid.FromString("10000000-0000-0000-0000-000000000502"))
	followers := &handlerFakeFollowerService{}
	handler := NewFollowerHandler(followers, &handlerFakeUserService{})
	request := authenticatedRequest(http.MethodPost, "/api/followers/follow", currentUserID)
	request.Body = io.NopCloser(bytes.NewBufferString(`{"following_id":"` + targetID.String() + `"}`))
	recorder := httptest.NewRecorder()

	handler.Follow(recorder, request)

	if recorder.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusAccepted, recorder.Body.String())
	}
	if followers.lastFollowFollowerID != currentUserID || followers.lastFollowFollowingID != targetID {
		t.Fatalf("follow args = %s/%s", followers.lastFollowFollowerID, followers.lastFollowFollowingID)
	}
}

func TestFollowerHandlerFollowRejectsInvalidUUID(t *testing.T) {
	handler := NewFollowerHandler(&handlerFakeFollowerService{}, &handlerFakeUserService{})
	request := authenticatedRequest(http.MethodPost, "/api/followers/follow", uuid.Must(uuid.NewV4()))
	request.Body = io.NopCloser(bytes.NewBufferString(`{"following_id":"not-a-uuid"}`))
	recorder := httptest.NewRecorder()

	handler.Follow(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestFollowerHandlerAcceptFollowCallsService(t *testing.T) {
	currentUserID := uuid.Must(uuid.FromString("10000000-0000-0000-0000-000000000501"))
	followerID := uuid.Must(uuid.FromString("10000000-0000-0000-0000-000000000503"))
	followers := &handlerFakeFollowerService{}
	handler := NewFollowerHandler(followers, &handlerFakeUserService{})
	request := authenticatedRequest(http.MethodPost, "/api/followers/accept", currentUserID)
	request.Body = io.NopCloser(bytes.NewBufferString(`{"follower_id":"` + followerID.String() + `"}`))
	recorder := httptest.NewRecorder()

	handler.AcceptFollow(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if followers.lastAcceptFollowerID != followerID || followers.lastAcceptFollowingID != currentUserID {
		t.Fatalf("accept args = %s/%s", followers.lastAcceptFollowerID, followers.lastAcceptFollowingID)
	}
}

func TestFollowerHandlerGetFollowersReturnsSafeProfiles(t *testing.T) {
	currentUserID := uuid.Must(uuid.FromString("10000000-0000-0000-0000-000000000501"))
	followerID := uuid.Must(uuid.FromString("10000000-0000-0000-0000-000000000504"))
	followers := &handlerFakeFollowerService{
		followers: []*models.User{{
			ID:        followerID,
			Email:     "follower@example.com",
			PassHash:  "secret-hash",
			FirstName: "Follower",
			LastName:  "User",
			DOB:       time.Date(1999, 5, 13, 0, 0, 0, 0, time.UTC),
			IsPublic:  true,
		}},
	}
	handler := NewFollowerHandler(followers, &handlerFakeUserService{})
	request := authenticatedRequest(http.MethodGet, "/api/followers/followers", currentUserID)
	recorder := httptest.NewRecorder()

	handler.GetFollowers(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if followers.lastGetFollowersID != currentUserID {
		t.Fatalf("GetFollowers id = %s, want %s", followers.lastGetFollowersID, currentUserID)
	}
	var envelope responseEnvelope
	decodeHandlerResponse(t, recorder, &envelope)
	var users []map[string]any
	if err := json.Unmarshal(envelope.Data, &users); err != nil {
		t.Fatalf("decode users: %v", err)
	}
	if len(users) != 1 || users[0]["email"] != "follower@example.com" {
		t.Fatalf("users = %#v", users)
	}
	if _, ok := users[0]["pass_hash"]; ok {
		t.Fatalf("response leaked pass_hash: %#v", users[0])
	}
}

func TestFollowerHandlerGetFollowingForPrivateProfileRequiresAcceptedFollow(t *testing.T) {
	currentUserID := uuid.Must(uuid.FromString("10000000-0000-0000-0000-000000000501"))
	targetID := uuid.Must(uuid.FromString("10000000-0000-0000-0000-000000000505"))
	users := &handlerFakeUserService{
		getByIDUser: &models.User{ID: targetID, Email: "private@example.com", IsPublic: false},
	}
	followers := &handlerFakeFollowerService{followStatus: "none"}
	handler := NewFollowerHandler(followers, users)
	request := authenticatedRequest(http.MethodGet, "/api/followers/following?user_id="+targetID.String(), currentUserID)
	recorder := httptest.NewRecorder()

	handler.GetFollowing(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
	}
	if users.lastGetByID != targetID {
		t.Fatalf("GetByID id = %s, want %s", users.lastGetByID, targetID)
	}
	if followers.lastStatusFollowerID != currentUserID || followers.lastStatusFollowingID != targetID {
		t.Fatalf("status args = %s/%s", followers.lastStatusFollowerID, followers.lastStatusFollowingID)
	}
}

type handlerFakeFollowerService struct {
	followers               []*models.User
	following               []*models.User
	followStatus            string
	followErr               error
	unfollowErr             error
	acceptErr               error
	rejectErr               error
	getFollowersErr         error
	getFollowingErr         error
	statusErr               error
	lastFollowFollowerID    uuid.UUID
	lastFollowFollowingID   uuid.UUID
	lastUnfollowFollowerID  uuid.UUID
	lastUnfollowFollowingID uuid.UUID
	lastAcceptFollowerID    uuid.UUID
	lastAcceptFollowingID   uuid.UUID
	lastRejectFollowerID    uuid.UUID
	lastRejectFollowingID   uuid.UUID
	lastGetFollowersID      uuid.UUID
	lastGetFollowingID      uuid.UUID
	lastStatusFollowerID    uuid.UUID
	lastStatusFollowingID   uuid.UUID
}

func (s *handlerFakeFollowerService) Follow(followerID, followingID uuid.UUID) (string, error) {
	s.lastFollowFollowerID = followerID
	s.lastFollowFollowingID = followingID
	if s.followErr != nil {
		return "", s.followErr
	}
	return "accepted", nil
}

func (s *handlerFakeFollowerService) Unfollow(followerID, followingID uuid.UUID) error {
	s.lastUnfollowFollowerID = followerID
	s.lastUnfollowFollowingID = followingID
	return s.unfollowErr
}

func (s *handlerFakeFollowerService) AcceptFollow(followerID, followingID uuid.UUID) error {
	s.lastAcceptFollowerID = followerID
	s.lastAcceptFollowingID = followingID
	return s.acceptErr
}

func (s *handlerFakeFollowerService) RejectFollow(followerID, followingID uuid.UUID) error {
	s.lastRejectFollowerID = followerID
	s.lastRejectFollowingID = followingID
	return s.rejectErr
}

func (s *handlerFakeFollowerService) GetFollowers(userID uuid.UUID) ([]*models.User, error) {
	s.lastGetFollowersID = userID
	if s.getFollowersErr != nil {
		return nil, s.getFollowersErr
	}
	return s.followers, nil
}

func (s *handlerFakeFollowerService) GetFollowing(userID uuid.UUID) ([]*models.User, error) {
	s.lastGetFollowingID = userID
	if s.getFollowingErr != nil {
		return nil, s.getFollowingErr
	}
	return s.following, nil
}

func (s *handlerFakeFollowerService) GetFollowStatus(followerID, followingID uuid.UUID) (string, error) {
	s.lastStatusFollowerID = followerID
	s.lastStatusFollowingID = followingID
	if s.statusErr != nil {
		return "", s.statusErr
	}
	if s.followStatus == "" {
		return "none", nil
	}
	return s.followStatus, nil
}

var errFollowerHandlerTest = errors.New("follower handler test error")
