package services

import (
	"context" 
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
)
type PermissionService struct {

}

func NewPermissionService(followerRepo any) *PermissionService {
	return &PermissionService{}
}

func (s *PermissionService) CanViewPost(ctx context.Context, viewerID *string, post *models.Post) (bool, error) {
	return true, nil
}