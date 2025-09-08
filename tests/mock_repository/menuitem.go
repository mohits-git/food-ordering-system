package mockrepository

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MenuItemRepository struct {
	mock.Mock
}

func (m *MenuItemRepository) SaveMenuItem(cxt context.Context, item domain.MenuItem) (int, error) {
	args := m.Called(cxt, item)
	return args.Int(0), args.Error(1)
}

func (m *MenuItemRepository) UpdateMenuItemAvailability(cxt context.Context, id int, available bool) error {
	args := m.Called(cxt, id, available)
	return args.Error(0)
}

func (m *MenuItemRepository) FindMenuItemsByRestaurantId(cxt context.Context, restaurantId int) ([]domain.MenuItem, error) {
	args := m.Called(cxt, restaurantId)
	return args.Get(0).([]domain.MenuItem), args.Error(1)
}

func (m *MenuItemRepository) FindMenuItemById(cxt context.Context, id int) (domain.MenuItem, error) {
	args := m.Called(cxt, id)
	return args.Get(0).(domain.MenuItem), args.Error(1)
}
