package services

import (
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
)

func TestPostServiceHomeFeedPaginationDefaultsAndHasMore(t *testing.T) {
	posts := newFakePostRepository()
	viewerID := uuid.Must(uuid.FromString("10000000-0000-0000-0000-000000000001"))
	posts.homeRows = makePostRows(t, 21)
	service := NewPostService(posts, newFakeUserRepository(), newFakeFollowersRepository(), newFakeGroupMembershipRepository())

	response, err := service.GetHomeFeed(viewerID, 0, 0)
	if err != nil {
		t.Fatalf("GetHomeFeed returned error: %v", err)
	}
	if posts.lastLimit != DefaultFeedLimit+1 {
		t.Fatalf("repo limit = %d, want %d", posts.lastLimit, DefaultFeedLimit+1)
	}
	if len(response.Data) != DefaultFeedLimit {
		t.Fatalf("response post count = %d, want %d", len(response.Data), DefaultFeedLimit)
	}
	if !response.Pagination.HasMore {
		t.Fatal("expected has_more=true when repository returns limit+1 rows")
	}
	if response.Pagination.Limit != DefaultFeedLimit || response.Pagination.Offset != 0 {
		t.Fatalf("pagination = %#v, want default limit and zero offset", response.Pagination)
	}
}

func TestPostServiceRejectsInvalidPagination(t *testing.T) {
	service := NewPostService(newFakePostRepository(), newFakeUserRepository(), newFakeFollowersRepository(), newFakeGroupMembershipRepository())
	viewerID := uuid.Must(uuid.FromString("10000000-0000-0000-0000-000000000001"))

	tests := []struct {
		name   string
		limit  int
		offset int
	}{
		{name: "limit too small", limit: -1, offset: 0},
		{name: "limit too large", limit: MaxFeedLimit + 1, offset: 0},
		{name: "offset negative", limit: 20, offset: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetHomeFeed(viewerID, tt.limit, tt.offset)
			if !errors.Is(err, ErrInvalidPagination) {
				t.Fatalf("error = %v, want ErrInvalidPagination", err)
			}
		})
	}
}

func TestPostServiceProfileVisibility(t *testing.T) {
	viewerID := uuid.Must(uuid.FromString("10000000-0000-0000-0000-000000000001"))
	profileID := uuid.Must(uuid.FromString("10000000-0000-0000-0000-000000000002"))

	tests := []struct {
		name      string
		profile   *models.User
		status    models.Status
		wantError error
	}{
		{
			name:    "public profile allowed",
			profile: &models.User{ID: profileID, Email: "public@example.com", IsPublic: true},
		},
		{
			name:    "accepted follower allowed",
			profile: &models.User{ID: profileID, Email: "private@example.com", IsPublic: false},
			status:  models.Accepted,
		},
		{
			name:      "non follower forbidden",
			profile:   &models.User{ID: profileID, Email: "private@example.com", IsPublic: false},
			status:    "none",
			wantError: ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := newFakeUserRepository()
			users.add(tt.profile)
			followers := newFakeFollowersRepository()
			if tt.status != "" && tt.status != "none" {
				followers.status[followerKey{followerID: viewerID, followeeID: profileID}] = tt.status
			}
			posts := newFakePostRepository()
			posts.profileRows = makePostRows(t, 1)
			service := NewPostService(posts, users, followers, newFakeGroupMembershipRepository())

			_, err := service.GetProfilePosts(profileID, viewerID, 20, 0)
			if !errors.Is(err, tt.wantError) {
				t.Fatalf("error = %v, want %v", err, tt.wantError)
			}
		})
	}
}

