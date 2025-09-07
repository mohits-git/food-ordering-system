package ports

import "context"

type AuthenticationService interface {
	Login(ctx context.Context, email, password string) (token string, err error)
	Logout(ctx context.Context, token string) error
}
