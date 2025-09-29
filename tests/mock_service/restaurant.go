package mockservice

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/stretchr/testify/mock"
)

type RestaurantService struct {
	mock.Mock
}

func (s *RestaurantService) CreateRestaurant(ctx context.Context, restaurantName, restaurantImage string) (int, error) {
	args := s.Called(ctx, restaurantName, restaurantImage)
	return args.Int(0), args.Error(1)
}

func (s *RestaurantService) GetAllRestaurants(ctx context.Context) ([]domain.Restaurant, error) {
	args := s.Called(ctx)
	return args.Get(0).([]domain.Restaurant), args.Error(1)
}
