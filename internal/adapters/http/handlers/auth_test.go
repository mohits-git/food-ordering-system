package handlers

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/mohits-git/food-ordering-system/internal/adapters/http/dtos"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
	"github.com/mohits-git/food-ordering-system/internal/utils/authctx"
	mockservice "github.com/mohits-git/food-ordering-system/tests/mock_service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_handlers_AuthHandler_NewAuthHandler(t *testing.T) {
	mockAuthService := &mockservice.AuthenticationService{}
	handler := NewAuthHandler(mockAuthService)
	require.NotNil(t, handler, "expected NewAuthHandler to return a non-nil handler")
}

func Test_handlers_AuthHandler_HandleLogin(t *testing.T) {
	mockAuthService := &mockservice.AuthenticationService{}
	handler := NewAuthHandler(mockAuthService)
	require.NotNil(t, handler, "expected NewAuthHandler to return a non-nil handler")

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.LoginRequest{
		Email:    "test@example.com",
		Password: "12345678",
	})
	require.NoError(t, err, "expected no error while encoding login request")
	req := httptest.NewRequest("POST", "/api/auth/login", buf)

	w := httptest.NewRecorder()
	mockAuthService.On("Login", mock.Anything, "test@example.com", "12345678").Return(
		"mocked-jwt-token", nil).Once()

	handler.HandleLogin(w, req)
	res := w.Result()

	require.Equal(t, 200, res.StatusCode, "expected status code 200")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")
	require.Equal(t, res.Header.Get("Set-Cookie"), "token=mocked-jwt-token; HttpOnly; Path=/api/; SameSite=Strict", "expected Set-Cookie header to be set with the token")

	defer res.Body.Close()
	loginResponse, err := decodeResponse[dtos.LoginResponse](res)
	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, "mocked-jwt-token", loginResponse.Token, "expected token to be 'mocked-jwt-token'")
	ok := mockAuthService.AssertExpectations(t)
	require.True(t, ok, "expected all expectations to be met for mockAuthService but some were not")
}

func Test_handlers_AuthHandler_HandleLogin_InvalidRequest(t *testing.T) {
	mockAuthService := &mockservice.AuthenticationService{}
	handler := NewAuthHandler(mockAuthService)
	require.NotNil(t, handler, "expected NewAuthHandler to return a non-nil handler")

	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBufferString("invalid-json"))
	w := httptest.NewRecorder()

	handler.HandleLogin(w, req)
	res := w.Result()

	require.Equal(t, 400, res.StatusCode, "expected status code 400 for invalid request")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 400, errorResponse.Status, "expected error status to be 400")
}

func Test_handlers_AuthHandler_HandleLogin_Unauthorized(t *testing.T) {
	mockAuthService := &mockservice.AuthenticationService{}
	handler := NewAuthHandler(mockAuthService)
	require.NotNil(t, handler, "expected NewAuthHandler to return a non-nil handler")

	mockAuthService.On("Login", mock.Anything, "test@example.com", "12345678").Return(
		"", apperr.NewAppError(apperr.ErrUnauthorized, "invalid credentials", nil)).Once()

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.LoginRequest{
		Email:    "test@example.com",
		Password: "12345678",
	})
	require.NoError(t, err, "expected no error while encoding login request")
	req := httptest.NewRequest("POST", "/api/auth/login", buf)
	w := httptest.NewRecorder()

	handler.HandleLogin(w, req)
	res := w.Result()

	require.Equal(t, 401, res.StatusCode, "expected status code 401 for unauthorized request")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 401, errorResponse.Status, "expected error status to be 401")
}

func Test_handlers_AuthHandler_HandleLogin_NotFound(t *testing.T) {
	mockAuthService := &mockservice.AuthenticationService{}
	handler := NewAuthHandler(mockAuthService)
	require.NotNil(t, handler, "expected NewAuthHandler to return a non-nil handler")

	mockAuthService.On("Login", mock.Anything, "test@example.com", "12345678").Return(
		"", apperr.NewAppError(apperr.ErrNotFound, "user not found", nil)).Once()

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.LoginRequest{
		Email:    "test@example.com",
		Password: "12345678",
	})
	require.NoError(t, err, "expected no error while encoding login request")
	req := httptest.NewRequest("POST", "/api/auth/login", buf)
	w := httptest.NewRecorder()

	handler.HandleLogin(w, req)
	res := w.Result()

	require.Equal(t, 401, res.StatusCode, "expected status code 401 for unauthorized request")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 401, errorResponse.Status, "expected error status to be 401")
}

