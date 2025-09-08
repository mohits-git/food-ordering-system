package services

import (
	"testing"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/utils/authctx"
	mockrepository "github.com/mohits-git/food-ordering-system/tests/mock_repository"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_services_OrderService_NewOrderService(t *testing.T) {
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewOrderService(&mockOrderRepo, &mockMenuItemRepo)
	require.NotNil(t, service)
}

func Test_services_OrderService_getRestaurantItemsMap(t *testing.T) {
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewOrderService(&mockOrderRepo, &mockMenuItemRepo)

	mockMenuItemRepo.On("FindMenuItemsByRestaurantId", mock.Anything, 1).
		Return([]domain.MenuItem{
			{ID: 1, Name: "Item 1", Price: 100, Available: true, RestaurantID: 1},
			{ID: 2, Name: "Item 2", Price: 200, Available: false, RestaurantID: 1},
		}, nil)

	itemsMap, err := service.getRestaurantItemsMap(t.Context(), 1)
	require.NoError(t, err)
	require.Len(t, itemsMap, 2)
	require.True(t, itemsMap[1])
	require.False(t, itemsMap[2])
	mockMenuItemRepo.AssertExpectations(t)
}

func Test_services_OrderService_getRestaurantItemsMap_when_no_items(t *testing.T) {
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewOrderService(&mockOrderRepo, &mockMenuItemRepo)

	mockMenuItemRepo.On("FindMenuItemsByRestaurantId", mock.Anything, 1).
		Return([]domain.MenuItem{}, nil)

	itemsMap, err := service.getRestaurantItemsMap(t.Context(), 1)
	require.Error(t, err)
	require.Nil(t, itemsMap)
	mockMenuItemRepo.AssertExpectations(t)
}

func Test_services_OrderService_CreateOrder(t *testing.T) {
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewOrderService(&mockOrderRepo, &mockMenuItemRepo)

	order := domain.Order{
		CustomerID:   1,
		RestaurantID: 1,
		OrderItems: []domain.OrderItem{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	}

	authCtx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.CUSTOMER,
	})

	mockMenuItemRepo.On("FindMenuItemsByRestaurantId", mock.Anything, 1).
		Return([]domain.MenuItem{
			{ID: 1, Name: "Item 1", Price: 100, Available: true, RestaurantID: 1},
			{ID: 2, Name: "Item 2", Price: 200, Available: true, RestaurantID: 1},
		}, nil)

	mockOrderRepo.On("SaveOrder", mock.Anything, order).
		Return(1, nil)

	id, err := service.CreateOrder(authCtx, order)
	require.NoError(t, err)
	require.Equal(t, 1, id)
	mockMenuItemRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}

func Test_services_OrderService_CreateOrder_when_invalid(t *testing.T) {
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewOrderService(&mockOrderRepo, &mockMenuItemRepo)

	order := domain.Order{
		CustomerID:   1,
		RestaurantID: 1,
		OrderItems: []domain.OrderItem{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	}

	authCtx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.CUSTOMER,
	})

	mockMenuItemRepo.On("FindMenuItemsByRestaurantId", mock.Anything, 1).
		Return([]domain.MenuItem{
			{ID: 1, Name: "Item 1", Price: 100, Available: true, RestaurantID: 1},
			{ID: 2, Name: "Item 2", Price: 200, Available: false, RestaurantID: 1},
		}, nil)

	id, err := service.CreateOrder(authCtx, order)
	require.Error(t, err)
	require.Equal(t, 0, id)
	mockMenuItemRepo.AssertExpectations(t)
	mockOrderRepo.AssertNotCalled(t, "SaveOrder", mock.Anything, mock.Anything)
}

func Test_services_OrderService_CreateOrder_when_unauthorized(t *testing.T) {
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewOrderService(&mockOrderRepo, &mockMenuItemRepo)

	order := domain.Order{
		CustomerID:   1,
		RestaurantID: 1,
		OrderItems: []domain.OrderItem{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	}

	id, err := service.CreateOrder(t.Context(), order)
	require.Error(t, err)
	require.Equal(t, 0, id)
	mockMenuItemRepo.AssertNotCalled(t, "FindMenuItemsByRestaurantId", mock.Anything, mock.Anything)
	mockOrderRepo.AssertNotCalled(t, "SaveOrder", mock.Anything, mock.Anything)
}

func Test_services_OrderService_CreateOrder_when_forbidden(t *testing.T) {
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewOrderService(&mockOrderRepo, &mockMenuItemRepo)

	order := domain.Order{
		CustomerID:   1,
		RestaurantID: 1,
		OrderItems: []domain.OrderItem{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	}

	authCtx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 2,
		Role:   domain.CUSTOMER,
	})

	id, err := service.CreateOrder(authCtx, order)
	require.Error(t, err)
	require.Equal(t, 0, id)
	mockMenuItemRepo.AssertNotCalled(t, "FindMenuItemsByRestaurantId", mock.Anything, mock.Anything)
	mockOrderRepo.AssertNotCalled(t, "SaveOrder", mock.Anything, mock.Anything)
}

