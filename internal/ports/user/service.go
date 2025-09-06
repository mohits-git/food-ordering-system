package ports

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
)

type UserService interface {
  GetUserById(ctx context.Context, id string) (domain.User, error)
  CreateUser(ctx context.Context, user domain.User) error
}
