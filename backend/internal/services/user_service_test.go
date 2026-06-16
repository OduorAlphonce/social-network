package services

import (
	"errors"
	"testing"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
)

func TestUserServiceRegisterRequiresProjectRegistrationFields(t *testing.T) {
	service := NewUserService(newFakeUserRepository(), newFakeSessionRepository())

	_, err := service.Register(&models.CreateUserRequest{
		Email:       "amina@example.com",
		Password:    "secret1",
		FirstName:   "Amina",
		DateOfBirth: "1998-04-12",
		IsPublic:    true,
	})

	if err == nil {
		t.Fatal("expected missing last name to be rejected")
	}
}

func TestUserServiceRegisterStoresHashAndReturnsSafeProfile(t *testing.T) {
	users := newFakeUserRepository()
	service := NewUserService(users, newFakeSessionRepository())

	response, err := service.Register(&models.CreateUserRequest{
		Email:       "amina@example.com",
		Password:    "secret1",
		FirstName:   "Amina",
		LastName:    "Njeri",
		DateOfBirth: "1998-04-12",
		Avatar:      "/uploads/avatars/amina.png",
		Nickname:    "amina",
		AboutMe:     "Weekend hiker",
		IsPublic:    true,
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	stored := users.byEmail["amina@example.com"]
	if stored == nil {
		t.Fatal("expected user to be stored")
	}
	if stored.PassHash == "secret1" {
		t.Fatal("expected password to be hashed, not stored as plaintext")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(stored.PassHash), []byte("secret1")); err != nil {
		t.Fatalf("stored hash does not match password: %v", err)
	}
	if response.Email != "amina@example.com" || response.FirstName != "Amina" || response.LastName != "Njeri" {
		t.Fatalf("unexpected response profile: %#v", response)
	}
}

func TestUserServiceLoginCreatesSessionForValidPassword(t *testing.T) {
	users := newFakeUserRepository()
	sessions := newFakeSessionRepository()
	service := NewUserService(users, sessions)

	hash, err := bcrypt.GenerateFromPassword([]byte("secret1"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword returned error: %v", err)
	}
	userID := uuid.Must(uuid.FromString("6f5d9a18-5c4f-4b7a-9e9a-7a5d2efc44b1"))
	users.add(&models.User{
		ID:        userID,
		Email:     "amina@example.com",
		PassHash:  string(hash),
		FirstName: "Amina",
		LastName:  "Njeri",
		IsPublic:  true,
	})

	session, err := service.Login("amina@example.com", "secret1")
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	if session.UserID != userID {
		t.Fatalf("session user id = %s, want %s", session.UserID, userID)
	}
	if sessions.byID[session.ID] == nil {
		t.Fatal("expected created session to be persisted")
	}
}

func TestUserServiceLoginRejectsInvalidPassword(t *testing.T) {
	users := newFakeUserRepository()
	service := NewUserService(users, newFakeSessionRepository())

	hash, err := bcrypt.GenerateFromPassword([]byte("secret1"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword returned error: %v", err)
	}
	users.add(&models.User{
		ID:        uuid.Must(uuid.FromString("6f5d9a18-5c4f-4b7a-9e9a-7a5d2efc44b1")),
		Email:     "amina@example.com",
		PassHash:  string(hash),
		FirstName: "Amina",
		LastName:  "Njeri",
		IsPublic:  true,
	})

	if _, err := service.Login("amina@example.com", "wrong-password"); err == nil {
		t.Fatal("expected invalid password to be rejected")
	}
}

type fakeUserRepository struct {
	byID    map[uuid.UUID]*models.User
	byEmail map[string]*models.User
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		byID:    map[uuid.UUID]*models.User{},
		byEmail: map[string]*models.User{},
	}
}

func (r *fakeUserRepository) add(user *models.User) {
	r.byID[user.ID] = user
	r.byEmail[user.Email] = user
}

func (r *fakeUserRepository) CreateUser(user *models.User) error {
	r.add(user)
	return nil
}

func (r *fakeUserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	user := r.byID[id]
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *fakeUserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := r.byEmail[email]
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *fakeUserRepository) UpdateUserProfile(id uuid.UUID) (*models.User, error) {
	return r.GetUserByID(id)
}

func (r *fakeUserRepository) DeleteUser(id uuid.UUID) error {
	delete(r.byID, id)
	return nil
}

type fakeSessionRepository struct {
	byID map[uuid.UUID]*models.Session
}

func newFakeSessionRepository() *fakeSessionRepository {
	return &fakeSessionRepository{byID: map[uuid.UUID]*models.Session{}}
}

func (r *fakeSessionRepository) CreateSession(session *models.Session) error {
	r.byID[session.ID] = session
	return nil
}

func (r *fakeSessionRepository) GetSessionByID(id uuid.UUID) (*models.Session, error) {
	session := r.byID[id]
	if session == nil {
		return nil, errors.New("session not found")
	}
	return session, nil
}

func (r *fakeSessionRepository) DeleteSession(id uuid.UUID) error {
	delete(r.byID, id)
	return nil
}