func Test_services_OrderService_GetOrderById(t *testing.T) {
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewOrderService(&mockOrderRepo, &mockMenuItemRepo)

	order := domain.Order{
		ID:           1,
		CustomerID:   1,
		RestaurantID: 1,
		OrderItems: []domain.OrderItem{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	}

	authCtx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.CUSTOMER,
	})

	mockOrderRepo.On("FindOrderById", mock.Anything, 1).
		Return(order, nil)

	fetchedOrder, err := service.GetOrderById(authCtx, 1)
	require.NoError(t, err)
	require.Equal(t, order, fetchedOrder)
	mockOrderRepo.AssertExpectations(t)
}

func Test_services_OrderService_GetOrderById_when_invalid_id(t *testing.T) {
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewOrderService(&mockOrderRepo, &mockMenuItemRepo)

	authCtx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.CUSTOMER,
	})

	fetchedOrder, err := service.GetOrderById(authCtx, 0)
	require.Error(t, err)
	require.Equal(t, domain.Order{}, fetchedOrder)
	mockOrderRepo.AssertNotCalled(t, "FindOrderById", mock.Anything, mock.Anything)
}

func Test_services_OrderService_GetOrderById_when_unauthorized(t *testing.T) {
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewOrderService(&mockOrderRepo, &mockMenuItemRepo)

	fetchedOrder, err := service.GetOrderById(t.Context(), 1)
	require.Error(t, err)
	require.Equal(t, domain.Order{}, fetchedOrder)
	mockOrderRepo.AssertNotCalled(t, "FindOrderById", mock.Anything, mock.Anything)
}

func Test_services_OrderService_GetOrderById_when_forbidden(t *testing.T) {
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewOrderService(&mockOrderRepo, &mockMenuItemRepo)

	order := domain.Order{
		ID:           1,
		CustomerID:   1,
		RestaurantID: 1,
		OrderItems: []domain.OrderItem{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	}

	authCtx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 2,
		Role:   domain.CUSTOMER,
	})

	mockOrderRepo.On("FindOrderById", mock.Anything, 1).
		Return(order, nil)

	fetchedOrder, err := service.GetOrderById(authCtx, 1)
	require.Error(t, err)
	require.Equal(t, domain.Order{}, fetchedOrder)
	mockOrderRepo.AssertExpectations(t)
}

func Test_services_OrderService_AddOrderItem(t *testing.T) {
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewOrderService(&mockOrderRepo, &mockMenuItemRepo)

	order := domain.Order{
		ID:           1,
		CustomerID:   1,
		RestaurantID: 1,
		OrderItems: []domain.OrderItem{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	}

	newItem := domain.OrderItem{MenuItemID: 3, Quantity: 1}

	authCtx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.CUSTOMER,
	})

	mockOrderRepo.On("FindOrderById", mock.Anything, 1).
		Return(order, nil)

	mockMenuItemRepo.On("FindMenuItemById", mock.Anything, 3).
		Return(domain.MenuItem{ID: 3, Name: "Item 3", Price: 150, Available: true, RestaurantID: 1}, nil)

	mockOrderRepo.On("UpdateOrder", mock.Anything, domain.Order{
		ID:           1,
		CustomerID:   1,
		RestaurantID: 1,
		OrderItems: []domain.OrderItem{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
			{MenuItemID: 3, Quantity: 1},
		},
	}).Return(nil)

	err := service.AddOrderItem(authCtx, 1, newItem)
	require.NoError(t, err)
	mockOrderRepo.AssertExpectations(t)
	mockMenuItemRepo.AssertExpectations(t)
}

