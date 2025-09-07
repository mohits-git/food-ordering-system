package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/ports"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
)

type UserSerivce struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserSerivce {
	return &UserSerivce{repo: repo}
}

func (s *UserSerivce) CreateUser(ctx context.Context, user domain.User) (int, error) {
	exist, err := s.repo.FindUserByEmail(ctx, user.Email)
	if err == nil && exist.ID != 0 {
		return 0, apperr.NewAppError(apperr.ErrConflict, fmt.Sprintf("user with email %s already exists", user.Email), err)
	}
	id, err := s.repo.SaveUser(ctx, user)
	if err != nil {
		return 0, apperr.NewAppError(apperr.ErrInternal, "failed to create user", err)
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
