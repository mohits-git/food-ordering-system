package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_domain_NewRestaurant(t *testing.T) {
	type args struct {
		id      int
		name    string
		ownerID int
	}
	tests := []struct {
		name string
		args args
		want Restaurant
	}{
		{
			name: "create restaurant with valid data",
			args: args{
				id:      1,
				name:    "Test Restaurant",
				ownerID: 1,
			},
			want: Restaurant{
				ID:      1,
				Name:    "Test Restaurant",
				OwnerID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRestaurant(tt.args.id, tt.args.name, tt.args.ownerID)
			assert.Equal(t, tt.want, got, "NewRestaurant() = %v, want %v", got, tt.want)
		})
	}
}

func Test_domain_Restaurant_Validate(t *testing.T) {
	tests := []struct {
		name string
		r    Restaurant
		want bool
	}{
		{
			name: "valid restaurant",
			r: Restaurant{
				ID:      1,
				Name:    "Valid Restaurant",
				OwnerID: 1,
			},
			want: true,
		},
		{
			name: "invalid restaurant with empty name",
			r: Restaurant{
				ID:      2,
				Name:    "",
				OwnerID: 1,
			},
			want: false,
		},
		{
			name: "invalid restaurant with zero ownerID",
			r: Restaurant{
				ID:      3,
				Name:    "No Owner Restaurant",
				OwnerID: 0,
			},
			want: false,
		},
		{
			name: "invalid restaurant with negative ownerID",
			r: Restaurant{
				ID:      4,
				Name:    "Negative Owner Restaurant",
				OwnerID: -1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.r.Validate()
			assert.Equal(t, tt.want, got, "Restaurant.Validate() = %v, want %v", got, tt.want)
		})
	}
}
