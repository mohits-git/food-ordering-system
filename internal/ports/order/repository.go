package ports

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
)

type OrderRepository interface {
  SaveOrder(ctx context.Context, order domain.Order) error
  FindOrderById(ctx context.Context, id string) (domain.Order, error)
  UpdateOrder(ctx context.Context, order domain.Order) error
}
