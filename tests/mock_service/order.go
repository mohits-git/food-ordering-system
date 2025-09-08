package mockservice

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/stretchr/testify/mock"
)

type OrderService struct {
	mock.Mock
}

func (s *OrderService) CreateOrder(ctx context.Context, order domain.Order) (int, error) {
	args := s.Called(ctx, order)
	return args.Int(0), args.Error(1)
}

func (s *OrderService) GetOrderById(ctx context.Context, id int) (domain.Order, error) {
	args := s.Called(ctx, id)
	return args.Get(0).(domain.Order), args.Error(1)
}

func (s *OrderService) AddOrderItem(ctx context.Context, orderId int, item domain.OrderItem) error {
	args := s.Called(ctx, orderId, item)
	return args.Error(0)
}
