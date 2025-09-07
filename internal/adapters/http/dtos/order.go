package dtos

import "github.com/mohits-git/food-ordering-system/internal/domain"

type CreateOrderRequest struct {
	RestaurantID int             `json:"restaurant_id"`
	OrderItems   []OrderItemsDTO `json:"order_items"`
}

type OrderItemsDTO struct {
	MenuItemID int `json:"menu_item_id"`
	Quantity   int `json:"quantity"`
}

func (o *OrderItemsDTO) ToDomain() domain.OrderItem {
	return domain.OrderItem{
		MenuItemID: o.MenuItemID,
		Quantity:   o.Quantity,
	}
}

type AddOrderItemRequest struct {
	MenuItemID int `json:"menu_item_id"`
	Quantity   int `json:"quantity"`
}

type CreateOrderResponse struct {
	ID int `json:"id"`
}

type AddOrderItemResponse struct {
	ID int `json:"id"`
}

type GetOrderByIdResponse struct {
	ID           int             `json:"id"`
	CustomerID   int             `json:"customer_id"`
	RestaurantID int             `json:"restaurant_id"`
	OrderItems   []OrderItemsDTO `json:"order_items"`
}
