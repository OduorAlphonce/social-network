package repositories

import (
	"database/sql"
	"errors"
	"social-network/internal/models"
)

type FollowerRepository interface {
	Create(f *models.Follower) error
	UpdateStatus(followerID, followingID, status string) error
	Delete(followerID, followingID string) error
	GetStatus(followerID, followingID string) (string, error)
	GetFollowers(userID string) ([]*models.User, error)
	GetFollowing(userID string) ([]*models.User, error)
}

type sqliteFollowerRepository struct {
	db *sql.DB
}

func NewFollowerRepository(db *sql.DB) FollowerRepository {
	return &sqliteFollowerRepository{db: db}
}

func (r *sqliteFollowerRepository) Create(f *models.Follower) error {
	query := `INSERT INTO followers (follower_id, following_id, status, created_at) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, f.FollowerID, f.FollowingID, f.Status, f.CreatedAt)
	return err
}

func (r *sqliteFollowerRepository) UpdateStatus(followerID, followingID, status string) error {
	query := `UPDATE followers SET status = ? WHERE follower_id = ? AND following_id = ?`
	res, err := r.db.Exec(query, status, followerID, followingID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("follow request not found")
	}
	return nil
}

func (r *sqliteFollowerRepository) Delete(followerID, followingID string) error {
	query := `DELETE FROM followers WHERE follower_id = ? AND following_id = ?`
	res, err := r.db.Exec(query, followerID, followingID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("relationship not found")
	}
	return nil
}

func (r *sqliteFollowerRepository) GetStatus(followerID, followingID string) (string, error) {
	query := `SELECT status FROM followers WHERE follower_id = ? AND following_id = ?`
	var status string
	err := r.db.QueryRow(query, followerID, followingID).Scan(&status)
	if err == sql.ErrNoRows {
		return "none", nil
	}
	if err != nil {
		return "", err
	}
	return status, nil
}

func (r *sqliteFollowerRepository) GetFollowers(userID string) ([]*models.User, error) {
	query := `SELECT u.id, u.email, u.first_name, u.last_name, u.date_of_birth, u.avatar, u.nickname, u.about_me, u.is_public, u.created_at, u.updated_at 
	FROM followers f 
	JOIN users u ON f.follower_id = u.id 
	WHERE f.following_id = ? AND f.status = 'accepted'`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		var avatar, nickname, aboutMe sql.NullString
		err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.DateOfBirth, &avatar, &nickname, &aboutMe, &u.IsPublic, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		u.Avatar = avatar.String
		u.Nickname = nickname.String
		u.AboutMe = aboutMe.String
		users = append(users, u)
	}
	return users, nil
}

func (r *sqliteFollowerRepository) GetFollowing(userID string) ([]*models.User, error) {
	query := `SELECT u.id, u.email, u.first_name, u.last_name, u.date_of_birth, u.avatar, u.nickname, u.about_me, u.is_public, u.created_at, u.updated_at 
	FROM followers f 
	JOIN users u ON f.following_id = u.id 
	WHERE f.follower_id = ? AND f.status = 'accepted'`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		var avatar, nickname, aboutMe sql.NullString
		err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.DateOfBirth, &avatar, &nickname, &aboutMe, &u.IsPublic, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		u.Avatar = avatar.String
		u.Nickname = nickname.String
		u.AboutMe = aboutMe.String
		users = append(users, u)
	}
	return users, nil
}
