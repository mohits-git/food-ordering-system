package authctx

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
)

type userKeyType string

const userKey userKeyType = "user"

type UserClaims struct {
	UserID int
	Role   domain.UserRole
}

func NewUserClaims(userID int, role domain.UserRole) UserClaims {
	return UserClaims{
		UserID: userID,
		Role:   role,
	}
}

func WithUserClaims(ctx context.Context, claims *UserClaims) context.Context {
	return context.WithValue(ctx, userKey, claims)
}

func UserClaimsFromCtx(ctx context.Context) (*UserClaims, bool) {
	val := ctx.Value(userKey)
	if val == nil {
		return nil, false
	}

	claims, ok := val.(*UserClaims)
	return claims, ok
}
