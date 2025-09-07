package services

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/ports"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
	"github.com/mohits-git/food-ordering-system/internal/utils/authctx"
)

type RestaurantService struct {
	restaurantRepo ports.RestaurantRepository
}

func NewRestaurantService(restaurantRepo ports.RestaurantRepository) *RestaurantService {
	return &RestaurantService{
		restaurantRepo: restaurantRepo,
	}
}

func (s *RestaurantService) CreateRestaurant(ctx context.Context, restaurantName string) (int, error) {
	if restaurantName == "" {
		return 0, apperr.NewAppError(apperr.ErrInvalid, "restaurant name cannot be empty", nil)
	}

	user, ok := authctx.UserClaimsFromCtx(ctx)
	if !ok {
		return 0, apperr.NewAppError(apperr.ErrUnauthorized, "unauthorized", nil)
	}
	if user.Role != domain.OWNER {
		return 0, apperr.NewAppError(apperr.ErrForbidden, "forbidden", nil)
	}

	id, err := s.restaurantRepo.SaveRestaurant(ctx,
		domain.NewRestaurant(0, restaurantName, user.UserID))
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *RestaurantService) GetAllRestaurants(ctx context.Context) ([]domain.Restaurant, error) {
	restaurants, err := s.restaurantRepo.FindAllRestaurants(ctx)
	if err != nil {
		return nil, err
	}
	return restaurants, nil
}
