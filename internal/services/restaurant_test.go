package services

import (
	"testing"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
	"github.com/mohits-git/food-ordering-system/internal/utils/authctx"
	mockrepository "github.com/mohits-git/food-ordering-system/tests/mock_repository"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_services_RestaurantService_NewRestaurantService(t *testing.T) {
	mockRepo := mockrepository.RestaurantRepository{}
	service := NewRestaurantService(&mockRepo)
	require.NotNil(t, service)
}

func Test_services_RestaurantService_GetAllRestaurants(t *testing.T) {
	mockRepo := mockrepository.RestaurantRepository{}
	service := NewRestaurantService(&mockRepo)

	mockRepo.On("FindAllRestaurants", mock.Anything).
		Return([]domain.Restaurant{
			{ID: 1, Name: "Restaurant 1", OwnerID: 1},
			{ID: 2, Name: "Restaurant 2", OwnerID: 2},
		}, nil)

	restaurants, err := service.GetAllRestaurants(t.Context())
	require.NoError(t, err)
	require.Len(t, restaurants, 2)
	require.Equal(t, "Restaurant 1", restaurants[0].Name)
	require.Equal(t, "Restaurant 2", restaurants[1].Name)
	mockRepo.AssertExpectations(t)
}

func Test_services_RestaurantService_GetAllRestaurants_when_error(t *testing.T) {
	mockRepo := mockrepository.RestaurantRepository{}
	service := NewRestaurantService(&mockRepo)
	expectedErr := apperr.NewAppError(apperr.ErrInternal, "internal error", nil)

	mockRepo.On("FindAllRestaurants", mock.Anything).
		Return([]domain.Restaurant{}, expectedErr)

	restaurants, err := service.GetAllRestaurants(t.Context())
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, restaurants)
	mockRepo.AssertExpectations(t)
}

func Test_services_RestaurantService_CreateRestaurant(t *testing.T) {
	mockRepo := mockrepository.RestaurantRepository{}
	service := NewRestaurantService(&mockRepo)

	newRestaurant := domain.Restaurant{Name: "New Restaurant", OwnerID: 1}

	mockRepo.On("SaveRestaurant", mock.Anything, newRestaurant).
		Return(1, nil)

	ctx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.OWNER,
	})

	restaurantId, err := service.CreateRestaurant(ctx, newRestaurant.Name)
	require.NoError(t, err)
	require.Equal(t, 1, restaurantId)
	mockRepo.AssertExpectations(t)
}

func Test_services_RestaurantService_CreateRestaurant_when_invalid_name(t *testing.T) {
	mockRepo := mockrepository.RestaurantRepository{}
	service := NewRestaurantService(&mockRepo)

	ctx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.OWNER,
	})

	restaurantId, err := service.CreateRestaurant(ctx, "")
	require.Error(t, err)
	require.Equal(t, 0, restaurantId)
	mockRepo.AssertNotCalled(t, "SaveRestaurant", mock.Anything, mock.Anything)
}

func Test_services_RestaurantService_CreateRestaurant_when_unauthorized(t *testing.T) {
	mockRepo := mockrepository.RestaurantRepository{}
	service := NewRestaurantService(&mockRepo)

	restaurantId, err := service.CreateRestaurant(t.Context(), "New Restaurant")
	require.Error(t, err)
	require.Equal(t, 0, restaurantId)
	mockRepo.AssertNotCalled(t, "SaveRestaurant", mock.Anything, mock.Anything)
}

func Test_services_RestaurantService_CreateRestaurant_when_forbidden(t *testing.T) {
	mockRepo := mockrepository.RestaurantRepository{}
	service := NewRestaurantService(&mockRepo)

	ctx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.CUSTOMER,
	})

	restaurantId, err := service.CreateRestaurant(ctx, "New Restaurant")
	require.Error(t, err)
	require.Equal(t, 0, restaurantId)
	mockRepo.AssertNotCalled(t, "SaveRestaurant", mock.Anything, mock.Anything)
}

func Test_services_RestaurantService_CreateRestaurant_when_repo_error(t *testing.T) {
	mockRepo := mockrepository.RestaurantRepository{}
	service := NewRestaurantService(&mockRepo)
	expectedErr := apperr.NewAppError(apperr.ErrInternal, "internal error", nil)

	newRestaurant := domain.Restaurant{Name: "New Restaurant", OwnerID: 1}

	mockRepo.On("SaveRestaurant", mock.Anything, newRestaurant).
		Return(0, expectedErr)

	ctx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.OWNER,
	})

	restaurantId, err := service.CreateRestaurant(ctx, newRestaurant.Name)
	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, 0, restaurantId)
	mockRepo.AssertExpectations(t)
}
