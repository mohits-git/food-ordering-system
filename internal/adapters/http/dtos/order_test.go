package dtos

import (
	"testing"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/stretchr/testify/require"
)

func Test_dtos_OrderItemsDTO_ToDomain(t *testing.T) {
	orderItemsDTO := OrderItemsDTO{
		MenuItemID: 1,
		Quantity:   2,
	}

	expected := domain.OrderItem{
		MenuItemID: orderItemsDTO.MenuItemID,
		Quantity:   orderItemsDTO.Quantity,
	}

	result := orderItemsDTO.ToDomain()

	require.Equal(t, expected, result, "OrderItemsDTO.ToDomain did not return the expected domain.OrderItem")
}
