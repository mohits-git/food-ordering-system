package domain

import "regexp"

type UserRole string

const (
	CUSTOMER UserRole = "customer"
	OWNER    UserRole = "owner"
	ADMIN    UserRole = "admin"
)

func (r UserRole) IsValid() bool {
	switch r {
	case CUSTOMER, OWNER, ADMIN:
		return true
	}
	return false
}

type User struct {
	ID       int
	Name     string
	Email    string
	Role     UserRole
	Password string
}

func NewUser(id int, name, email, password string, role UserRole) User {
	return User{
		ID:       id,
		Name:     name,
		Email:    email,
		Role:     role,
		Password: password,
	}
}

func (u *User) Validate() bool {
	if u.Name == "" || u.Email == "" || u.Password == "" {
		return false
	}

	matched, err := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, u.Email)
	if err != nil || !matched {
		return false
	}

	return u.Role.IsValid()
}
