package ports

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
)

type UserRepository interface {
	FindUserById(ctx context.Context, id int) (domain.User, error)
	FindUserByEmail(ctx context.Context, email string) (domain.User, error)
	SaveUser(ctx context.Context, user domain.User) (int, error)
	// UpdateUser(domain.User) error
	// DeleteUser(id int) error
}
