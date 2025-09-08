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

func Test_services_MenuItemService_NewMenuItemsService(t *testing.T) {
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	mockRestaurantRepo := mockrepository.RestaurantRepository{}
	service := NewMenuItemsService(&mockMenuItemRepo, &mockRestaurantRepo)
	require.NotNil(t, service)
}

func Test_services_MenuItemService_GetAllMenuItemsByRestaurantId(t *testing.T) {
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	mockRestaurantRepo := mockrepository.RestaurantRepository{}
	service := NewMenuItemsService(&mockMenuItemRepo, &mockRestaurantRepo)

	restaurantId := 1
	mockMenuItemRepo.On("FindMenuItemsByRestaurantId", mock.Anything, restaurantId).
		Return([]domain.MenuItem{
			{ID: 1, Name: "Item 1", Price: 10.0, Available: true, RestaurantID: restaurantId},
			{ID: 2, Name: "Item 2", Price: 15.0, Available: false, RestaurantID: restaurantId},
		}, nil)

	items, err := service.GetAllMenuItemsByRestaurantId(t.Context(), restaurantId)
	require.NoError(t, err)
	require.Len(t, items, 2)
	require.Equal(t, "Item 1", items[0].Name)
	require.Equal(t, "Item 2", items[1].Name)
	mockMenuItemRepo.AssertExpectations(t)
}

func Test_services_MenuItemService_GetAllMenuItemsByRestaurantId_when_error(t *testing.T) {
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	mockRestaurantRepo := mockrepository.RestaurantRepository{}
	service := NewMenuItemsService(&mockMenuItemRepo, &mockRestaurantRepo)
	expectedErr := apperr.NewAppError(apperr.ErrInternal, "internal error", nil)

	restaurantId := 1
	mockMenuItemRepo.On("FindMenuItemsByRestaurantId", mock.Anything, restaurantId).
		Return([]domain.MenuItem{}, expectedErr)

	items, err := service.GetAllMenuItemsByRestaurantId(t.Context(), restaurantId)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, items)
	mockMenuItemRepo.AssertExpectations(t)
}

func Test_services_MenuItemService_CreateMenuItemForRestaurant(t *testing.T) {
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	mockRestaurantRepo := mockrepository.RestaurantRepository{}
	service := NewMenuItemsService(&mockMenuItemRepo, &mockRestaurantRepo)

	newItem := domain.MenuItem{Name: "New Item", Price: 20.0, Available: true, RestaurantID: 1}
	restaurant := domain.Restaurant{ID: 1, Name: "Restaurant 1", OwnerID: 1}

	mockRestaurantRepo.On("FindRestaurantById", mock.Anything, newItem.RestaurantID).
		Return(restaurant, nil)
	mockMenuItemRepo.On("SaveMenuItem", mock.Anything, newItem).
		Return(1, nil)

	ctx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.OWNER,
	})

	itemId, err := service.CreateMenuItemForRestaurant(ctx, newItem)
	require.NoError(t, err)
	require.Equal(t, 1, itemId)
	mockMenuItemRepo.AssertExpectations(t)
	mockRestaurantRepo.AssertExpectations(t)
}

func Test_services_MenuItemService_CreateMenuItemForRestaurant_when_invalid_data(t *testing.T) {
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	mockRestaurantRepo := mockrepository.RestaurantRepository{}
	service := NewMenuItemsService(&mockMenuItemRepo, &mockRestaurantRepo)

	newItem := domain.MenuItem{Name: "", Price: -20.0, Available: true, RestaurantID: 1}

	ctx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.OWNER,
	})

	itemId, err := service.CreateMenuItemForRestaurant(ctx, newItem)
	require.Error(t, err)
	require.Equal(t, 0, itemId)
	mockMenuItemRepo.AssertExpectations(t)
	mockRestaurantRepo.AssertExpectations(t)
}

func Test_services_MenuItemService_CreateMenuItemForRestaurant_when_unauthenticated(t *testing.T) {
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	mockRestaurantRepo := mockrepository.RestaurantRepository{}
	service := NewMenuItemsService(&mockMenuItemRepo, &mockRestaurantRepo)

	newItem := domain.MenuItem{Name: "New Item", Price: 20.0, Available: true, RestaurantID: 1}

	itemId, err := service.CreateMenuItemForRestaurant(t.Context(), newItem)
	require.Error(t, err)
	require.Equal(t, 0, itemId)
	mockMenuItemRepo.AssertExpectations(t)
	mockRestaurantRepo.AssertExpectations(t)
}

func Test_services_MenuItemService_CreateMenuItemForRestaurant_when_forbidden(t *testing.T) {
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	mockRestaurantRepo := mockrepository.RestaurantRepository{}
	service := NewMenuItemsService(&mockMenuItemRepo, &mockRestaurantRepo)

	newItem := domain.MenuItem{Name: "New Item", Price: 20.0, Available: true, RestaurantID: 1}

	ctx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.CUSTOMER,
	})

	itemId, err := service.CreateMenuItemForRestaurant(ctx, newItem)
	require.Error(t, err)
	require.Equal(t, 0, itemId)
	mockMenuItemRepo.AssertExpectations(t)
	mockRestaurantRepo.AssertExpectations(t)
}

