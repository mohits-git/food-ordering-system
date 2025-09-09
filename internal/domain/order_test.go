package domain

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_domain_OrderItem_Validate(t *testing.T) {
	type fields struct {
		MenuItemID int
		Quantity   int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "valid order item",
			fields: fields{
				MenuItemID: 1,
				Quantity:   2,
			},
			want: true,
		},
		{
			name: "invalid order item with zero MenuItemID",
			fields: fields{
				MenuItemID: 0,
				Quantity:   2,
			},
			want: false,
		},
		{
			name: "invalid order item with negative MenuItemID",
			fields: fields{
				MenuItemID: -1,
				Quantity:   2,
			},
			want: false,
		},
		{
			name: "invalid order item with zero Quantity",
			fields: fields{
				MenuItemID: 1,
				Quantity:   0,
			},
			want: false,
		},
		{
			name: "invalid order item with negative Quantity",
			fields: fields{
				MenuItemID: 1,
				Quantity:   -2,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oi := &OrderItem{
				MenuItemID: tt.fields.MenuItemID,
				Quantity:   tt.fields.Quantity,
			}
			got := oi.Validate()
			assert.Equal(t, tt.want, got, "OrderItem.Validate() = %v, want %v", got, tt.want)
		})
	}
}

func Test_domain_Order_Validate(t *testing.T) {
	type fields struct {
		ID           int
		CustomerID   int
		RestaurantID int
		Items        []OrderItem
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "valid order",
			fields: fields{
				ID:           1,
				CustomerID:   1,
				RestaurantID: 1,
				Items: []OrderItem{
					{MenuItemID: 1, Quantity: 2},
					{MenuItemID: 2, Quantity: 1},
				},
			},
			want: true,
		},
		{
			name: "invalid order with zero CustomerID",
			fields: fields{
				ID:           2,
				CustomerID:   0,
				RestaurantID: 1,
				Items: []OrderItem{
					{MenuItemID: 1, Quantity: 2},
				},
			},
			want: false,
		},
		{
			name: "invalid order with negative CustomerID",
			fields: fields{
				ID:           3,
				CustomerID:   -1,
				RestaurantID: 1,
				Items: []OrderItem{
					{MenuItemID: 1, Quantity: 2},
				},
			},
			want: false,
		},
		{
			name: "invalid order with zero RestaurantID",
			fields: fields{
				ID:           4,
				CustomerID:   1,
				RestaurantID: 0,
				Items: []OrderItem{
					{MenuItemID: 1, Quantity: 2},
				},
			},
			want: false,
		},
		{
			name: "invalid order with negative RestaurantID",
			fields: fields{
				ID:           5,
				CustomerID:   1,
				RestaurantID: -1,
				Items: []OrderItem{
					{MenuItemID: 1, Quantity: 2},
				},
			},
			want: false,
		},
		{
			name: "invalid order with empty Items",
			fields: fields{
				ID:           6,
				CustomerID:   1,
				RestaurantID: 1,
				Items:        []OrderItem{},
			},
			want: false,
		},
		{
			name: "invalid order with invalid OrderItem",
			fields: fields{
				ID:           7,
				CustomerID:   1,
				RestaurantID: 1,
				Items: []OrderItem{
					{MenuItemID: 0, Quantity: 2},
				},
			},
			want: false,
		},
		{
			name: "invalid order with unavailable MenuItemID",
			fields: fields{
				ID:           8,
				CustomerID:   1,
				RestaurantID: 1,
				Items: []OrderItem{
					{MenuItemID: 3, Quantity: 2},
				},
			},
			want: false,
		},
	}
	menuItems := map[int]bool{
		1: true,
		2: true,
		3: false,
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Order{
				ID:           tt.fields.ID,
				CustomerID:   tt.fields.CustomerID,
				RestaurantID: tt.fields.RestaurantID,
				OrderItems:   tt.fields.Items,
			}
			got := o.Validate(menuItems)
			assert.Equal(t, tt.want, got, "Order.Validate() = %v, want %v", got, tt.want)
		})
	}
}

