package dtos

import "github.com/mohits-git/food-ordering-system/internal/domain"

type AddMenuItemRequest struct {
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Available bool    `json:"available"`
	ImageURL  string  `json:"image_url"`
}

type UpdateMenuItemAvailabilityRequest struct {
	Available bool `json:"available"`
}

type AddMenuItemResponse struct {
	ID int `json:"id"`
}

type UpdateMenuItemResponse struct{}

type MenuItemResponse struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Available bool    `json:"available"`
	ImageURL  string  `json:"image_url"`
}

func NewMenuItemResponse(item domain.MenuItem) MenuItemResponse {
	return MenuItemResponse{
		ID:        item.ID,
		Name:      item.Name,
		Price:     item.Price,
		Available: item.Available,
		ImageURL:  item.ImageURL,
	}
}

type GetMenuItemsResponse struct {
	RestaurantID int                `json:"restaurant_id"`
	Items        []MenuItemResponse `json:"items"`
}
