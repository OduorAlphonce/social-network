package services

import (
    "context"
    "errors"

    "learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
)

var (
    ErrPostNotFound = errors.New("post not found")
    ErrPostForbidden = errors.New("access to this post is forbidden")
)

type PostRepository interface {
    FindByID(ctx context.Context, id string) (*models.Post, error)
}

type InteractionRepository interface {
    GetViewerVote(ctx context.Context, viewerID *string, postID string) (*string, error)
}

type PermissionChecker interface {
    CanViewPost(ctx context.Context, viewerID *string, post *models.Post) (bool, error)
}

type PostService struct {
    postRepo    PostRepository
    interaction InteractionRepository
    permChecker PermissionChecker
}

func NewPostService(pr PostRepository, ir InteractionRepository, pc PermissionChecker) *PostService {
    return &PostService{
        postRepo:    pr,
        interaction: ir,
        permChecker: pc,
    }
}

func (s *PostService) GetSinglePost(ctx context.Context, postID string, viewerID *string) (any, error) {
    post, err := s.postRepo.FindByID(ctx, postID)
    if err != nil {
        return nil, ErrPostNotFound
    }

    canAccess, err := s.permChecker.CanViewPost(ctx, viewerID, post)
    if err != nil {
        return nil, err
    }
    if !canAccess {
        return nil, ErrPostForbidden
    }

    if post.Deleted {
        return models.PostTombstone{
            ID:      post.ID,
            Deleted: true,
        }, nil
    }

    var viewerVote *string
    if viewerID != nil {
        vote, err := s.interaction.GetViewerVote(ctx, viewerID, postID)
        if err == nil {
            viewerVote = vote
        }
    }

    dto := models.PostDTO{
        ID:         post.ID,
        Content:    post.Content,
        Privacy:    post.Privacy,
        GroupID:    post.GroupID,
        Deleted:    false,
        CreatedAt:  post.CreatedAt,
        UpdatedAt:  post.UpdatedAt,
        ViewerVote: viewerVote,
    }

    return dto, nil
}