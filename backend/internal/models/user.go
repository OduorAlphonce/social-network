package models

import "time"

type User struct {
	ID          string    `json:"id" db:"id"`
	Email       string    `json:"email" db:"email"`
	Password    string    `json:"-" db:"password"` // Hashed password
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	DateOfBirth string    `json:"date_of_birth" db:"date_of_birth"`
	Avatar      string    `json:"avatar,omitempty" db:"avatar"`
	Nickname    string    `json:"nickname,omitempty" db:"nickname"`
	AboutMe     string    `json:"about_me,omitempty" db:"about_me"`
	IsPublic    bool      `json:"is_public" db:"is_public"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateUserRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	Avatar      string `json:"avatar,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	AboutMe     string `json:"about_me,omitempty"`
	IsPublic    bool   `json:"is_public"`
}

type UserResponse struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth string    `json:"date_of_birth"`
	Avatar      string    `json:"avatar,omitempty"`
	Nickname    string    `json:"nickname,omitempty"`
	AboutMe     string    `json:"about_me,omitempty"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
}
