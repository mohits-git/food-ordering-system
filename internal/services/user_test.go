package services

import (
	"context"
	"testing"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
	mockpasswordhasher "github.com/mohits-git/food-ordering-system/tests/mock_password_hasher"
	mockrepository "github.com/mohits-git/food-ordering-system/tests/mock_repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_sqlite_NewUserService(t *testing.T) {
	mockRepo := mockrepository.UserRepository{}
	mockPasswordHasher := mockpasswordhasher.PasswordHasher{}
	userSerivce := NewUserService(&mockRepo, &mockPasswordHasher)
	require.NotNil(t, userSerivce, "required NewUserService() to return non-nil value but got nil")
	_, ok := userSerivce.(*UserSerivce)
	require.True(t, ok, "required sqlite.NewUserSerivce() to return sqlite repository but got some tother type")
}

func Test_sqlite_CreateUser_when_valid_user(t *testing.T) {
	// build/mock
	mockRepo := mockrepository.UserRepository{}
	mockPasswordHasher := mockpasswordhasher.PasswordHasher{}
	userService := NewUserService(&mockRepo, &mockPasswordHasher)
	user := domain.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Role:     domain.CUSTOMER,
		Password: "12345678",
	}

	mockPasswordHasher.On("HashPassword", user.Password).Return("hashedPassword", nil)
	// passing context as mock.Anything
	mockRepo.On("SaveUser", mock.Anything, mock.MatchedBy(func(u domain.User) bool {
		return u.Email == user.Email &&
			u.Name == user.Name &&
			u.Role == user.Role &&
			u.Password == "hashedPassword"
	})).Return(1, nil)

	// call
	id, err := userService.CreateUser(context.TODO(), user)

	// assertions
	ok := mockPasswordHasher.AssertExpectations(t)
	assert.True(t, ok, "expected all expectations to be met for mockPasswordHasher but some were not")
	ok = mockRepo.AssertExpectations(t)
	assert.True(t, ok, "expected all expectations to be met for mockRepo but some were not")
	assert.NoError(t, err, "expected CreateUser() to not return error but got %v", err)
	assert.Equal(t, 1, id, "expected CreateUser() to return id 1 but got %d", id)
}

func Test_sqlite_CreateUser_when_invalid_user(t *testing.T) {
	// build/mock
	mockRepo := mockrepository.UserRepository{}
	mockPasswordHasher := mockpasswordhasher.PasswordHasher{}
	userService := NewUserService(&mockRepo, &mockPasswordHasher)
	user := domain.User{
		Name:     "",
		Email:    "test@example.com",
		Role:     domain.CUSTOMER,
		Password: "12345678",
	}

	// call
	id, err := userService.CreateUser(context.TODO(), user)

	// assertions
	assert.Equal(t, 0, id, "expected CreateUser() to return id 0 but got %d", id)
	require.Error(t, err, "expected CreateUser() to return error but got nil")
	appErr, ok := err.(*apperr.AppError)
	require.True(t, ok, "expected error to be of type *apperr.AppError but got %T", err)
	assert.Equal(t, apperr.ErrInvalid, appErr.Code, "expected error code to be apperr.ErrInvalid but got %s", appErr.Code)
}

func Test_sqlite_CreateUser_when_password_too_long(t *testing.T) {
	// build/mock
	mockRepo := mockrepository.UserRepository{}
	mockPasswordHasher := mockpasswordhasher.PasswordHasher{}
	userService := NewUserService(&mockRepo, &mockPasswordHasher)
	user := domain.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Role:     domain.CUSTOMER,
		Password: "12345678",
	}

	mockPasswordHasher.On("HashPassword", user.Password).Return("", apperr.NewAppError(apperr.ErrInvalid, "password too long", nil))

	// call
	id, err := userService.CreateUser(context.TODO(), user)

	// assertions
	ok := mockPasswordHasher.AssertExpectations(t)
	assert.Equal(t, 0, id, "expected CreateUser() to return id 0 but got %d", id)
	require.Error(t, err, "expected CreateUser() to return error but got nil")
	appErr, ok := err.(*apperr.AppError)
	require.True(t, ok, "expected error to be of type *apperr.AppError but got %T", err)
	assert.Equal(t, apperr.ErrInvalid, appErr.Code, "expected error code to be apperr.ErrInvalid but got %s", appErr.Code)
}

