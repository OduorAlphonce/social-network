package services

import (
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
)

// PostService reads timeline, profile, and group post feeds.
type PostService interface {
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

// NewPostService creates a service for authenticated post feed reads.
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
