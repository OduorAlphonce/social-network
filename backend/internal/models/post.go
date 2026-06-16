package models

import "time"

type Post struct {
	ID        string     `db:"id"`
	UserID    string     `db:"user_id"`
	Content   string     `db:"content"`
	Privacy   string     `db:"privacy"`
	GroupID   *string    `db:"group_id"`
	Deleted   bool       `db:"deleted"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
type PostTombstone struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}
type PostDTO struct {
	ID         string     `json:"id"`
	Content    string     `json:"content"`
	Privacy    string     `json:"privacy"`
	GroupID    *string    `json:"group_id,omitempty"`
	Deleted    bool       `json:"deleted"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	ViewerVote *string    `json:"viewer_vote"`
}