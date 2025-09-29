package ports

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
)

type RestaurantRepository interface {
	SaveRestaurant(cxt context.Context, restaurant domain.Restaurant) (int, error)
	FindAllRestaurants(cxt context.Context) ([]domain.Restaurant, error)
	FindRestaurantById(cxt context.Context, id int) (domain.Restaurant, error)
}
