package ports

import "context"

type TokenService interface {
  GenerateToken(ctx context.Context, userID string) (token string, err error)
  ValidateToken(ctx context.Context, token string) (userID string, err error)
}
