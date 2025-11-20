package auth_test

import (
	"links-shortener/internal/auth"
	"links-shortener/internal/user"
	"testing"
)

type MockUserRepo struct{}

func (m *MockUserRepo) Create(u *user.User) (*user.User, error) {
	return &user.User{
		Email: u.Email,
	}, nil
}

func (m *MockUserRepo) FindByEmail(email string) (*user.User, error) {
	return nil, nil
}

func TestAuthService_RegisterSuccess(t *testing.T) {
	const initialEmail = "test@example.com"
	authService := auth.NewAuthService(&MockUserRepo{})
	email, err := authService.Register(initialEmail, "password", "Вася")
	if err != nil {
		t.Fatal(err)
	}
	if email != initialEmail {
		t.Fatalf("expected %s, got %s", initialEmail, email)
	}
}
