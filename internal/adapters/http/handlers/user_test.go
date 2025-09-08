package handlers

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/mohits-git/food-ordering-system/internal/adapters/http/dtos"
	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
	mockservice "github.com/mohits-git/food-ordering-system/tests/mock_service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_handlers_NewUserHandler(t *testing.T) {
	mockUserService := mockservice.UserService{}
	userHandler := NewUserHandler(&mockUserService)
	require.NotNil(t, userHandler, "expected NewUserHandler() to return non-nil value")
}

func Test_handlers_HandleCreateUser_when_valid_input(t *testing.T) {
	user := domain.User{
		ID:       1,
		Name:     "Mohit",
		Email:    "test@example.com",
		Role:     "customer",
		Password: "123456",
	}
	createReq := dtos.CreateUserRequest{
		Name:     "Mohit",
		Email:    "test@example.com",
		Role:     "customer",
		Password: "123456",
	}

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, createReq)
	require.NoError(t, err, "expected no error while encoding create user request to JSON")

	req := httptest.NewRequest("POST", "/api/users", buf)
	w := httptest.NewRecorder()

	mockUserService := mockservice.UserService{}
	userHandler := NewUserHandler(&mockUserService)

	mockUserService.On("CreateUser", mock.Anything, mock.MatchedBy(func(u domain.User) bool {
		return u.Email == user.Email &&
			u.Name == user.Name &&
			u.Role == user.Role &&
			u.Password == user.Password
	})).Return(1, nil).Once()

	userHandler.HandleCreateUser(w, req)

	resp := w.Result()
	require.Equal(t, 201, resp.StatusCode, "expected status code 201 Created")
	defer resp.Body.Close()
	body, err := decodeResponse[dtos.CreateUserResponse](resp)
	require.NoError(t, err, "expected no error while decoding response body")
	require.Equal(t, user.ID, body.UserID, "expected user ID to be 1 as mock service returns 1 value")
	mockUserService.AssertExpectations(t)
}

func Test_handlers_HandleCreateUser_when_invalid_body(t *testing.T) {
	createReq := struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Role     int    `json:"role"`
		Password string `json:"password"`
	}{
		Name:     "Mohit",
		Email:    "test@example.com",
		Role:     1,
		Password: "123456",
	}

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, createReq)
	require.NoError(t, err, "expected no error while encoding create user request to JSON")

	req := httptest.NewRequest("POST", "/api/users", buf)
	w := httptest.NewRecorder()

	mockUserService := mockservice.UserService{}
	userHandler := NewUserHandler(&mockUserService)

	userHandler.HandleCreateUser(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, 400, resp.StatusCode, "expected status code 400 Bad Request")
}

func Test_handlers_HandleCreateUser_when_invalid_input(t *testing.T) {
	createReq := dtos.CreateUserRequest{
		Name:     "Mohit",
		Email:    "invalid-email",
		Role:     "customer",
		Password: "123456",
	}
	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, createReq)
	require.NoError(t, err, "expected no error while encoding create user request to JSON")

	req := httptest.NewRequest("POST", "/api/users", buf)
	w := httptest.NewRecorder()

	mockUserService := mockservice.UserService{}
	userHandler := NewUserHandler(&mockUserService)

	mockUserService.On("CreateUser", mock.Anything, mock.MatchedBy(func(u domain.User) bool {
		return u.Email == createReq.Email &&
			u.Name == createReq.Name &&
			u.Role == domain.UserRole(createReq.Role) &&
			u.Password == createReq.Password
	})).Return(0, apperr.NewAppError(apperr.ErrInvalid, "invalid user data", nil)).Once()

	userHandler.HandleCreateUser(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, 400, resp.StatusCode, "expected status code 400 Bad Request")
	mockUserService.AssertExpectations(t)
}

func Test_handlers_HandleCreateUser_when_user_already_exists(t *testing.T) {
	createReq := dtos.CreateUserRequest{
		Name:     "Mohit",
		Email:    "test@example.com",
		Role:     "customer",
		Password: "123456",
	}
	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, createReq)
	require.NoError(t, err, "expected no error while encoding create user request to JSON")

	req := httptest.NewRequest("POST", "/api/users", buf)
	w := httptest.NewRecorder()

	mockUserService := mockservice.UserService{}
	userHandler := NewUserHandler(&mockUserService)

	mockUserService.On("CreateUser", mock.Anything, mock.MatchedBy(func(u domain.User) bool {
		return u.Email == createReq.Email &&
			u.Name == createReq.Name &&
			u.Role == domain.UserRole(createReq.Role) &&
			u.Password == createReq.Password
	})).Return(0, apperr.NewAppError(apperr.ErrConflict, "user exists", nil)).Once()

	userHandler.HandleCreateUser(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, 409, resp.StatusCode, "expected status code 409 Conflict")
	mockUserService.AssertExpectations(t)
}

