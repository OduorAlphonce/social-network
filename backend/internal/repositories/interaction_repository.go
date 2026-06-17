package repositories

import (
	"context"
	"database/sql"
)

type InteractionRepository struct {
	db *sql.DB
}

func NewInteractionRepository(db *sql.DB) *InteractionRepository {
	return &InteractionRepository{db: db}
}

func (r *InteractionRepository) GetViewerVote(ctx context.Context, viewerID *string, postID string) (*string, error) {
	return nil, nil
}