func Test_services_MenuItemService_CreateMenuItemForRestaurant_when_not_owner(t *testing.T) {
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	mockRestaurantRepo := mockrepository.RestaurantRepository{}
	service := NewMenuItemsService(&mockMenuItemRepo, &mockRestaurantRepo)

	newItem := domain.MenuItem{Name: "New Item", Price: 20.0, Available: true, RestaurantID: 1}
	restaurant := domain.Restaurant{ID: 1, Name: "Restaurant 1", OwnerID: 2}

	mockRestaurantRepo.On("FindRestaurantById", mock.Anything, newItem.RestaurantID).
		Return(restaurant, nil)

	ctx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.OWNER,
	})

	itemId, err := service.CreateMenuItemForRestaurant(ctx, newItem)
	require.Error(t, err)
	require.Equal(t, 0, itemId)
	mockMenuItemRepo.AssertExpectations(t)
	mockRestaurantRepo.AssertExpectations(t)
}

func Test_services_MenuItemService_CreateMenuItemForRestaurant_when_restaurant_not_found(t *testing.T) {
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	mockRestaurantRepo := mockrepository.RestaurantRepository{}
	service := NewMenuItemsService(&mockMenuItemRepo, &mockRestaurantRepo)
	expectedErr := apperr.NewAppError(apperr.ErrNotFound, "restaurant not found", nil)

	newItem := domain.MenuItem{Name: "New Item", Price: 20.0, Available: true, RestaurantID: 1}

	mockRestaurantRepo.On("FindRestaurantById", mock.Anything, newItem.RestaurantID).
		Return(domain.Restaurant{}, expectedErr)

	ctx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.OWNER,
	})

	itemId, err := service.CreateMenuItemForRestaurant(ctx, newItem)
	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, 0, itemId)
	mockMenuItemRepo.AssertExpectations(t)
	mockRestaurantRepo.AssertExpectations(t)
}

func Test_services_MenuItemService_CreateMenuItemForRestaurant_when_repo_error(t *testing.T) {
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	mockRestaurantRepo := mockrepository.RestaurantRepository{}
	service := NewMenuItemsService(&mockMenuItemRepo, &mockRestaurantRepo)
	expectedErr := apperr.NewAppError(apperr.ErrInternal, "internal error", nil)

	newItem := domain.MenuItem{Name: "New Item", Price: 20.0, Available: true, RestaurantID: 1}
	restaurant := domain.Restaurant{ID: 1, Name: "Restaurant 1", OwnerID: 1}

	mockRestaurantRepo.On("FindRestaurantById", mock.Anything, newItem.RestaurantID).
		Return(restaurant, nil)
	mockMenuItemRepo.On("SaveMenuItem", mock.Anything, newItem).
		Return(0, expectedErr)

	ctx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.OWNER,
	})

	itemId, err := service.CreateMenuItemForRestaurant(ctx, newItem)
	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, 0, itemId)
	mockMenuItemRepo.AssertExpectations(t)
	mockRestaurantRepo.AssertExpectations(t)
}

func Test_services_MenuItemService_UpdateAvailability(t *testing.T) {
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	mockRestaurantRepo := mockrepository.RestaurantRepository{}
	service := NewMenuItemsService(&mockMenuItemRepo, &mockRestaurantRepo)

	itemId := 1
	available := false
	menuItem := domain.MenuItem{ID: itemId, Name: "Item 1", Price: 10.0, Available: true, RestaurantID: 1}
	restaurant := domain.Restaurant{ID: 1, Name: "Restaurant 1", OwnerID: 1}

	mockMenuItemRepo.On("FindMenuItemById", mock.Anything, itemId).
		Return(menuItem, nil)
	mockRestaurantRepo.On("FindRestaurantById", mock.Anything, menuItem.RestaurantID).
		Return(restaurant, nil)
	mockMenuItemRepo.On("UpdateMenuItemAvailability", mock.Anything, itemId, available).
		Return(nil)

	ctx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.OWNER,
	})

	err := service.UpdateAvailability(ctx, itemId, available)
	require.NoError(t, err)
	mockMenuItemRepo.AssertExpectations(t)
	mockRestaurantRepo.AssertExpectations(t)
}

func Test_services_MenuItemService_UpdateAvailability_when_invalid_id(t *testing.T) {
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	mockRestaurantRepo := mockrepository.RestaurantRepository{}
	service := NewMenuItemsService(&mockMenuItemRepo, &mockRestaurantRepo)

	itemId := 0
	available := false

	ctx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.OWNER,
	})

	err := service.UpdateAvailability(ctx, itemId, available)
	require.Error(t, err)
	mockMenuItemRepo.AssertExpectations(t)
	mockRestaurantRepo.AssertExpectations(t)
}

