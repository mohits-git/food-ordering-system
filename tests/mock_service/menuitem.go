package mockservice

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MenuItemService struct {
	mock.Mock
}

func (s *MenuItemService) CreateMenuItemForRestaurant(ctx context.Context, item domain.MenuItem) (int, error) {
	args := s.Called(ctx, item)
	return args.Int(0), args.Error(1)
}

func (s *MenuItemService) GetAllMenuItemsByRestaurantId(ctx context.Context, restaurantId int) ([]domain.MenuItem, error) {
	args := s.Called(ctx, restaurantId)
	return args.Get(0).([]domain.MenuItem), args.Error(1)
}

func (s *MenuItemService) UpdateAvailability(ctx context.Context, id int, available bool) error {
	args := s.Called(ctx, id, available)
	return args.Error(0)
}
