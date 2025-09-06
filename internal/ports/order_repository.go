package ports

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
)

type OrderRepository interface {
  SaveOrder(ctx context.Context, order domain.Order) (int, error)
  FindOrderById(ctx context.Context, id int) (domain.Order, error)
  UpdateOrder(ctx context.Context, order domain.Order) error
}
