package repositories

import (
	"database/sql"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"

	"github.com/gofrs/uuid/v5"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	var userID string

	err := r.db.QueryRow(
		"SELECT id,email,password_hash,first_name,last_name,dob,avatar,nickname,about_me,follower_count,following_count,created_at FROM users WHERE id = ?", id.String(),
	).Scan(&userID, &user.Email, &user.PassHash, &user.FirstName, &user.LastName, &user.DOB, &user.Avatar, &user.Nickname, &user.AboutMe, &user.FollowerCount, &user.FollowingCount, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	user.ID, err = uuid.FromString(userID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}