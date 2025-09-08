package services

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/ports"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
	"github.com/mohits-git/food-ordering-system/internal/utils/authctx"
)

type OrderService struct {
	orderRepo    ports.OrderRepository
	menuItemRepo ports.MenuItemRepository
}

func NewOrderService(orderRepo ports.OrderRepository, menuItemRepo ports.MenuItemRepository) *OrderService {
	return &OrderService{orderRepo, menuItemRepo}
}

func (s *OrderService) getRestaurantItemsMap(ctx context.Context, restaurantId int) (map[int]bool, error) {
	restaurantItems, err := s.menuItemRepo.FindMenuItemsByRestaurantId(ctx, restaurantId)
	if err != nil {
		return nil, err
	}
	if len(restaurantItems) == 0 {
		return nil, apperr.NewAppError(apperr.ErrInvalid, "restaurant has no menu items", nil)
	}
	restaurantItemMap := make(map[int]bool)
	for _, item := range restaurantItems {
		restaurantItemMap[item.ID] = item.Available
	}
	return restaurantItemMap, nil
}

func (s *OrderService) CreateOrder(ctx context.Context, order domain.Order) (int, error) {
	user, ok := authctx.UserClaimsFromCtx(ctx)
	if !ok {
		return 0, apperr.NewAppError(apperr.ErrUnauthorized, "user not authenticated", nil)
	}
	if user.Role != domain.CUSTOMER || user.UserID != order.CustomerID {
		return 0, apperr.NewAppError(apperr.ErrForbidden, "only customers can create orders", nil)
	}

	restaurantItemsMap, err := s.getRestaurantItemsMap(ctx, order.RestaurantID)
	if err != nil {
		return 0, err
	}

	if ok := order.Validate(restaurantItemsMap); !ok {
		return 0, apperr.NewAppError(apperr.ErrInvalid, "invalid order data", nil)
	}

	id, err := s.orderRepo.SaveOrder(ctx, order)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *OrderService) GetOrderById(ctx context.Context, id int) (domain.Order, error) {
	if id <= 0 {
		return domain.Order{}, apperr.NewAppError(apperr.ErrInvalid, "invalid order id", nil)
	}

	user, ok := authctx.UserClaimsFromCtx(ctx)
	if !ok {
		return domain.Order{}, apperr.NewAppError(apperr.ErrUnauthorized, "user not authenticated", nil)
	}

	if user.Role != domain.CUSTOMER {
		return domain.Order{}, apperr.NewAppError(apperr.ErrForbidden, "only customers can access orders", nil)
	}

	order, err := s.orderRepo.FindOrderById(ctx, id)
	if err != nil {
		return domain.Order{}, err
	}

	if order.CustomerID != user.UserID {
		return domain.Order{}, apperr.NewAppError(apperr.ErrForbidden, "access to the order is forbidden", nil)
	}

	return order, nil
}

func (s *OrderService) AddOrderItem(ctx context.Context, orderId int, item domain.OrderItem) error {
	if orderId <= 0 || !item.Validate() {
		return apperr.NewAppError(apperr.ErrInvalid, "invalid input data", nil)
	}

	// authorize access to orders resource
	user, ok := authctx.UserClaimsFromCtx(ctx)
	if !ok {
		return apperr.NewAppError(apperr.ErrUnauthorized, "user not authenticated", nil)
	}
	if user.Role != domain.CUSTOMER {
		return apperr.NewAppError(apperr.ErrForbidden, "only customers can modify orders", nil)
	}

	// authorize access to order
	order, err := s.orderRepo.FindOrderById(ctx, orderId)
	if err != nil {
		return err
	}
	if order.CustomerID != user.UserID {
		return apperr.NewAppError(apperr.ErrForbidden, "access to the order is forbidden", nil)
	}

	// validation - check if item belongs to the same restaurant
	menuItem, err := s.menuItemRepo.FindMenuItemById(ctx, item.MenuItemID)
	if err != nil {
		return err
	}
	if menuItem.ID == 0 || menuItem.RestaurantID != order.RestaurantID {
		return apperr.NewAppError(apperr.ErrInvalid, "menu item does not belong to the restaurant of the order", nil)
	}
	if !menuItem.Available {
		return apperr.NewAppError(apperr.ErrInvalid, "menu item is not available", nil)
	}

	// save order
	order.AddItem(item.MenuItemID, item.Quantity)
	if err := s.orderRepo.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}
