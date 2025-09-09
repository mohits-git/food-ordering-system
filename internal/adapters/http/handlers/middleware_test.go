package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mohits-git/food-ordering-system/internal/utils/authctx"
	mocktokenprovider "github.com/mohits-git/food-ordering-system/tests/mock_token_provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_handlers_NewAuthMiddleware(t *testing.T) {
	mockTokenProvider := &mocktokenprovider.TokenProvider{}

	middleware := NewAuthMiddleware(mockTokenProvider)

	require.NotNil(t, middleware, "NewAuthMiddleware returned nil")
	require.Equal(t, mockTokenProvider, middleware.tokenProvider, "NewAuthMiddleware did not set the tokenProvider correctly")
}

func Test_handlers_Authenticated(t *testing.T) {
	mockTokenProvider := &mocktokenprovider.TokenProvider{}
	middleware := NewAuthMiddleware(mockTokenProvider)

	mockTokenProvider.On("ValidateToken", "valid-token").Return(authctx.UserClaims{
		UserID: 1,
		Role:   "user",
	}, nil)

	handler := middleware.Authenticated(func(w http.ResponseWriter, r *http.Request) {
		userClaims, ok := authctx.UserClaimsFromCtx(r.Context())
		assert.True(t, ok, "UserClaims not found in context")
		assert.Equal(t, 1, userClaims.UserID, "UserID does not match")
		w.WriteHeader(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	handler.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Result().StatusCode, "Expected status OK for valid token")
}

func Test_handlers_Authenticated_MissingToken(t *testing.T) {
	mockTokenProvider := &mocktokenprovider.TokenProvider{}
	middleware := NewAuthMiddleware(mockTokenProvider)

	handler := middleware.Authenticated(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode, "Expected status Unauthorized for missing token")
}

func Test_handlers_Authenticated_InvalidToken(t *testing.T) {
	mockTokenProvider := &mocktokenprovider.TokenProvider{}
	middleware := NewAuthMiddleware(mockTokenProvider)

	mockTokenProvider.On("ValidateToken", "invalid-token").Return(authctx.UserClaims{}, assert.AnError)

	handler := middleware.Authenticated(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	handler.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode, "Expected status Unauthorized for invalid token")
}

func Test_handlers_WithToken(t *testing.T) {
	mockTokenProvider := &mocktokenprovider.TokenProvider{}
	middleware := NewAuthMiddleware(mockTokenProvider)

	mockTokenProvider.On("ValidateToken", "valid-token").Return(authctx.UserClaims{
		UserID: 1,
		Role:   "user",
	}, nil)

	handler := middleware.WithToken(func(w http.ResponseWriter, r *http.Request) {
		token, ok := authctx.TokenFromCtx(r.Context())
		assert.True(t, ok, "Token not found in context")
		assert.Equal(t, "valid-token", token, "Token does not match")
		w.WriteHeader(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	handler.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Result().StatusCode, "Expected status OK for valid token")
}

func Test_handlers_WithToken_MissingToken(t *testing.T) {
	mockTokenProvider := &mocktokenprovider.TokenProvider{}
	middleware := NewAuthMiddleware(mockTokenProvider)

	handler := middleware.WithToken(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode, "Expected status Unauthorized for missing token")
}

func Test_handlers_WithToken_InvalidToken(t *testing.T) {
	mockTokenProvider := &mocktokenprovider.TokenProvider{}
	middleware := NewAuthMiddleware(mockTokenProvider)

	mockTokenProvider.On("ValidateToken", "invalid-token").Return(authctx.UserClaims{}, assert.AnError)

	handler := middleware.WithToken(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	handler.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode, "Expected status Unauthorized for invalid token")
}