func Test_handlers_AuthHandler_HandleLogin_BadRequest(t *testing.T) {
	mockAuthService := &mockservice.AuthenticationService{}
	handler := NewAuthHandler(mockAuthService)
	require.NotNil(t, handler, "expected NewAuthHandler to return a non-nil handler")

	mockAuthService.On("Login", mock.Anything, "test@example.com", "12345678").Return(
		"", apperr.NewAppError(apperr.ErrInvalid, "bad request", nil)).Once()

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.LoginRequest{
		Email:    "test@example.com",
		Password: "12345678",
	})
	require.NoError(t, err, "expected no error while encoding login request")
	req := httptest.NewRequest("POST", "/api/auth/login", buf)
	w := httptest.NewRecorder()

	handler.HandleLogin(w, req)
	res := w.Result()

	require.Equal(t, 400, res.StatusCode, "expected status code 400 for bad request")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 400, errorResponse.Status, "expected error status to be 400")
}

func Test_handlers_AuthHandler_HandleLogin_InternalServerError(t *testing.T) {
	mockAuthService := &mockservice.AuthenticationService{}
	handler := NewAuthHandler(mockAuthService)
	require.NotNil(t, handler, "expected NewAuthHandler to return a non-nil handler")

	mockAuthService.On("Login", mock.Anything, "test@example.com", "12345678").Return(
		"", apperr.NewAppError(apperr.ErrInternal, "internal server error", nil)).Once()

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.LoginRequest{
		Email:    "test@example.com",
		Password: "12345678",
	})
	require.NoError(t, err, "expected no error while encoding login request")
	req := httptest.NewRequest("POST", "/api/auth/login", buf)
	w := httptest.NewRecorder()

	handler.HandleLogin(w, req)
	res := w.Result()

	require.Equal(t, 500, res.StatusCode, "expected status code 500 for internal server error")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 500, errorResponse.Status, "expected error status to be 500")
}

func Test_handlers_AuthHandler_HandleLogout_Success(t *testing.T) {
	mockAuthService := &mockservice.AuthenticationService{}
	handler := NewAuthHandler(mockAuthService)
	require.NotNil(t, handler, "expected NewAuthHandler to return a non-nil handler")

	req := httptest.NewRequest("POST", "/api/auth/logout", nil)
	req = req.WithContext(authctx.WithToken(req.Context(), "mocked-jwt-token"))
	w := httptest.NewRecorder()

	mockAuthService.On("Logout", mock.Anything, "mocked-jwt-token").Return(nil).Once()

	handler.HandleLogout(w, req)
	res := w.Result()

	require.Equal(t, 200, res.StatusCode, "expected status code 200 for successful logout")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")
	require.Equal(t, res.Header.Get("Set-Cookie"), "token=; HttpOnly; Path=/api/; Max-Age=0; SameSite=Strict", "expected Set-Cookie header to clear the token")

	defer res.Body.Close()
	logoutResponse, err := decodeResponse[dtos.LogoutResponse](res)
	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, dtos.LogoutResponse{}, logoutResponse, "expected empty logout response")
	ok := mockAuthService.AssertExpectations(t)
	require.True(t, ok, "expected all expectations to be met for mockAuthService but some were not")
}

func Test_handlers_AuthHandler_HandleLogout_Unauthorized(t *testing.T) {
	mockAuthService := &mockservice.AuthenticationService{}
	handler := NewAuthHandler(mockAuthService)
	require.NotNil(t, handler, "expected NewAuthHandler to return a non-nil handler")

	req := httptest.NewRequest("POST", "/api/auth/logout", nil)
	w := httptest.NewRecorder()

	handler.HandleLogout(w, req)
	res := w.Result()

	require.Equal(t, 401, res.StatusCode, "expected status code 401 for unauthorized logout")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 401, errorResponse.Status, "expected error status to be 401")
}

func Test_handlers_AuthHandler_HandleLogout_InternalServerError(t *testing.T) {
	mockAuthService := &mockservice.AuthenticationService{}
	handler := NewAuthHandler(mockAuthService)
	require.NotNil(t, handler, "expected NewAuthHandler to return a non-nil handler")

	req := httptest.NewRequest("POST", "/api/auth/logout", nil)
	req = req.WithContext(authctx.WithToken(req.Context(), "mocked-jwt-token"))
	w := httptest.NewRecorder()

	mockAuthService.On("Logout", mock.Anything, "mocked-jwt-token").Return(
		apperr.NewAppError(apperr.ErrInternal, "internal server error", nil)).Once()

	handler.HandleLogout(w, req)
	res := w.Result()

	require.Equal(t, 500, res.StatusCode, "expected status code 500 for internal server error")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 500, errorResponse.Status, "expected error status to be 500")
}