func Test_sqlite_CreateUser_when_email_already_exists(t *testing.T) {
	// build/mock
	mockRepo := mockrepository.UserRepository{}
	mockPasswordHasher := mockpasswordhasher.PasswordHasher{}
	userService := NewUserService(&mockRepo, &mockPasswordHasher)
	user := domain.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Role:     domain.CUSTOMER,
		Password: "12345678",
	}
	mockPasswordHasher.On("HashPassword", user.Password).Return("hashedPassword", nil)
	// passing context as mock.Anything
	mockRepo.On("SaveUser", mock.Anything, mock.MatchedBy(func(u domain.User) bool {
		return u.Email == user.Email &&
			u.Name == user.Name &&
			u.Role == user.Role &&
			u.Password == "hashedPassword"
	})).Return(0, apperr.NewAppError(apperr.ErrConflict, "email already exists", nil))

	// call
	id, err := userService.CreateUser(context.TODO(), user)

	// assertions
	ok := mockPasswordHasher.AssertExpectations(t)
	assert.True(t, ok, "expected all expectations to be met for mockPasswordHasher but some were not")
	ok = mockRepo.AssertExpectations(t)
	assert.True(t, ok, "expected all expectations to be met for mockRepo but some were not")
	assert.Equal(t, 0, id, "expected CreateUser() to return id 0 but got %d", id)
	require.Error(t, err, "expected CreateUser() to return error but got nil")
	appErr, ok := err.(*apperr.AppError)
	require.True(t, ok, "expected error to be of type *apperr.AppError but got %T", err)
	assert.Equal(t, apperr.ErrConflict, appErr.Code, "expected error code to be apperr.ErrConflict but got %s", appErr.Code)
}

func Test_sqlite_GetUserById_when_user_exists(t *testing.T) {
	// build/mock
	mockRepo := mockrepository.UserRepository{}
	mockPasswordHasher := mockpasswordhasher.PasswordHasher{}
	userService := NewUserService(&mockRepo, &mockPasswordHasher)
	user := domain.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Role:     domain.CUSTOMER,
		Password: "hashedPassword",
	}
	// passing context as mock.Anything
	mockRepo.On("FindUserById", mock.Anything, user.ID).Return(user, nil)
	// call
	fetchedUser, err := userService.GetUserById(context.TODO(), user.ID)
	// assertions
	ok := mockRepo.AssertExpectations(t)
	assert.True(t, ok, "expected all expectations to be met for mockRepo but some were not")
	assert.NoError(t, err, "expected GetUserById() to not return error but got %v", err)
	assert.Equal(t, user, fetchedUser, "expected GetUserById() to return user %v but got %v", user, fetchedUser)
}

func Test_sqlite_GetUserById_when_user_not_exists(t *testing.T) {
	// build/mock
	mockRepo := mockrepository.UserRepository{}
	mockPasswordHasher := mockpasswordhasher.PasswordHasher{}
	userService := NewUserService(&mockRepo, &mockPasswordHasher)

	// passing context as mock.Anything
	mockRepo.On("FindUserById", mock.Anything, 1).Return(domain.User{}, apperr.NewAppError(apperr.ErrNotFound, "user not found", nil))
	// call
	fetchedUser, err := userService.GetUserById(context.TODO(), 1)
	// assertions
	ok := mockRepo.AssertExpectations(t)
	assert.True(t, ok, "expected all expectations to be met for mockRepo but some were not")
	assert.Equal(t, domain.User{}, fetchedUser, "expected GetUserById() to return empty user but got %v", fetchedUser)
	require.Error(t, err, "expected GetUserById() to return error but got nil")
	appErr, ok := err.(*apperr.AppError)
	require.True(t, ok, "expected error to be of type *apperr.AppError but got %T", err)
	assert.Equal(t, apperr.ErrNotFound, appErr.Code, "expected error code to be apperr.ErrNotFound but got %s", appErr.Code)
}
