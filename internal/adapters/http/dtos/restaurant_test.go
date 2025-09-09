package dtos

import (
	"testing"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/stretchr/testify/require"
)

func Test_dtos_NewRestaurantDTO(t *testing.T) {
	restaurant := domain.Restaurant{
		ID:      1,
		Name:    "Test Restaurant",
		OwnerID: 2,
	}

	expected := RestaurantDTO{
		ID:      restaurant.ID,
		Name:    restaurant.Name,
		OwnerID: restaurant.OwnerID,
	}

	result := NewRestaurantDTO(restaurant)

	require.Equal(t, expected, result, "NewRestaurantDTO did not return the expected RestaurantDTO")
}

func Test_dtos_NewRestaurant(t *testing.T) {
	dto := RestaurantDTO{
		ID:      1,
		Name:    "Test Restaurant",
		OwnerID: 2,
	}

	expected := domain.Restaurant{
		ID:      dto.ID,
		Name:    dto.Name,
		OwnerID: 2,
	}

	result := NewRestaurant(dto)

	require.Equal(t, expected, result, "RestaurantDTO.ToDomain did not return the expected domain.Restaurant")
}
