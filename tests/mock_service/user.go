package mockservice

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/stretchr/testify/mock"
)

type UserService struct {
	mock.Mock
}

func (s *UserService) GetUserById(ctx context.Context, id int) (domain.User, error) {
	args := s.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (s *UserService) CreateUser(ctx context.Context, user domain.User) (int, error) {
	args := s.Called(ctx, user)
	return args.Int(0), args.Error(1)
}