func TestPostServiceProfileOwnerCanViewOwnPrivateProfile(t *testing.T) {
	userID := uuid.Must(uuid.FromString("10000000-0000-0000-0000-000000000001"))
	posts := newFakePostRepository()
	posts.profileRows = makePostRows(t, 1)
	service := NewPostService(posts, newFakeUserRepository(), newFakeFollowersRepository(), newFakeGroupMembershipRepository())

	if _, err := service.GetProfilePosts(userID, userID, 20, 0); err != nil {
		t.Fatalf("GetProfilePosts owner returned error: %v", err)
	}
}

func TestPostServiceGroupFeedRequiresAcceptedMembership(t *testing.T) {
	viewerID := uuid.Must(uuid.FromString("10000000-0000-0000-0000-000000000001"))
	groupID := uuid.Must(uuid.FromString("20000000-0000-0000-0000-000000000001"))
	groups := newFakeGroupMembershipRepository()
	service := NewPostService(newFakePostRepository(), newFakeUserRepository(), newFakeFollowersRepository(), groups)

	_, err := service.GetGroupFeed(groupID, viewerID, 20, 0)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("error = %v, want ErrForbidden", err)
	}

	groups.accepted[groupMemberKey{groupID: groupID, userID: viewerID}] = true
	if _, err := service.GetGroupFeed(groupID, viewerID, 20, 0); err != nil {
		t.Fatalf("GetGroupFeed accepted member returned error: %v", err)
	}
}

type fakePostRepository struct {
	homeRows    []*models.PostWithAuthor
	profileRows []*models.PostWithAuthor
	groupRows   []*models.PostWithAuthor
	lastLimit   int
	lastOffset  int
}

func newFakePostRepository() *fakePostRepository {
	return &fakePostRepository{}
}

func (r *fakePostRepository) CreatePost(post *models.Post) error {
	return nil
}

func (r *fakePostRepository) GetPostByID(id, viewerID uuid.UUID) (*models.PostWithAuthor, error) {
	return nil, errors.New("not implemented")
}

func (r *fakePostRepository) ListPosts(query models.PostQuery, viewerID uuid.UUID) ([]*models.PostWithAuthor, error) {
	return nil, errors.New("not implemented")
}

func (r *fakePostRepository) ListHomeFeed(viewerID uuid.UUID, limit, offset int) ([]*models.PostWithAuthor, error) {
	r.lastLimit = limit
	r.lastOffset = offset
	return r.homeRows, nil
}

func (r *fakePostRepository) ListProfilePosts(profileUserID, viewerID uuid.UUID, limit, offset int) ([]*models.PostWithAuthor, error) {
	r.lastLimit = limit
	r.lastOffset = offset
	return r.profileRows, nil
}

func (r *fakePostRepository) ListGroupFeed(groupID, viewerID uuid.UUID, limit, offset int) ([]*models.PostWithAuthor, error) {
	r.lastLimit = limit
	r.lastOffset = offset
	return r.groupRows, nil
}

type groupMemberKey struct {
	groupID uuid.UUID
	userID  uuid.UUID
}

type fakeGroupMembershipRepository struct {
	accepted map[groupMemberKey]bool
}

func newFakeGroupMembershipRepository() *fakeGroupMembershipRepository {
	return &fakeGroupMembershipRepository{accepted: map[groupMemberKey]bool{}}
}

func (r *fakeGroupMembershipRepository) IsAcceptedGroupMember(groupID, userID uuid.UUID) (bool, error) {
	return r.accepted[groupMemberKey{groupID: groupID, userID: userID}], nil
}

func makePostRows(t *testing.T, count int) []*models.PostWithAuthor {
	t.Helper()
	authorID := uuid.Must(uuid.FromString("10000000-0000-0000-0000-000000000009"))
	rows := make([]*models.PostWithAuthor, 0, count)
	for i := 0; i < count; i++ {
		postID, err := uuid.NewV4()
		if err != nil {
			t.Fatalf("uuid.NewV4 returned error: %v", err)
		}
		rows = append(rows, &models.PostWithAuthor{
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
			ViewerVote: models.ViewerVoteNone,
		})
	}
	return rows
}
