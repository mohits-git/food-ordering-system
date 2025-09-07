package ports

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
)

type OrderService interface {
  CreateOrder(ctx context.Context, order domain.Order) (int, error)
  GetOrderById(ctx context.Context, id int) (domain.Order, error)
  AddOrderItem(ctx context.Context, orderId int, item domain.OrderItem) error
}
