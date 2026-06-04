package services

import (
	"errors"
	"regexp"
	"time"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/repositories"
)

type UserService interface {
	Register(req *models.CreateUserRequest) (*models.UserResponse, error)
	Login(email, password string) (*models.Session, error)
	Logout(sessionID string) error
	Authenticate(sessionID string) (*models.User, error)
	GetByID(id uuid.UUID) (*models.User, error)
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
	existingUser, _ := s.userRepo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Parse Date of Birth
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return nil, errors.New("invalid date of birth format, must be YYYY-MM-DD")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	user := &models.User{
		ID:             userID,
		Email:          req.Email,
		PassHash:       string(hashedPassword),
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		DOB:            dob,
		Avatar:         req.Avatar,
		Nickname:       req.Nickname,
		AboutMe:        req.AboutMe,
		IsPublic:       req.IsPublic,
		FollowerCount:  0,
		FollowingCount: 0,
		CreatedAt:      now,
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		DateOfBirth: user.DOB.Format("2006-01-02"),
		Avatar:      user.Avatar,
		Nickname:    user.Nickname,
		AboutMe:     user.AboutMe,
		IsPublic:    user.IsPublic,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (s *userService) Login(email, password string) (*models.Session, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	sessionID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	session := &models.Session{
		ID:        sessionID,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours validity
		CreatedAt: time.Now(),
	}

	err = s.sessionRepo.CreateSession(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *userService) Logout(sessionID string) error {
	sessUUID, err := uuid.FromString(sessionID)
	if err != nil {
		return err
	}
	return s.sessionRepo.DeleteSession(sessUUID)
}

func (s *userService) Authenticate(sessionID string) (*models.User, error) {
	sessUUID, err := uuid.FromString(sessionID)
	if err != nil {
		return nil, err
	}

	session, err := s.sessionRepo.GetSessionByID(sessUUID)
	if err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		_ = s.sessionRepo.DeleteSession(sessUUID)
		return nil, errors.New("session expired")
	}

	return s.userRepo.GetUserByID(session.UserID)
}

func (s *userService) GetByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}
