package repositories

import (
	"database/sql"
	"errors"
	"social-network/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
}

type sqliteUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &sqliteUserRepository{db: db}
}

func (r *sqliteUserRepository) Create(u *models.User) error {
	query := `INSERT INTO users (id, email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, is_public, created_at, updated_at) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, u.ID, u.Email, u.Password, u.FirstName, u.LastName, u.DateOfBirth, u.Avatar, u.Nickname, u.AboutMe, u.IsPublic, u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *sqliteUserRepository) GetByID(id string) (*models.User, error) {
	query := `SELECT id, email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, is_public, created_at, updated_at FROM users WHERE id = ?`
	u := &models.User{}
	var avatar, nickname, aboutMe sql.NullString
	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.DateOfBirth, &avatar, &nickname, &aboutMe, &u.IsPublic, &u.CreatedAt, &u.UpdatedAt)
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

func (r *sqliteUserRepository) GetByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, is_public, created_at, updated_at FROM users WHERE email = ?`
	u := &models.User{}
	var avatar, nickname, aboutMe sql.NullString
	err := r.db.QueryRow(query, email).Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.DateOfBirth, &avatar, &nickname, &aboutMe, &u.IsPublic, &u.CreatedAt, &u.UpdatedAt)
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
