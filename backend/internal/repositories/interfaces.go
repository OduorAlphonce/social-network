package repositories

import (
	"github.com/gofrs/uuid/v5"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUserProfile(id uuid.UUID) (*models.User, error)
	DeleteUser(id uuid.UUID) error
}

type SessionRepository interface {
	CreateSession(session *models.Session) error
	GetSessionByID(id uuid.UUID) (*models.Session, error)
	DeleteSession(id uuid.UUID) error
}

type FollowersRepository interface {
	Follow(followerID, followeeID uuid.UUID, status models.Status) error
	Unfollow(followerID, followeeID uuid.UUID) error
	AcceptFollower(followerID, followeeID uuid.UUID) error
	RejectFollower(followerID, followeeID uuid.UUID) error
	GetFollowers(userID uuid.UUID) ([]*models.User, error)
	GetFollowing(userID uuid.UUID) ([]*models.User, error)
	GetStatus(followerID, followeeID uuid.UUID) (models.Status, error)
}
