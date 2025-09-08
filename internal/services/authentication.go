package services

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/ports"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
	"github.com/mohits-git/food-ordering-system/internal/utils/authctx"
)

type AuthenticationService struct {
	userRepo       ports.UserRepository
	tokenProvider  ports.TokenProvider
	passwordHasher ports.PasswordHasher
}

func NewAuthenticationService(
	userRepo ports.UserRepository,
	tokenProvider ports.TokenProvider,
	passwordHasher ports.PasswordHasher,
) *AuthenticationService {
	return &AuthenticationService{
		userRepo:       userRepo,
		tokenProvider:  tokenProvider,
		passwordHasher: passwordHasher,
	}
}

func (s *AuthenticationService) Login(ctx context.Context, email, password string) (token string, err error) {
	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	match, err := s.passwordHasher.ComparePassword(user.Password, password)
	if err != nil {
		return "", err
	}
	if !match {
		return "", apperr.NewAppError(apperr.ErrUnauthorized, "invalid email or password", nil)
	}

	claims := authctx.NewUserClaims(user.ID, user.Role)
	token, err = s.tokenProvider.GenerateToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthenticationService) Logout(ctx context.Context, token string) error {
	// no op logout for stateless JWT
	// if we add blacklisting, we can implement it here
	// or refresh token mechanism
	return nil
}
