package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	PassHash  string    `db:"password_hash" json:"pass_hash"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	DOB       time.Time `db:"dob" json: dob`
	avatar    string    `db:"avatar" json:"avatar"`
	Nickname  string    `db:"nickname" json:"nickname"`
	AboutMe   string    `db:"about_me" json:"about_me"`

	FollowerCount  int `db:"follower_count" json:"follower_count"`
	FollowingCount int `db:"following_count" json:"following_count"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Session struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Status string

const (
	Pending  Status = "pending"
	Accepted Status = "accepted"
)

type Follower struct {
	FollowerID uuid.UUID `db:"follower_id" json:"follower_id"`
	FolloweeID uuid.UUID `db:"followee_id" json:"followee_id"`
	Status     Status    `db:"status" json:"status"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
