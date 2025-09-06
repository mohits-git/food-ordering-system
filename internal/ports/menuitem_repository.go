package ports

import (
  "context"
  "github.com/mohits-git/food-ordering-system/internal/domain"
)

type MenuItemRepository interface {
  SaveMenuItem(cxt context.Context, item domain.MenuItem) (int, error)
  UpdateMenuItem(cxt context.Context, item domain.MenuItem) error
  FindMenuItemById(cxt context.Context, id int) (domain.MenuItem, error)
  FindMenuItemsByRestaurantId(cxt context.Context, restaurantId string) ([]domain.MenuItem, error)
  FindAvailableMenuItemsByRestaurantId(cxt context.Context, restaurantId string) ([]domain.MenuItem, error)
}
