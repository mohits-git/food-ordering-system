package mockrepository

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (u *UserRepository) FindUserById(ctx context.Context, id int) (domain.User, error) {
	args := u.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (u *UserRepository) FindUserByEmail(ctx context.Context, email string) (domain.User, error) {
	args := u.Called(ctx, email)
	return args.Get(0).(domain.User), args.Error(1)
}

func (u *UserRepository) SaveUser(ctx context.Context, user domain.User) (int, error) {
	args := u.Called(ctx, user)
	return args.Int(0), args.Error(1)
}
