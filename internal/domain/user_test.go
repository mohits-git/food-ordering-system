package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_domain_UserRole_IsValid(t *testing.T) {
	type fields struct {
		r UserRole
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "valid role CUSTOMER",
			fields: fields{
				r: CUSTOMER,
			},
			want: true,
		},
		{
			name: "valid role OWNER",
			fields: fields{
				r: OWNER,
			},
			want: true,
		},
		{
			name: "valid role ADMIN",
			fields: fields{
				r: ADMIN,
			},
			want: true,
		},
		{
			name: "invalid role GUEST",
			fields: fields{
				r: "guest",
			},
			want: false,
		},
		{
			name: "invalid empty role",
			fields: fields{
				r: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := tt.fields.r
			got := u.IsValid()
			assert.Equal(t, tt.want, got, "UserRole.IsValid() = %v, want %v", got, tt.want)
		})
	}
}

func Test_domain_User_NewUser(t *testing.T) {
	user := NewUser(1, "testuser", "test@example.com", "12345678", CUSTOMER)
	expected := User{ID: 1, Name: "testuser", Email: "test@example.com", Password: "12345678", Role: CUSTOMER}
	assert.Equal(t, expected, user)
}

func Test_domain_User_Validate(t *testing.T) {
	tests := []struct {
		name string
		u    User
		want bool
	}{
		{
			name: "valid user",
			u: User{
				ID:       1,
				Name:     "John Doe",
				Email:    "test@example.com",
				Password: "password123",
				Role:     CUSTOMER,
			},
			want: true,
		},
		{
			name: "invalid user with empty name",
			u: User{
				ID:       2,
				Name:     "",
				Email:    "",
				Password: "password123",
				Role:     CUSTOMER,
			},
			want: false,
		},
		{
			name: "invalid user with invalid email",
			u: User{
				ID:       3,
				Name:     "Jane Doe",
				Email:    "invalid-email",
				Password: "password123",
				Role:     CUSTOMER,
			},
			want: false,
		},
		{
			name: "invalid user with invalid role",
			u: User{
				ID:       5,
				Name:     "Bob",
				Email:    "test@example.com",
				Password: "password123",
				Role:     "guest",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.u.Validate()
			assert.Equal(t, tt.want, got, "User.Validate() = %v, want %v", got, tt.want)
		})
	}
}
