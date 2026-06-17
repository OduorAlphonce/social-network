package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/repositories"
)

const (
	// DefaultFeedLimit is the limit used when a feed request omits limit.
	DefaultFeedLimit = 20
	// MaxFeedLimit is the largest accepted page size for post feeds.
	MaxFeedLimit = 100
)

var (
	// ErrForbidden means the current user is not allowed to access the feed.
	ErrForbidden = errors.New("forbidden")
	// ErrInvalidPagination means limit or offset is outside accepted bounds.
	ErrInvalidPagination = errors.New("invalid pagination")
	// ErrPostNotFound means the requested post does not exist.
	ErrPostNotFound = errors.New("post not found")
	// ErrPostForbidden means the current user is not allowed to access a post.
	ErrPostForbidden = errors.New("access to this post is forbidden")
)

// PostService reads timeline, profile, group, and single-post views.
type PostService interface {
	GetSinglePost(ctx context.Context, postID string, viewerID *string) (models.PostResponse, error)
	GetHomeFeed(viewerID uuid.UUID, limit, offset int) (*models.PostListResponse, error)
	GetProfilePosts(profileUserID, viewerID uuid.UUID, limit, offset int) (*models.PostListResponse, error)
	GetGroupFeed(groupID, viewerID uuid.UUID, limit, offset int) (*models.PostListResponse, error)
}

type postService struct {
	postRepo        repositories.PostRepository
	userRepo        repositories.UserRepository
	followerRepo    repositories.FollowersRepository
	groupMemberRepo repositories.GroupMembershipRepository
}

// NewPostService creates a service for authenticated post reads.
func NewPostService(
	postRepo repositories.PostRepository,
	userRepo repositories.UserRepository,
	followerRepo repositories.FollowersRepository,
	groupMemberRepo repositories.GroupMembershipRepository,
) PostService {
	return &postService{
		postRepo:        postRepo,
		userRepo:        userRepo,
		followerRepo:    followerRepo,
		groupMemberRepo: groupMemberRepo,
	}
}

func (s *postService) GetSinglePost(ctx context.Context, postID string, viewerID *string) (models.PostResponse, error) {
	id, err := uuid.FromString(postID)
	if err != nil {
		return nil, ErrPostNotFound
	}

	viewerUUID, ok := parseOptionalViewerID(viewerID)
	if !ok {
		return nil, ErrPostForbidden
	}

	row, err := s.postRepo.GetPostByID(id, viewerUUID)
	if err != nil {
		return nil, ErrPostNotFound
	}
	if err := s.canViewPost(row, viewerUUID); err != nil {
		return nil, err
	}

	post, err := models.MapPostResponse(row)
	if err != nil {
		return nil, fmt.Errorf("map single post: %w", err)
	}
	return post, nil
}

func (s *postService) GetHomeFeed(viewerID uuid.UUID, limit, offset int) (*models.PostListResponse, error) {
	page, err := normalizeFeedPagination(limit, offset)
	if err != nil {
		return nil, err
	}
	rows, err := s.postRepo.ListHomeFeed(viewerID, page.fetchLimit(), page.offset)
	if err != nil {
		return nil, err
	}
	return mapPostFeed("Posts returned.", rows, page)
}

func (s *postService) GetProfilePosts(profileUserID, viewerID uuid.UUID, limit, offset int) (*models.PostListResponse, error) {
	page, err := normalizeFeedPagination(limit, offset)
	if err != nil {
		return nil, err
	}
	if err := s.canViewProfile(profileUserID, viewerID); err != nil {
		return nil, err
	}
	rows, err := s.postRepo.ListProfilePosts(profileUserID, viewerID, page.fetchLimit(), page.offset)
	if err != nil {
		return nil, err
	}
	return mapPostFeed("Posts returned.", rows, page)
}

func (s *postService) GetGroupFeed(groupID, viewerID uuid.UUID, limit, offset int) (*models.PostListResponse, error) {
	page, err := normalizeFeedPagination(limit, offset)
	if err != nil {
		return nil, err
	}
	accepted, err := s.groupMemberRepo.IsAcceptedGroupMember(groupID, viewerID)
	if err != nil {
		return nil, err
	}
	if !accepted {
		return nil, ErrForbidden
	}
	rows, err := s.postRepo.ListGroupFeed(groupID, viewerID, page.fetchLimit(), page.offset)
	if err != nil {
		return nil, err
	}
	return mapPostFeed("Posts returned.", rows, page)
}

func (s *postService) canViewPost(row *models.PostWithAuthor, viewerID uuid.UUID) error {
	if row == nil {
		return ErrPostNotFound
	}
	if row.Post.GroupID != nil {
		accepted, err := s.groupMemberRepo.IsAcceptedGroupMember(*row.Post.GroupID, viewerID)
		if err != nil {
			return err
		}
		if !accepted {
			return ErrPostForbidden
		}
		return nil
	}
	if row.Post.UserID == nil {
		return ErrPostForbidden
	}
	if *row.Post.UserID == viewerID {
		return nil
	}
	switch row.Post.Privacy {
	case models.PostPrivacyPublic:
		return nil
	case models.PostPrivacyAlmostPrivate:
		status, err := s.followerRepo.GetStatus(viewerID, *row.Post.UserID)
		if err != nil {
			return err
		}
		if status == models.Accepted {
			return nil
		}
	case models.PostPrivacyPrivate:
		return ErrPostForbidden
	}
	return ErrPostForbidden
}

func (s *postService) canViewProfile(profileUserID, viewerID uuid.UUID) error {
	if profileUserID == viewerID {
		return nil
	}
	profileUser, err := s.userRepo.GetUserByID(profileUserID)
	if err != nil {
		return err
	}
	if profileUser.IsPublic {
		return nil
	}
	status, err := s.followerRepo.GetStatus(viewerID, profileUserID)
	if err != nil {
		return err
	}
	if status != models.Accepted {
		return ErrForbidden
	}
	return nil
}

type feedPage struct {
	limit  int
	offset int
}

func (p feedPage) fetchLimit() int {
	return p.limit + 1
}

func normalizeFeedPagination(limit, offset int) (feedPage, error) {
	if limit == 0 {
		limit = DefaultFeedLimit
	}
	if limit < 1 || limit > MaxFeedLimit || offset < 0 {
		return feedPage{}, ErrInvalidPagination
	}
	return feedPage{limit: limit, offset: offset}, nil
}

func mapPostFeed(message string, rows []*models.PostWithAuthor, page feedPage) (*models.PostListResponse, error) {
	hasMore := len(rows) > page.limit
	if hasMore {
		rows = rows[:page.limit]
	}

	posts := make([]models.PostResponse, 0, len(rows))
	for _, row := range rows {
		post, err := models.MapPostResponse(row)
		if err != nil {
			return nil, fmt.Errorf("map post feed: %w", err)
		}
		posts = append(posts, post)
	}

	return &models.PostListResponse{
		Status:  "success",
		Message: message,
		Data:    posts,
		Errors:  nil,
		Pagination: models.Pagination{
			Limit:   page.limit,
			Offset:  page.offset,
			HasMore: hasMore,
		},
	}, nil
}

func parseOptionalViewerID(viewerID *string) (uuid.UUID, bool) {
	if viewerID == nil || *viewerID == "" {
		return uuid.Nil, false
	}
	id, err := uuid.FromString(*viewerID)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}
