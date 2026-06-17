package repositories

import (
	"context"
	"database/sql"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
)
type PostRepository struct {
	db *sql.DB 
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) FindByID(ctx context.Context, id string) (*models.Post, error) {
	return nil, nil
}