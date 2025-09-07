package services

import (
	"context"
	"errors"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/ports"
)

type UserSerivce struct {
	repo           ports.UserRepository
	passwordHasher ports.PasswordHasher
}

func NewUserService(repo ports.UserRepository, passwordHasher ports.PasswordHasher) *UserSerivce {
	return &UserSerivce{
		repo:           repo,
		passwordHasher: passwordHasher,
	}
}

func (s *UserSerivce) CreateUser(ctx context.Context, user domain.User) (int, error) {
	hashedPassword, err := s.passwordHasher.HashPassword(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = hashedPassword

	id, err := s.repo.SaveUser(ctx, user)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *UserSerivce) GetUserById(ctx context.Context, id int) (domain.User, error) {
	user, err := s.repo.FindUserById(ctx, id)
	if err != nil {
		return domain.User{}, errors.New("failed to get user by id: " + err.Error())
	}
	return user, nil
}
