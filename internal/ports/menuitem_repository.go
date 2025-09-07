package ports

import (
  "context"
  "github.com/mohits-git/food-ordering-system/internal/domain"
)

type MenuItemRepository interface {
  SaveMenuItem(cxt context.Context, item domain.MenuItem) (int, error)
  UpdateMenuItemAvailability(cxt context.Context, id int, available bool) error
  FindMenuItemsByRestaurantId(cxt context.Context, restaurantId int) ([]domain.MenuItem, error)
  FindMenuItemById(cxt context.Context, id int) (domain.MenuItem, error)
}