func Test_handlers_HandleCreateUser_when_internal_server_error(t *testing.T) {
	createReq := dtos.CreateUserRequest{
		Name:     "Mohit",
		Email:    "test@example.com",
		Role:     "customer",
		Password: "123456",
	}
	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, createReq)
	require.NoError(t, err, "expected no error while encoding create user request to JSON")

	req := httptest.NewRequest("POST", "/api/users", buf)
	w := httptest.NewRecorder()

	mockUserService := mockservice.UserService{}
	userHandler := NewUserHandler(&mockUserService)

	mockUserService.On("CreateUser", mock.Anything, mock.MatchedBy(func(u domain.User) bool {
		return u.Email == createReq.Email &&
			u.Name == createReq.Name &&
			u.Role == domain.UserRole(createReq.Role) &&
			u.Password == createReq.Password
	})).Return(0, apperr.NewAppError(apperr.ErrInternal, "database error", nil)).Once()

	userHandler.HandleCreateUser(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, 500, resp.StatusCode, "expected status code 500 Internal Server Error")
	mockUserService.AssertExpectations(t)
}

func Test_handlers_HandleGetUserById_user_exists(t *testing.T) {
	user := domain.User{
		ID:       1,
		Name:     "Mohit",
		Email:    "test@example.com",
		Role:     "customer",
		Password: "123456",
	}

	// test req, resp
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/users/1", nil)
	req.SetPathValue("id", "1")

	// mock user service
	mockUserService := mockservice.UserService{}
	mockUserService.On(
		"GetUserById",
		mock.Anything,
		1,
	).Return(user, nil).Once()

	// handler
	userHandler := NewUserHandler(&mockUserService)
	userHandler.HandleGetUserById(w, req)

	// assertions
	resp := w.Result()
	require.Equal(t, 200, resp.StatusCode, "expected status code 200 OK")

	body, err := decodeResponse[dtos.GetUserResponse](resp)
	defer resp.Body.Close()

	require.NoError(t, err, "expected no error while decoding response body")
	require.Equal(t, user.ID, body.UserID, "expected user ID to be %d but got %d", user.ID, body.UserID)
	require.Equal(t, user.Name, body.Name, "expected user Name to be %s but got %s", user.Name, body.Name)
	require.Equal(t, user.Email, body.Email, "expected user Email to be %s but got %s", user.Email, body.Email)
	role := domain.UserRole(body.Role)
	ok := role.IsValid()
	require.Truef(t, ok, "expected user role to be a valid UserRole but got %s", role)
	require.Equal(t, user.Role, role, "expected user Role to be %s but got %s", user.Role, body.Role)

	mockUserService.AssertExpectations(t)
}

func Test_handlers_HandleGetUserById_when_invalid_id(t *testing.T) {
	// test req, resp
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/users/abc", nil)
	req.SetPathValue("id", "abc")

	// mock user service
	mockUserService := mockservice.UserService{}

	// handler
	userHandler := NewUserHandler(&mockUserService)
	userHandler.HandleGetUserById(w, req)

	// assertions
	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, 400, resp.StatusCode, "expected status code 400 Bad Request")
}

func Test_handlers_HandleGetUserById_user_not_found(t *testing.T) {
	// test req, resp
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/users/1", nil)
	req.SetPathValue("id", "1")

	// mock user service
	mockUserService := mockservice.UserService{}
	mockUserService.On(
		"GetUserById",
		mock.Anything,
		1,
	).Return(domain.User{}, apperr.NewAppError(apperr.ErrNotFound, "user not found", nil)).Once()

	// handler
	userHandler := NewUserHandler(&mockUserService)
	userHandler.HandleGetUserById(w, req)

	// assertions
	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, 404, resp.StatusCode, "expected status code 404 Not Found")
	mockUserService.AssertExpectations(t)
}

func Test_handlers_HandleGetUserById_internal_server_error(t *testing.T) {
	// test req, resp
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/users/1", nil)
	req.SetPathValue("id", "1")

	// mock user service
	mockUserService := mockservice.UserService{}
	mockUserService.On(
		"GetUserById",
		mock.Anything,
		1,
	).Return(domain.User{}, apperr.NewAppError(apperr.ErrInternal, "database error", nil)).Once()

	// handler
	userHandler := NewUserHandler(&mockUserService)
	userHandler.HandleGetUserById(w, req)

	// assertions
	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, 500, resp.StatusCode, "expected status code 500 Internal Server Error")
	mockUserService.AssertExpectations(t)
}
