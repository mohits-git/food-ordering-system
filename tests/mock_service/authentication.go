package mockservice

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type AuthenticationService struct {
	mock.Mock
}

func (s *AuthenticationService) Login(ctx context.Context, email, password string) (token string, err error) {
	args := s.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

func (s *AuthenticationService) Logout(ctx context.Context, token string) error {
	args := s.Called(ctx, token)
	return args.Error(0)
}
