package mockrepository

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/stretchr/testify/mock"
)

type OrderRepository struct {
	mock.Mock
}

func (o *OrderRepository) SaveOrder(ctx context.Context, order domain.Order) (int, error) {
	args := o.Called(ctx, order)
	return args.Int(0), args.Error(1)
}

func (o *OrderRepository) FindOrderById(ctx context.Context, id int) (domain.Order, error) {
	args := o.Called(ctx, id)
	return args.Get(0).(domain.Order), args.Error(1)
}

func (o *OrderRepository) UpdateOrder(ctx context.Context, order domain.Order) error {
	args := o.Called(ctx, order)
	return args.Error(0)
}
