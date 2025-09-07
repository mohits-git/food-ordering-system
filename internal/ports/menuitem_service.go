package ports

import (
  "context"

  "github.com/mohits-git/food-ordering-system/internal/domain"
)

type MenuItemService interface {
  CreateMenuItemForRestaurant(ctx context.Context, item domain.MenuItem) (int, error)
  GetAllMenuItemsByRestaurantId(ctx context.Context, restaurantId int) ([]domain.MenuItem, error)
  UpdateAvailability(ctx context.Context, id int, available bool) error
}
