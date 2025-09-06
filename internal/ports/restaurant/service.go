package ports

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
)

type RestaurantService interface {
  CreateRestaurant(ctx context.Context, menuItem domain.MenuItem) error
  GetRestaurantById(ctx context.Context, id string) (domain.Restaurant, error)
  GetAllRestaurants(ctx context.Context) ([]domain.Restaurant, error)
}
