package dtos

import "github.com/mohits-git/food-ordering-system/internal/domain"

type CreateRestaurantRequest struct {
	Name string `json:"name"`
}

type CreateRestaurantResponse struct {
	ID int `json:"id"`
}

type RestaurantDTO struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	OwnerID int    `json:"owner_id"`
}

func NewRestaurantDTO(restaurant domain.Restaurant) RestaurantDTO {
	return RestaurantDTO{
		ID:      restaurant.ID,
		Name:    restaurant.Name,
		OwnerID: restaurant.OwnerID,
	}
}

type GetRestaurantsResponse struct {
	Restaurants []RestaurantDTO `json:"restaurants"`
}

func NewRestaurant(restaurant RestaurantDTO) domain.Restaurant {
	return domain.Restaurant{
		ID:      restaurant.ID,
		Name:    restaurant.Name,
		OwnerID: restaurant.OwnerID,
	}
}
