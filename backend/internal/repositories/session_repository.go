package repositories

import (
	"database/sql"
	"errors"

	"github.com/gofrs/uuid/v5"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
)

type sqliteSessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) SessionRepository {
	return &sqliteSessionRepository{db: db}
}

func (r *sqliteSessionRepository) CreateSession(s *models.Session) error {
	query := `INSERT INTO sessions (id, user_id, expires_at, created_at) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, s.ID, s.UserID, s.ExpiresAt, s.CreatedAt)
	return err
}

func (r *sqliteSessionRepository) GetSessionByID(id uuid.UUID) (*models.Session, error) {
	query := `SELECT id, user_id, expires_at, created_at FROM sessions WHERE id = ?`
	s := &models.Session{}
	var expiresAt, createdAt string
	err := r.db.QueryRow(query, id).Scan(&s.ID, &s.UserID, &expiresAt, &createdAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("session not found")
	}
	if err != nil {
		return nil, err
	}
	parsedExpiresAt, err := parseSQLiteTime(expiresAt)
	if err != nil {
		return nil, err
	}
	parsedCreatedAt, err := parseSQLiteTime(createdAt)
	if err != nil {
		return nil, err
	}
	s.ExpiresAt = parsedExpiresAt
	s.CreatedAt = parsedCreatedAt
	return s, nil
}

func (r *sqliteSessionRepository) DeleteSession(id uuid.UUID) error {
	query := `DELETE FROM sessions WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
