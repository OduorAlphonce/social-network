package repositories

import (
	"database/sql"
	"errors"
	"social-network/internal/models"
)

type SessionRepository interface {
	Create(session *models.Session) error
	GetByID(id string) (*models.Session, error)
	Delete(id string) error
}

type sqliteSessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) SessionRepository {
	return &sqliteSessionRepository{db: db}
}

func (r *sqliteSessionRepository) Create(s *models.Session) error {
	query := `INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, s.ID, s.UserID, s.ExpiresAt)
	return err
}

func (r *sqliteSessionRepository) GetByID(id string) (*models.Session, error) {
	query := `SELECT id, user_id, expires_at FROM sessions WHERE id = ?`
	s := &models.Session{}
	err := r.db.QueryRow(query, id).Scan(&s.ID, &s.UserID, &s.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("session not found")
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *sqliteSessionRepository) Delete(id string) error {
	query := `DELETE FROM sessions WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
