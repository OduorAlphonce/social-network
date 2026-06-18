package services

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/repositories"
)

// PermissionService centralizes access control checks across the social network.
type PermissionService struct {
	followerRepo    repositories.FollowersRepository
	userRepo        repositories.UserRepository
	groupMemberRepo repositories.GroupMembershipRepository
}

// NewPermissionService instantiates a new PermissionService with the required repositories.
func NewPermissionService(
	followerRepo repositories.FollowersRepository,
	userRepo repositories.UserRepository,
	groupMemberRepo repositories.GroupMembershipRepository,
) *PermissionService {
	return &PermissionService{
		followerRepo:    followerRepo,
		userRepo:        userRepo,
		groupMemberRepo: groupMemberRepo,
	}
}

// CanViewPost checks if a viewer has access to read a specific post.
func (s *PermissionService) CanViewPost(ctx context.Context, viewerID *string, post *models.Post) (bool, error) {
	if post == nil {
		return false, nil
	}

	var viewerUUID uuid.UUID
	if viewerID != nil && *viewerID != "" {
		id, err := uuid.FromString(*viewerID)
		if err != nil {
			return false, nil
		}
		viewerUUID = id
	}

	// Group posts require accepted membership in the group
	if post.GroupID != nil {
		if viewerUUID == uuid.Nil {
			return false, nil
		}
		return s.groupMemberRepo.IsAcceptedGroupMember(*post.GroupID, viewerUUID)
	}

	if post.UserID == nil {
		return false, nil
	}

	// A user can always view their own post
	if *post.UserID == viewerUUID {
		return true, nil
	}

	switch post.Privacy {
	case models.PostPrivacyPublic:
		return true, nil
	case models.PostPrivacyAlmostPrivate:
		if viewerUUID == uuid.Nil {
			return false, nil
		}
		status, err := s.followerRepo.GetStatus(viewerUUID, *post.UserID)
		if err != nil {
			return false, err
		}
		return status == models.Accepted, nil
	case models.PostPrivacyPrivate:
		// Private posts are only visible to the owner (handled above)
		return false, nil
	}

	return false, nil
}

// IsFollowing checks if followerID is currently following followeeID with an accepted status.
func (s *PermissionService) IsFollowing(ctx context.Context, followerID, followeeID uuid.UUID) (bool, error) {
	status, err := s.followerRepo.GetStatus(followerID, followeeID)
	if err != nil {
		return false, err
	}
	return status == models.Accepted, nil
}

// CanViewProfile checks if viewerID is authorized to view profileUserID's profile.
func (s *PermissionService) CanViewProfile(ctx context.Context, viewerID, profileUserID uuid.UUID) (bool, error) {
	if viewerID == profileUserID {
		return true, nil
	}
	profileUser, err := s.userRepo.GetUserByID(profileUserID)
	if err != nil {
		return false, err
	}
	if profileUser.IsPublic {
		return true, nil
	}
	status, err := s.followerRepo.GetStatus(viewerID, profileUserID)
	if err != nil {
		return false, err
	}
	return status == models.Accepted, nil
}

// IsGroupMember checks if userID is an accepted member of groupID.
func (s *PermissionService) IsGroupMember(ctx context.Context, groupID, userID uuid.UUID) (bool, error) {
	return s.groupMemberRepo.IsAcceptedGroupMember(groupID, userID)
}
