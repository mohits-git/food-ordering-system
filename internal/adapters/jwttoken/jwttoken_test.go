package jwttoken

import (
	"testing"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/utils/authctx"
)

func Test_jwttoken_NewJWTService(t *testing.T) {
  jwtService := NewJWTService("mysecretkey", "myissuer", "myaudience")
  if jwtService.secretKey != "mysecretkey" {
    t.Errorf("expected secretKey to be 'mysecretkey', got %s", jwtService.secretKey)
  }
  if jwtService.issuer != "myissuer" {
    t.Errorf("expected issuer to be 'myissuer', got %s", jwtService.issuer)
  }
  if jwtService.audience != "myaudience" {
    t.Errorf("expected audience to be 'myaudience', got %s", jwtService.audience)
  }
}

func Test_jwttoken_GenerateToken(t *testing.T) {
  jwtService := NewJWTService("mysecretkey", "myissuer", "myaudience")
  claims := authctx.UserClaims{
    UserID: 1,
    Role:   domain.CUSTOMER,
  }

  token, err := jwtService.GenerateToken(claims)
  if err != nil {
    t.Errorf("expected no error, got %v", err)
  }
  if token == "" {
    t.Error("expected token to be non-empty")
  }
}

func Test_jwttoken_ValidateToken_when_valid(t *testing.T) {
  jwtService := NewJWTService("mysecretkey", "myissuer", "myaudience")
  claims := authctx.UserClaims{
    UserID: 1,
    Role:   domain.CUSTOMER,
  }

  token, err := jwtService.GenerateToken(claims)
  if err != nil {
    t.Fatalf("expected no error generating token, got %v", err)
  }

  validatedClaims, err := jwtService.ValidateToken(token)
  if err != nil {
    t.Errorf("expected no error validating token, got %v", err)
  }
  if validatedClaims.UserID != claims.UserID {
    t.Errorf("expected UserID to be %d, got %d", claims.UserID, validatedClaims.UserID)
  }
  if validatedClaims.Role != claims.Role {
    t.Errorf("expected Role to be %s, got %s", claims.Role, validatedClaims.Role)
  }
}

func Test_jwttoken_ValidateToken_when_signed_with_invalid_secret(t *testing.T) {
  claims := authctx.UserClaims{
    UserID: 1,
    Role:   domain.CUSTOMER,
  }

  jwtService := NewJWTService("mysecretkey", "myissuer", "myaudience")
  token, err := jwtService.GenerateToken(claims)
  if err != nil {
    t.Fatalf("expected no error generating token, got %v", err)
  }

  anotherJWTService := NewJWTService("anothersecretkey", "myissuer", "myaudience")

  _, err = anotherJWTService.ValidateToken(token)
  if err == nil {
    t.Fatalf("expected error while validating token, got %v", err)
  }
}


func Test_jwttoken_ValidateToken_when_signed_with_invalid_iss_or_aud(t *testing.T) {
  claims := authctx.UserClaims{
    UserID: 1,
    Role:   domain.CUSTOMER,
  }

  jwtService := NewJWTService("mysecretkey", "myissuer", "myaudience")
  token, err := jwtService.GenerateToken(claims)
  if err != nil {
    t.Fatalf("expected no error generating token, got %v", err)
  }

  anotherJWTService := NewJWTService("mysecretkey", "anotherissuer", "myaudience")

  _, err = anotherJWTService.ValidateToken(token)
  if err == nil {
    t.Errorf("expected error while validating token, got %v", err)
  }

  anotherJWTService = NewJWTService("mysecretkey", "myissuer", "anotheraudience")
  _, err = anotherJWTService.ValidateToken(token)
  if err == nil {
    t.Errorf("expected error while validating token, got %v", err)
  }
}
