package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Group struct {
	ID          uuid.UUID `db:"id" json:"id"`
	CreatorID   uuid.UUID `db:"creator_id" json:"creator_id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type GroupMember struct {
	GroupID uuid.UUID `db:"group_id" json:"group_id"`
	UserID  uuid.UUID `db:"user_id" json:"user_id"`
	Status  string    `db:"status" json:"status"` // 'pending_invite', 'pending_request', 'accepted'
}

type CreateGroupRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type InviteUserRequest struct {
	UserID string `json:"user_id"`
}

type GroupResponse struct {
	ID          uuid.UUID `json:"id"`
	CreatorID   uuid.UUID `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	IsMember    bool      `json:"is_member"`
	Status      string    `json:"status,omitempty"` // membership status
}
