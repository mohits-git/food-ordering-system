package dtos

import (
	"testing"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/stretchr/testify/require"
)

func Test_dtos_NewMenuItemResponse(t *testing.T) {
	item := domain.MenuItem{
		ID:           1,
		Name:         "Pizza",
		Price:        9.99,
		Available:    true,
		RestaurantID: 2,
	}

	expected := MenuItemResponse{
		ID:        item.ID,
		Name:      item.Name,
		Price:     item.Price,
		Available: item.Available,
	}

	result := NewMenuItemResponse(item)

	require.Equal(t, expected, result, "NewMenuItemResponse did not return the expected MenuItemResponse")
}
