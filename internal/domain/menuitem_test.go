package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_domain_NewMenuItem(t *testing.T) {
	type args struct {
		id           int
		name         string
		price        float64
		available    bool
		restaurantId int
	}
	tests := []struct {
		name string
		args args
		want MenuItem
	}{
		{
			name: "valid menu item",
			args: args{
				id:           1,
				name:         "Pizza",
				price:        9.99,
				available:    true,
				restaurantId: 1,
			},
			want: MenuItem{
				ID:           1,
				Name:         "Pizza",
				Price:        9.99,
				Available:    true,
				RestaurantID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMenuItem(tt.args.id, tt.args.name, tt.args.price, tt.args.available, tt.args.restaurantId)
			assert.Equal(t, tt.want, got, "NewMenuItem() = %v, want %v", got, tt.want)
		})
	}
}

func Test_domain_MenuItem_Validate(t *testing.T) {
	tests := []struct {
		name string
		m    MenuItem
		want bool
	}{
		{
			name: "valid menu item",
			m: MenuItem{
				ID:           1,
				Name:         "Burger",
				Price:        5.99,
				Available:    true,
				RestaurantID: 1,
			},
			want: true,
		},
		{
			name: "invalid menu item with empty name",
			m: MenuItem{
				ID:           2,
				Name:         "",
				Price:        5.99,
				Available:    true,
				RestaurantID: 1,
			},
			want: false,
		},
		{
			name: "invalid menu item with negative price",
			m: MenuItem{
				ID:           3,
				Name:         "Salad",
				Price:        -1.00,
				Available:    true,
				RestaurantID: 1,
			},
			want: false,
		},
		{
			name: "invalid menu item with zero restaurant ID",
			m: MenuItem{
				ID:           4,
				Name:         "Pasta",
				Price:        7.99,
				Available:    true,
				RestaurantID: 0,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.Validate()
			assert.Equal(t, tt.want, got, "MenuItem.Validate() = %v, want %v", got, tt.want)
		})
	}
}

func Test_domain_MenuItem_IsAvailable(t *testing.T) {
	tests := []struct {
		name string
		m    MenuItem
		want bool
	}{
		{
			name: "menu item is available",
			m: MenuItem{
				ID:           1,
				Name:         "Sushi",
				Price:        12.99,
				Available:    true,
				RestaurantID: 1,
			},
			want: true,
		},
		{
			name: "menu item is not available",
			m: MenuItem{
				ID:           2,
				Name:         "Steak",
				Price:        19.99,
				Available:    false,
				RestaurantID: 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.IsAvailable()
			assert.Equal(t, tt.want, got, "MenuItem.IsAvailable() = %v, want %v", got, tt.want)
		})
	}
}
