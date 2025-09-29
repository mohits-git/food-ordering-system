package ports

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
)

type RestaurantService interface {
	CreateRestaurant(ctx context.Context, restaurantName string, restaurantImage string) (int, error)
	GetAllRestaurants(ctx context.Context) ([]domain.Restaurant, error)
}
