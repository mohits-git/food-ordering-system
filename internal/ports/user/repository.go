package ports

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
)

type UserRepository interface {
  FindUserById(ctx context.Context, id string) (domain.User, error)
  SaveUser(ctx context.Context, user domain.User) error
  // UpdateUser(domain.User) error
  // DeleteUser(id string) error
}
