package repositories

import (
	"database/sql"
	"errors"

	"github.com/gofrs/uuid/v5"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
)

type sqliteUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &sqliteUserRepository{db: db}
}

func (r *sqliteUserRepository) CreateUser(u *models.User) error {
	query := `INSERT INTO users (id, email, password_hash, first_name, last_name, dob, avatar, nickname, about_me, is_public, follower_count, following_count, created_at) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, u.ID, u.Email, u.PassHash, u.FirstName, u.LastName, u.DOB, u.Avatar, u.Nickname, u.AboutMe, u.IsPublic, u.FollowerCount, u.FollowingCount, u.CreatedAt)
	return err
}

func (r *sqliteUserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	query := `SELECT id, email, password_hash, first_name, last_name, dob, avatar, nickname, about_me, is_public, follower_count, following_count, created_at FROM users WHERE id = ?`
	u := &models.User{}
	var avatar, nickname, aboutMe sql.NullString
	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.Email, &u.PassHash, &u.FirstName, &u.LastName, &u.DOB, &avatar, &nickname, &aboutMe, &u.IsPublic, &u.FollowerCount, &u.FollowingCount, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	u.Avatar = avatar.String
	u.Nickname = nickname.String
	u.AboutMe = aboutMe.String
	return u, nil
}

func (r *sqliteUserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, first_name, last_name, dob, avatar, nickname, about_me, is_public, follower_count, following_count, created_at FROM users WHERE email = ?`
	u := &models.User{}
	var avatar, nickname, aboutMe sql.NullString
	err := r.db.QueryRow(query, email).Scan(&u.ID, &u.Email, &u.PassHash, &u.FirstName, &u.LastName, &u.DOB, &avatar, &nickname, &aboutMe, &u.IsPublic, &u.FollowerCount, &u.FollowingCount, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	u.Avatar = avatar.String
	u.Nickname = nickname.String
	u.AboutMe = aboutMe.String
	return u, nil
}

// UpdateUserProfile is a stub to satisfy the interface; @aloduor is implementing this.
func (r *sqliteUserRepository) UpdateUserProfile(id uuid.UUID) (*models.User, error) {
	return nil, errors.New("not implemented")
}

// DeleteUser is a stub to satisfy the interface; @fcharles is implementing this.
func (r *sqliteUserRepository) DeleteUser(id uuid.UUID) error {
	return errors.New("not implemented")
}