func Test_domain_Order_AddItem(t *testing.T) {
	type fields struct {
		ID           int
		CustomerID   int
		RestaurantID int
		Items        []OrderItem
	}
	type args struct {
		menuItemID int
		quantity   int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []OrderItem
	}{
		{
			name: "add new item to empty order",
			fields: fields{
				ID:           1,
				CustomerID:   1,
				RestaurantID: 1,
				Items:        []OrderItem{},
			},
			args: args{
				menuItemID: 1,
				quantity:   2,
			},
			want: []OrderItem{
				{MenuItemID: 1, Quantity: 2},
			},
		},
		{
			name: "add new item to non-empty order",
			fields: fields{
				ID:           2,
				CustomerID:   1,
				RestaurantID: 1,
				Items: []OrderItem{
					{MenuItemID: 1, Quantity: 2},
				},
			},
			args: args{
				menuItemID: 2,
				quantity:   1,
			},
			want: []OrderItem{
				{MenuItemID: 1, Quantity: 2},
				{MenuItemID: 2, Quantity: 1},
			},
		},
		{
			name: "add existing item to order",
			fields: fields{
				ID:           3,
				CustomerID:   1,
				RestaurantID: 1,
				Items: []OrderItem{
					{MenuItemID: 1, Quantity: 2},
				},
			},
			args: args{
				menuItemID: 1,
				quantity:   3,
			},
			want: []OrderItem{
				{MenuItemID: 1, Quantity: 5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Order{
				ID:           tt.fields.ID,
				CustomerID:   tt.fields.CustomerID,
				RestaurantID: tt.fields.RestaurantID,
				OrderItems:   tt.fields.Items,
			}
			o.AddItem(tt.args.menuItemID, tt.args.quantity)
			if !slices.EqualFunc(o.OrderItems, tt.want, func(a, b OrderItem) bool {
				return a.MenuItemID == b.MenuItemID && a.Quantity == b.Quantity
			}) {
				t.Errorf("Order.AddItem() = %v, want %v", o.OrderItems, tt.want)
			}
		})
	}
}

func Test_domain_Order_ClearItems(t *testing.T) {
	type fields struct {
		ID           int
		CustomerID   int
		RestaurantID int
		Items        []OrderItem
	}
	tests := []struct {
		name   string
		fields fields
		want   []OrderItem
	}{
		{
			name: "clear items from non-empty order",
			fields: fields{
				ID:           1,
				CustomerID:   1,
				RestaurantID: 1,
				Items: []OrderItem{
					{MenuItemID: 1, Quantity: 2},
					{MenuItemID: 2, Quantity: 1},
				},
			},
			want: []OrderItem{},
		},
		{
			name: "clear items from already empty order",
			fields: fields{
				ID:           2,
				CustomerID:   1,
				RestaurantID: 1,
				Items:        []OrderItem{},
			},
			want: []OrderItem{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Order{
				ID:           tt.fields.ID,
				CustomerID:   tt.fields.CustomerID,
				RestaurantID: tt.fields.RestaurantID,
				OrderItems:   tt.fields.Items,
			}
			o.ClearItems()
			if !slices.EqualFunc(o.OrderItems, tt.want, func(a, b OrderItem) bool {
				return a.MenuItemID == b.MenuItemID && a.Quantity == b.Quantity
			}) {
				t.Errorf("Order.ClearItems() = %v, want %v", o.OrderItems, tt.want)
			}
		})
	}
}

func Test_domain_Order_TotalPrice(t *testing.T) {
	type fields struct {
		ID           int
		CustomerID   int
		RestaurantID int
		Items        []OrderItem
	}
	tests := []struct {
		name      string
		fields    fields
		menuItems map[int]MenuItem
		want      float64
	}{
		{
			name: "calculate total price for order with multiple items",
			fields: fields{
				ID:           1,
				CustomerID:   1,
				RestaurantID: 1,
				Items: []OrderItem{
					{MenuItemID: 1, Quantity: 2},
					{MenuItemID: 2, Quantity: 1},
				},
			},
			menuItems: map[int]MenuItem{
				1: {ID: 1, Name: "Burger", Price: 5.00, Available: true, RestaurantID: 1},
				2: {ID: 2, Name: "Fries", Price: 3.00, Available: true, RestaurantID: 1},
			},
			want: 13.00,
		},
		{
			name: "calculate total price for order with single item",
			fields: fields{
				ID:           2,
				CustomerID:   1,
				RestaurantID: 1,
				Items: []OrderItem{
					{MenuItemID: 1, Quantity: 3},
				},
			},
			menuItems: map[int]MenuItem{
				1: {ID: 1, Name: "Pizza", Price: 8.00, Available: true, RestaurantID: 1},
			},
			want: 24.00,
		},
		{
			name: "calculate total price for order with no items",
			fields: fields{
				ID:           3,
				CustomerID:   1,
				RestaurantID: 1,
				Items:        []OrderItem{},
			},
			menuItems: map[int]MenuItem{
				1: {ID: 1, Name: "Salad", Price: 4.00, Available: true, RestaurantID: 1},
			},
			want: 0.00,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Order{
				ID:           tt.fields.ID,
				CustomerID:   tt.fields.CustomerID,
				RestaurantID: tt.fields.RestaurantID,
				OrderItems:   tt.fields.Items,
			}
			got := o.TotalPrice(tt.menuItems)
			assert.Equal(t, tt.want, got, "Order.TotalPrice() = %v, want %v", got, tt.want)
		})
	}
}
