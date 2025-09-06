package ports

import "context"

type AuthorizationService interface {
  Authorize(ctx context.Context, userID, resource, action string) (bool, error)
}
