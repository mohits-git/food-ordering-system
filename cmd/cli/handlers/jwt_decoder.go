package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/utils/authctx"
)

func decodeJwt(token string) authctx.UserClaims {
	parser := jwt.Parser{}
	jwtToken, parts, err := parser.ParseUnverified(token, jwt.MapClaims{})
	if err != nil || len(parts) < 2 {
		return authctx.UserClaims{}
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return authctx.UserClaims{}
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return authctx.UserClaims{}
	}

	role, ok := claims["role"].(string)
	if !ok {
		return authctx.UserClaims{}
	}

	return authctx.UserClaims{
		UserID: int(userIDFloat),
		Role:   domain.UserRole(role),
	}
}

