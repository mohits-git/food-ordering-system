package mockrepository

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/stretchr/testify/mock"
)

type RestaurantRepository struct {
	mock.Mock
}

func (r *RestaurantRepository) SaveRestaurant(cxt context.Context, restaurant domain.Restaurant) (int, error) {
	args := r.Called(cxt, restaurant)
	return args.Int(0), args.Error(1)
}

func (r *RestaurantRepository) FindAllRestaurants(cxt context.Context) ([]domain.Restaurant, error) {
	args := r.Called(cxt)
	return args.Get(0).([]domain.Restaurant), args.Error(1)
}

func (r *RestaurantRepository) FindRestaurantById(cxt context.Context, id int) (domain.Restaurant, error) {
	args := r.Called(cxt, id)
	return args.Get(0).(domain.Restaurant), args.Error(1)
}
