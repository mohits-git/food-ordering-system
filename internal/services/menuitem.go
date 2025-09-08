package services

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/ports"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
	"github.com/mohits-git/food-ordering-system/internal/utils/authctx"
)

type MenuItemService struct {
	menuItemRepo   ports.MenuItemRepository
	restaurantRepo ports.RestaurantRepository
}

func NewMenuItemsService(menuItemsRepo ports.MenuItemRepository, restaurantRepo ports.RestaurantRepository) *MenuItemService {
	return &MenuItemService{menuItemsRepo, restaurantRepo}
}

func (m *MenuItemService) CreateMenuItemForRestaurant(ctx context.Context, item domain.MenuItem) (int, error) {
	if !item.Validate() {
		return 0, apperr.NewAppError(apperr.ErrInvalid, "invalid menu item data", nil)
	}

	user, ok := authctx.UserClaimsFromCtx(ctx)
	if !ok {
		return 0, apperr.NewAppError(apperr.ErrUnauthorized, "unauthenticated user", nil)
	}
	if user.Role != domain.OWNER {
		return 0, apperr.NewAppError(apperr.ErrForbidden, "only restaurant owners can add menu items", nil)
	}

	restaurant, err := m.restaurantRepo.FindRestaurantById(ctx, item.RestaurantID)
	if err != nil {
		return 0, err
	}
	if restaurant.OwnerID != user.UserID {
		return 0, apperr.NewAppError(apperr.ErrForbidden, "only restaurant owners can add menu items", nil)
	}

	return m.menuItemRepo.SaveMenuItem(ctx, item)
}

func (m *MenuItemService) GetAllMenuItemsByRestaurantId(ctx context.Context, restaurantId int) ([]domain.MenuItem, error) {
	items, err := m.menuItemRepo.FindMenuItemsByRestaurantId(ctx, restaurantId)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (m *MenuItemService) UpdateAvailability(ctx context.Context, id int, available bool) error {
	if id <= 0 {
		return apperr.NewAppError(apperr.ErrInvalid, "invalid menu item id", nil)
	}

	user, ok := authctx.UserClaimsFromCtx(ctx)
	if !ok {
		return apperr.NewAppError(apperr.ErrUnauthorized, "unauthenticated user", nil)
	}
	if user.Role != domain.OWNER {
		return apperr.NewAppError(apperr.ErrForbidden, "only restaurant owners can update menu items", nil)
	}

	item, err := m.menuItemRepo.FindMenuItemById(ctx, id)
	if err != nil {
		return err
	}
	restaurant, err := m.restaurantRepo.FindRestaurantById(ctx, item.RestaurantID)
	if err != nil {
		return err
	}
	if restaurant.OwnerID != user.UserID {
		return apperr.NewAppError(apperr.ErrForbidden, "only restaurant owners can update menu items", nil)
	}

	return m.menuItemRepo.UpdateMenuItemAvailability(ctx, id, available)
}