func Test_services_MenuItemService_UpdateAvailability_when_unauthenticated(t *testing.T) {
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	mockRestaurantRepo := mockrepository.RestaurantRepository{}
	service := NewMenuItemsService(&mockMenuItemRepo, &mockRestaurantRepo)

	itemId := 1
	available := false

	err := service.UpdateAvailability(t.Context(), itemId, available)
	require.Error(t, err)
	mockMenuItemRepo.AssertExpectations(t)
	mockRestaurantRepo.AssertExpectations(t)
}

func Test_services_MenuItemService_UpdateAvailability_when_forbidden(t *testing.T) {
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	mockRestaurantRepo := mockrepository.RestaurantRepository{}
	service := NewMenuItemsService(&mockMenuItemRepo, &mockRestaurantRepo)

	itemId := 1
	available := false

	ctx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.CUSTOMER,
	})

	err := service.UpdateAvailability(ctx, itemId, available)
	require.Error(t, err)
	mockMenuItemRepo.AssertExpectations(t)
	mockRestaurantRepo.AssertExpectations(t)
}

func Test_services_MenuItemService_UpdateAvailability_when_not_owner(t *testing.T) {
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	mockRestaurantRepo := mockrepository.RestaurantRepository{}
	service := NewMenuItemsService(&mockMenuItemRepo, &mockRestaurantRepo)

	itemId := 1
	available := false
	menuItem := domain.MenuItem{ID: itemId, Name: "Item 1", Price: 10.0, Available: true, RestaurantID: 1}
	restaurant := domain.Restaurant{ID: 1, Name: "Restaurant 1", OwnerID: 2}

	mockMenuItemRepo.On("FindMenuItemById", mock.Anything, itemId).
		Return(menuItem, nil)
	mockRestaurantRepo.On("FindRestaurantById", mock.Anything, menuItem.RestaurantID).
		Return(restaurant, nil)

	ctx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.OWNER,
	})

	err := service.UpdateAvailability(ctx, itemId, available)
	require.Error(t, err)
	mockMenuItemRepo.AssertExpectations(t)
	mockRestaurantRepo.AssertExpectations(t)
}

func Test_services_MenuItemService_UpdateAvailability_when_menu_item_not_found(t *testing.T) {
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	mockRestaurantRepo := mockrepository.RestaurantRepository{}
	service := NewMenuItemsService(&mockMenuItemRepo, &mockRestaurantRepo)
	expectedErr := apperr.NewAppError(apperr.ErrNotFound, "menu item not found", nil)

	itemId := 1
	available := false

	mockMenuItemRepo.On("FindMenuItemById", mock.Anything, itemId).
		Return(domain.MenuItem{}, expectedErr)

	ctx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.OWNER,
	})

	err := service.UpdateAvailability(ctx, itemId, available)
	require.ErrorIs(t, err, expectedErr)
	mockMenuItemRepo.AssertExpectations(t)
	mockRestaurantRepo.AssertExpectations(t)
}

func Test_services_MenuItemService_UpdateAvailability_when_restaurant_not_found(t *testing.T) {
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	mockRestaurantRepo := mockrepository.RestaurantRepository{}
	service := NewMenuItemsService(&mockMenuItemRepo, &mockRestaurantRepo)
	expectedErr := apperr.NewAppError(apperr.ErrNotFound, "restaurant not found", nil)

	itemId := 1
	available := false
	menuItem := domain.MenuItem{ID: itemId, Name: "Item 1", Price: 10.0, Available: true, RestaurantID: 1}

	mockMenuItemRepo.On("FindMenuItemById", mock.Anything, itemId).
		Return(menuItem, nil)
	mockRestaurantRepo.On("FindRestaurantById", mock.Anything, menuItem.RestaurantID).
		Return(domain.Restaurant{}, expectedErr)

	ctx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.OWNER,
	})

	err := service.UpdateAvailability(ctx, itemId, available)
	require.ErrorIs(t, err, expectedErr)
	mockMenuItemRepo.AssertExpectations(t)
	mockRestaurantRepo.AssertExpectations(t)
}

func Test_services_MenuItemService_UpdateAvailability_when_repo_error(t *testing.T) {
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	mockRestaurantRepo := mockrepository.RestaurantRepository{}
	service := NewMenuItemsService(&mockMenuItemRepo, &mockRestaurantRepo)
	expectedErr := apperr.NewAppError(apperr.ErrInternal, "internal error", nil)

	itemId := 1
	available := false
	menuItem := domain.MenuItem{ID: itemId, Name: "Item 1", Price: 10.0, Available: true, RestaurantID: 1}
	restaurant := domain.Restaurant{ID: 1, Name: "Restaurant 1", OwnerID: 1}

	mockMenuItemRepo.On("FindMenuItemById", mock.Anything, itemId).
		Return(menuItem, nil)
	mockRestaurantRepo.On("FindRestaurantById", mock.Anything, menuItem.RestaurantID).
		Return(restaurant, nil)
	mockMenuItemRepo.On("UpdateMenuItemAvailability", mock.Anything, itemId, available).
		Return(expectedErr)

	ctx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.OWNER,
	})

	err := service.UpdateAvailability(ctx, itemId, available)
	require.ErrorIs(t, err, expectedErr)
	mockMenuItemRepo.AssertExpectations(t)
	mockRestaurantRepo.AssertExpectations(t)
}
