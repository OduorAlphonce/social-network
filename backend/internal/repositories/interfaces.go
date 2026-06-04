package repositories

import (
	"github.com/gofrs/uuid/v5"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id uuid.UUID) (*models.User, error)
	UpdateUserProfile(id uuid.UUID) (*models.User, error)
	DeleteUser(id uuid.UUID) error
}

type SessionRepository interface {
	CreateSession(session *models.Session) error
	GetSessionByID(id uuid.UUID) (*models.Session, error)
	DeleteSession(id uuid.UUID) error
}

type FollowersRepository interface {
	Follow(id uuid.UUID) error
	Unfollow(id uuid.UUID) error
	AcceptFollower(status models.Status) error
	RejectFollower(status models.Status) error
	GetFollowers() ([]*models.User, error)
	GetFollowing() ([]*models.User, error)
}
