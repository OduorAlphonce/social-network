package repositories

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
)

type sqliteFollowerRepository struct {
	db *sql.DB
}

func NewFollowerRepository(db *sql.DB) FollowersRepository {
	return &sqliteFollowerRepository{db: db}
}

func (r *sqliteFollowerRepository) Follow(followerID, followeeID uuid.UUID, status models.Status) error {
	query := `INSERT INTO followers (follower_id, followee_id, status, created_at) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, followerID, followeeID, string(status), time.Now())
	return err
}

func (r *sqliteFollowerRepository) Unfollow(followerID, followeeID uuid.UUID) error {
	query := `DELETE FROM followers WHERE follower_id = ? AND followee_id = ?`
	res, err := r.db.Exec(query, followerID, followeeID)
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

func (r *sqliteFollowerRepository) AcceptFollower(followerID, followeeID uuid.UUID) error {
	query := `UPDATE followers SET status = ? WHERE follower_id = ? AND followee_id = ?`
	res, err := r.db.Exec(query, string(models.Accepted), followerID, followeeID)
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

func (r *sqliteFollowerRepository) RejectFollower(followerID, followeeID uuid.UUID) error {
	query := `DELETE FROM followers WHERE follower_id = ? AND followee_id = ?`
	res, err := r.db.Exec(query, followerID, followeeID)
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

func (r *sqliteFollowerRepository) GetStatus(followerID, followeeID uuid.UUID) (models.Status, error) {
	query := `SELECT status FROM followers WHERE follower_id = ? AND followee_id = ?`
	var status string
	err := r.db.QueryRow(query, followerID, followeeID).Scan(&status)
	if err == sql.ErrNoRows {
		return "none", nil
	}
	if err != nil {
		return "", err
	}
	return models.Status(status), nil
}

func (r *sqliteFollowerRepository) GetFollowers(userID uuid.UUID) ([]*models.User, error) {
	query := `SELECT u.id, u.email, u.password_hash, u.first_name, u.last_name, u.dob, u.avatar, u.nickname, u.about_me, u.is_public, u.follower_count, u.following_count, u.created_at 
	FROM followers f 
	JOIN users u ON f.follower_id = u.id 
	WHERE f.followee_id = ? AND f.status = 'accepted'`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		var avatar, nickname, aboutMe sql.NullString
		var dob, createdAt string
		err := rows.Scan(&u.ID, &u.Email, &u.PassHash, &u.FirstName, &u.LastName, &dob, &avatar, &nickname, &aboutMe, &u.IsPublic, &u.FollowerCount, &u.FollowingCount, &createdAt)
		if err != nil {
			return nil, err
		}
		parsedDOB, err := parseSQLiteTime(dob)
		if err != nil {
			return nil, err
		}
		parsedCreatedAt, err := parseSQLiteTime(createdAt)
		if err != nil {
			return nil, err
		}
		u.DOB = parsedDOB
		u.CreatedAt = parsedCreatedAt
		u.Avatar = avatar.String
		u.Nickname = nickname.String
		u.AboutMe = aboutMe.String
		users = append(users, u)
	}
	return users, nil
}

func (r *sqliteFollowerRepository) GetFollowing(userID uuid.UUID) ([]*models.User, error) {
	query := `SELECT u.id, u.email, u.password_hash, u.first_name, u.last_name, u.dob, u.avatar, u.nickname, u.about_me, u.is_public, u.follower_count, u.following_count, u.created_at 
	FROM followers f 
	JOIN users u ON f.followee_id = u.id 
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
		var dob, createdAt string
		err := rows.Scan(&u.ID, &u.Email, &u.PassHash, &u.FirstName, &u.LastName, &dob, &avatar, &nickname, &aboutMe, &u.IsPublic, &u.FollowerCount, &u.FollowingCount, &createdAt)
		if err != nil {
			return nil, err
		}
		parsedDOB, err := parseSQLiteTime(dob)
		if err != nil {
			return nil, err
		}
		parsedCreatedAt, err := parseSQLiteTime(createdAt)
		if err != nil {
			return nil, err
		}
		u.DOB = parsedDOB
		u.CreatedAt = parsedCreatedAt
		u.Avatar = avatar.String
		u.Nickname = nickname.String
		u.AboutMe = aboutMe.String
		users = append(users, u)
	}
	return users, nil
}

func (r *sqliteFollowerRepository) GetPendingFollowers(userID uuid.UUID) ([]*models.User, error) {
	query := `SELECT u.id, u.email, u.password_hash, u.first_name, u.last_name, u.dob, u.avatar, u.nickname, u.about_me, u.is_public, u.follower_count, u.following_count, u.created_at 
	FROM followers f 
	JOIN users u ON f.follower_id = u.id 
	WHERE f.followee_id = ? AND f.status = 'pending'`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		var avatar, nickname, aboutMe sql.NullString
		var dob, createdAt string
		err := rows.Scan(&u.ID, &u.Email, &u.PassHash, &u.FirstName, &u.LastName, &dob, &avatar, &nickname, &aboutMe, &u.IsPublic, &u.FollowerCount, &u.FollowingCount, &createdAt)
		if err != nil {
			return nil, err
		}
		parsedDOB, err := parseSQLiteTime(dob)
		if err != nil {
			return nil, err
		}
		parsedCreatedAt, err := parseSQLiteTime(createdAt)
		if err != nil {
			return nil, err
		}
		u.DOB = parsedDOB
		u.CreatedAt = parsedCreatedAt
		u.Avatar = avatar.String
		u.Nickname = nickname.String
		u.AboutMe = aboutMe.String
		users = append(users, u)
	}
	return users, nil
}
