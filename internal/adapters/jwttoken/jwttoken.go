package jwttoken

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/utils/authctx"
)

type JWTService struct {
	secretKey string
	issuer    string
	audience  string
}

func NewJWTService(secretKey, issuer, audience string) *JWTService {
	return &JWTService{
		secretKey: secretKey,
		issuer:    issuer,
		audience:  audience,
	}
}

func (s *JWTService) GenerateToken(claims authctx.UserClaims) (string, error) {
	jwtClaims := jwt.MapClaims{
		"iss":     s.issuer,
		"aud":     s.audience,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
		"user_id": claims.UserID,
		"role":    claims.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *JWTService) ValidateToken(tokenString string) (authctx.UserClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return authctx.UserClaims{}, err
	}

	if !token.Valid {
		return authctx.UserClaims{}, jwt.ErrTokenMalformed
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return authctx.UserClaims{}, jwt.ErrTokenInvalidClaims
	}
	if claims["iss"] != s.issuer || claims["aud"] != s.audience {
		return authctx.UserClaims{}, jwt.ErrTokenInvalidClaims
	}

	userClaims, err := s.ExtractClaims(claims)
	return userClaims, err
}

func (s *JWTService) ExtractClaims(claims jwt.MapClaims) (authctx.UserClaims, error) {
	userID, ok := claims["user_id"].(int)
	if !ok {
		return authctx.UserClaims{}, jwt.ErrTokenInvalidClaims
	}

	role, ok := claims["role"].(domain.UserRole)
	if !ok {
		return authctx.UserClaims{}, jwt.ErrTokenInvalidClaims
	}

	return authctx.UserClaims{
		UserID: userID,
		Role:   role,
	}, nil
}
