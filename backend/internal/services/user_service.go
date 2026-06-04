package services

import (
	"errors"
	"regexp"
	"time"

	"social-network/internal/models"
	"social-network/internal/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(req *models.CreateUserRequest) (*models.UserResponse, error)
	Login(email, password string) (*models.Session, error)
	Logout(sessionID string) error
	Authenticate(sessionID string) (*models.User, error)
	GetByID(id string) (*models.User, error)
}

type userService struct {
	userRepo    repositories.UserRepository
	sessionRepo repositories.SessionRepository
}

func NewUserService(ur repositories.UserRepository, sr repositories.SessionRepository) UserService {
	return &userService{
		userRepo:    ur,
		sessionRepo: sr,
	}
}

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func (s *userService) Register(req *models.CreateUserRequest) (*models.UserResponse, error) {
	// Validation
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" || req.DateOfBirth == "" {
		return nil, errors.New("missing required fields")
	}

	if !emailRegex.MatchString(req.Email) {
		return nil, errors.New("invalid email format")
	}

	if len(req.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	// Check if email already exists
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userID := uuid.New().String()
	now := time.Now()

	user := &models.User{
		ID:          userID,
		Email:       req.Email,
		Password:    string(hashedPassword),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DateOfBirth: req.DateOfBirth,
		Avatar:      req.Avatar,
		Nickname:    req.Nickname,
		AboutMe:     req.AboutMe,
		IsPublic:    req.IsPublic,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		DateOfBirth: user.DateOfBirth,
		Avatar:      user.Avatar,
		Nickname:    user.Nickname,
		AboutMe:     user.AboutMe,
		IsPublic:    user.IsPublic,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (s *userService) Login(email, password string) (*models.Session, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	sessionID := uuid.New().String()
	session := &models.Session{
		ID:        sessionID,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours validity
	}

	err = s.sessionRepo.Create(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *userService) Logout(sessionID string) error {
	return s.sessionRepo.Delete(sessionID)
}

func (s *userService) Authenticate(sessionID string) (*models.User, error) {
	session, err := s.sessionRepo.GetByID(sessionID)
	if err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		s.sessionRepo.Delete(sessionID)
		return nil, errors.New("session expired")
	}

	return s.userRepo.GetByID(session.UserID)
}

func (s *userService) GetByID(id string) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

