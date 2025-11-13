package auth

import (
	"errors"
	"links-shortener/internal/user"
)

type AuthService struct {
	UserRepo *user.UserRepository
}

func NewAuthService(userRepo *user.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

func (service *AuthService) Register(email, password, name string) (string, error) {
	existedUser, _ := service.UserRepo.FindByEmail(email)
	if existedUser != nil {
		return "", errors.New(ErrUserExists)
	}

	user := &user.User{
		Email:    email,
		Password: "",
		Name:     name,
	}

	_, err := service.UserRepo.Create(user)
	if err != nil {
		return "", err
	}
	return user.Email, nil
}
