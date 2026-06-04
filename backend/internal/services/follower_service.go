package services

import (
	"errors"
	"time"

	"social-network/internal/models"
	"social-network/internal/repositories"
)

type FollowerService interface {
	Follow(followerID, followingID string) (string, error)
	Unfollow(followerID, followingID string) error
	AcceptFollow(followerID, followingID string) error
	RejectFollow(followerID, followingID string) error
	GetFollowers(userID string) ([]*models.User, error)
	GetFollowing(userID string) ([]*models.User, error)
	GetFollowStatus(followerID, followingID string) (string, error)
}

type followerService struct {
	followerRepo repositories.FollowerRepository
	userRepo     repositories.UserRepository
}

func NewFollowerService(fr repositories.FollowerRepository, ur repositories.UserRepository) FollowerService {
	return &followerService{
		followerRepo: fr,
		userRepo:     ur,
	}
}

func (s *followerService) Follow(followerID, followingID string) (string, error) {
	if followerID == followingID {
		return "", errors.New("cannot follow yourself")
	}

	// Verify target user exists
	targetUser, err := s.userRepo.GetByID(followingID)
	if err != nil {
		return "", errors.New("target user not found")
	}

	// Check current status
	status, err := s.followerRepo.GetStatus(followerID, followingID)
	if err != nil {
		return "", err
	}

	if status != "none" {
		return status, errors.New("follow relationship or request already exists")
	}

	// If target user is public, follow directly. If private, send follow request (pending status).
	newStatus := "pending"
	if targetUser.IsPublic {
		newStatus = "accepted"
	}

	follower := &models.Follower{
		FollowerID:  followerID,
		FollowingID: followingID,
		Status:      newStatus,
		CreatedAt:   time.Now(),
	}

	err = s.followerRepo.Create(follower)
	if err != nil {
		return "", err
	}

	return newStatus, nil
}

func (s *followerService) Unfollow(followerID, followingID string) error {
	return s.followerRepo.Delete(followerID, followingID)
}

func (s *followerService) AcceptFollow(followerID, followingID string) error {
	status, err := s.followerRepo.GetStatus(followerID, followingID)
	if err != nil {
		return err
	}
	if status != "pending" {
		return errors.New("no pending follow request to accept")
	}
	return s.followerRepo.UpdateStatus(followerID, followingID, "accepted")
}

func (s *followerService) RejectFollow(followerID, followingID string) error {
	status, err := s.followerRepo.GetStatus(followerID, followingID)
	if err != nil {
		return err
	}
	if status != "pending" {
		return errors.New("no pending follow request to reject")
	}
	return s.followerRepo.Delete(followerID, followingID)
}

func (s *followerService) GetFollowers(userID string) ([]*models.User, error) {
	return s.followerRepo.GetFollowers(userID)
}

func (s *followerService) GetFollowing(userID string) ([]*models.User, error) {
	return s.followerRepo.GetFollowing(userID)
}

func (s *followerService) GetFollowStatus(followerID, followingID string) (string, error) {
	return s.followerRepo.GetStatus(followerID, followingID)
}
