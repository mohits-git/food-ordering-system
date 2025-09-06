package ports

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
)

type OrderService interface {
  CreateOrder(ctx context.Context, order domain.Order) error
  GetOrderById(ctx context.Context, id string) (domain.Order, error)
  AddOrderItem(ctx context.Context, orderId string, item domain.OrderItem) error
  RemoveOrderItem(ctx context.Context, orderId string, itemId string) error
}
