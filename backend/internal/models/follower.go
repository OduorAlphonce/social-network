package models

import "time"

type Follower struct {
	FollowerID  string    `json:"follower_id" db:"follower_id"`
	FollowingID string    `json:"following_id" db:"following_id"`
	Status      string    `json:"status" db:"status"` // 'pending', 'accepted'
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type FollowRequestInput struct {
	FollowingID string `json:"following_id"`
}

type AcceptRejectFollowInput struct {
	FollowerID string `json:"follower_id"`
}

type FollowStatusResponse struct {
	FollowerID  string `json:"follower_id"`
	FollowingID string `json:"following_id"`
	Status      string `json:"status"` // 'pending', 'accepted', 'none'
}
