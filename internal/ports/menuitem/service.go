package ports

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
)

type MenuItemService interface {
  CreateMenuItemForRestaurant(cxt context.Context, restaurantId string, item domain.MenuItem) error
  UpdateMenuItem(cxt context.Context, item domain.MenuItem) error
  GetMenuItemById(cxt context.Context, id string) (domain.MenuItem, error)
  GetAllMenuItemsByRestaurantId(cxt context.Context, restaurantId string) ([]domain.MenuItem, error)
  GetAvailableMenuItemsByRestaurantId(cxt context.Context, restaurantId string) ([]domain.MenuItem, error)
}