func Test_services_OrderService_AddOrderItem_when_invalid(t *testing.T) {
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewOrderService(&mockOrderRepo, &mockMenuItemRepo)

	newItem := domain.OrderItem{MenuItemID: 3, Quantity: 1}

	authCtx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.CUSTOMER,
	})

	err := service.AddOrderItem(authCtx, 0, newItem)
	require.Error(t, err)
	mockOrderRepo.AssertNotCalled(t, "FindOrderById", mock.Anything, mock.Anything)
	mockMenuItemRepo.AssertNotCalled(t, "FindMenuItemById", mock.Anything, mock.Anything)
	mockOrderRepo.AssertNotCalled(t, "UpdateOrder", mock.Anything, mock.Anything)
}

func Test_services_OrderService_AddOrderItem_when_unauthorized(t *testing.T) {
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewOrderService(&mockOrderRepo, &mockMenuItemRepo)

	newItem := domain.OrderItem{MenuItemID: 3, Quantity: 1}

	err := service.AddOrderItem(t.Context(), 1, newItem)
	require.Error(t, err)
	mockOrderRepo.AssertNotCalled(t, "FindOrderById", mock.Anything, mock.Anything)
	mockMenuItemRepo.AssertNotCalled(t, "FindMenuItemById", mock.Anything, mock.Anything)
	mockOrderRepo.AssertNotCalled(t, "UpdateOrder", mock.Anything, mock.Anything)
}

func Test_services_OrderService_AddOrderItem_when_forbidden(t *testing.T) {
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewOrderService(&mockOrderRepo, &mockMenuItemRepo)

	order := domain.Order{
		ID:           1,
		CustomerID:   1,
		RestaurantID: 1,
		OrderItems: []domain.OrderItem{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	}

	newItem := domain.OrderItem{MenuItemID: 3, Quantity: 1}

	authCtx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 2,
		Role:   domain.CUSTOMER,
	})

	mockOrderRepo.On("FindOrderById", mock.Anything, 1).
		Return(order, nil)

	err := service.AddOrderItem(authCtx, 1, newItem)
	require.Error(t, err)
	mockOrderRepo.AssertExpectations(t)
	mockMenuItemRepo.AssertNotCalled(t, "FindMenuItemById", mock.Anything, mock.Anything)
	mockOrderRepo.AssertNotCalled(t, "UpdateOrder", mock.Anything, mock.Anything)
}

func Test_services_OrderService_AddOrderItem_when_item_not_belong_to_restaurant(t *testing.T) {
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewOrderService(&mockOrderRepo, &mockMenuItemRepo)

	order := domain.Order{
		ID:           1,
		CustomerID:   1,
		RestaurantID: 1,
		OrderItems: []domain.OrderItem{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	}

	newItem := domain.OrderItem{MenuItemID: 3, Quantity: 1}

	authCtx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.CUSTOMER,
	})

	mockOrderRepo.On("FindOrderById", mock.Anything, 1).
		Return(order, nil)

	mockMenuItemRepo.On("FindMenuItemById", mock.Anything, 3).
		Return(domain.MenuItem{ID: 3, Name: "Item 3", Price: 150, Available: true, RestaurantID: 2}, nil)

	err := service.AddOrderItem(authCtx, 1, newItem)
	require.Error(t, err)
	mockOrderRepo.AssertExpectations(t)
	mockMenuItemRepo.AssertExpectations(t)
	mockOrderRepo.AssertNotCalled(t, "UpdateOrder", mock.Anything, mock.Anything)
}

func Test_services_OrderService_AddOrderItem_when_item_not_available(t *testing.T) {
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewOrderService(&mockOrderRepo, &mockMenuItemRepo)

	order := domain.Order{
		ID:           1,
		CustomerID:   1,
		RestaurantID: 1,
		OrderItems: []domain.OrderItem{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	}

	newItem := domain.OrderItem{MenuItemID: 3, Quantity: 1}

	authCtx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.CUSTOMER,
	})

	mockOrderRepo.On("FindOrderById", mock.Anything, 1).
		Return(order, nil)

	mockMenuItemRepo.On("FindMenuItemById", mock.Anything, 3).
		Return(domain.MenuItem{ID: 3, Name: "Item 3", Price: 150, Available: false, RestaurantID: 1}, nil)

	err := service.AddOrderItem(authCtx, 1, newItem)
	require.Error(t, err)
	mockOrderRepo.AssertExpectations(t)
	mockMenuItemRepo.AssertExpectations(t)
	mockOrderRepo.AssertNotCalled(t, "UpdateOrder", mock.Anything, mock.Anything)
}
