package sqlite

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mattn/go-sqlite3"
	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
	"github.com/stretchr/testify/require"
)

func Test_sqlite_NewUserRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Expected no error when creating sqlmock")
	defer db.Close()
	_ = mock // mock is not used in this test but can be used for further tests

	repo := NewUserRepository(db)
	require.NotNil(t, repo, "Expected NewUserRepository to return a non-nil repository")
	repository, ok := repo.(*UserRepository)
	require.True(t, ok, "Expected repository to be of type *UserRepository")
	require.Equal(t, db, repository.db, "Expected repository db to match the provided db")
}

func Test_sqlite_UserRepository_FindUserById(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Expected no error when creating sqlmock")
	defer db.Close()

	repo := NewUserRepository(db)
	require.NotNil(t, repo, "Expected NewUserRepository to return a non-nil repository")

	// Define test cases
	tests := []struct {
		name              string
		userId            int
		mockSetup         func()
		expectedUser      domain.User
		expectedErrorCode apperr.AppErrorCode
	}{
		{
			name:   "User found",
			userId: 1,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "password"}).
					AddRow(1, "John Doe", "test@example.com", "customer", "hashedpassword")
				mock.ExpectQuery("SELECT id, name, email, role, password FROM users WHERE id").
					WithArgs(int64(1)).
					WillReturnRows(rows)

			},
			expectedUser: domain.User{
				ID:       1,
				Name:     "John Doe",
				Email:    "test@example.com",
				Role:     "customer",
				Password: "hashedpassword",
			},
			expectedErrorCode: apperr.ErrNone,
		},
		{
			name:   "User not found",
			userId: 2,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "password"})
				mock.ExpectQuery("SELECT id, name, email, role, password FROM users WHERE id").
					WithArgs(int64(2)).
					WillReturnRows(rows)
			},
			expectedUser:      domain.User{},
			expectedErrorCode: apperr.ErrNotFound,
		},
		{
			name:   "Database error",
			userId: 3,
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, name, email, role, password FROM users WHERE id").
					WithArgs(int64(3)).
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:      domain.User{},
			expectedErrorCode: apperr.ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			user, err := repo.FindUserById(t.Context(), tt.userId)
			require.Equal(t, tt.expectedUser, user, "Expected user to match")
			if err == nil && tt.expectedErrorCode == apperr.ErrNone {
				err = mock.ExpectationsWereMet()
				require.NoError(t, err, "Expected all sqlmock expectations to be met")
				return
			}
			appErr, ok := err.(*apperr.AppError)
			require.True(t, ok, "Expected error to be of type *apperr.AppError")
			if tt.expectedErrorCode == apperr.ErrNone {
				require.NoError(t, err, "Expected no error")
			} else {
				require.Error(t, err, "Expected an error")
				require.Equal(t, tt.expectedErrorCode, appErr.Code, "Expected error code to match")
			}
			err = mock.ExpectationsWereMet()
			require.NoError(t, err, "Expected all sqlmock expectations to be met")
		})
	}
}

func Test_sqlite_UserRepository_FindUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Expected no error when creating sqlmock")
	defer db.Close()

	repo := NewUserRepository(db)
	require.NotNil(t, repo, "Expected NewUserRepository to return a non-nil repository")

	// Define test cases
	tests := []struct {
		name              string
		email             string
		mockSetup         func()
		expectedUser      domain.User
		expectedErrorCode apperr.AppErrorCode
	}{
		{
			name:  "User found",
			email: "test@example.com",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "password"}).
					AddRow(1, "John Doe", "test@example.com", "customer", "hashedpassword")
				mock.ExpectQuery("SELECT id, name, email, role, password FROM users WHERE email").
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			expectedUser: domain.User{
				ID:       1,
				Name:     "John Doe",
				Email:    "test@example.com",
				Role:     "customer",
				Password: "hashedpassword",
			},
			expectedErrorCode: apperr.ErrNone,
		},
		{
			name:  "User not found",
			email: "test@example.com",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "password"})
				mock.ExpectQuery("SELECT id, name, email, role, password FROM users WHERE email").
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			expectedUser:      domain.User{},
			expectedErrorCode: apperr.ErrNotFound,
		},
		{
			name:  "Database error",
			email: "test@example.com",
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, name, email, role, password FROM users WHERE email").
					WithArgs("test@example.com").
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:      domain.User{},
			expectedErrorCode: apperr.ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			user, err := repo.FindUserByEmail(t.Context(), tt.email)
			require.Equal(t, tt.expectedUser, user, "Expected user to match")
			if err == nil && tt.expectedErrorCode == apperr.ErrNone {
				err = mock.ExpectationsWereMet()
				require.NoError(t, err, "Expected all sqlmock expectations to be met")
				return
			}
			appErr, ok := err.(*apperr.AppError)
			require.True(t, ok, "Expected error to be of type *apperr.AppError")
			if tt.expectedErrorCode == apperr.ErrNone {
				require.NoError(t, err, "Expected no error")
			} else {
				require.Error(t, err, "Expected an error")
				require.Equal(t, tt.expectedErrorCode, appErr.Code, "Expected error code to match")
			}
			err = mock.ExpectationsWereMet()
			require.NoError(t, err, "Expected all sqlmock expectations to be met")
		})
	}
}

func Test_sqlite_UserRepository_SaveUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Expected no error when creating sqlmock")
	defer db.Close()

	repo := NewUserRepository(db)
	require.NotNil(t, repo, "Expected NewUserRepository to return a non-nil repository")

	// Define test cases
	tests := []struct {
		name              string
		user              domain.User
		mockSetup         func()
		expectedUserId    int
		expectedErrorCode apperr.AppErrorCode
	}{
		{
			name: "User created successfully",
			user: domain.User{
				Name:     "John Doe",
				Email:    "test@example.com",
				Role:     "customer",
				Password: "hashedpassword",
			},
			mockSetup: func() {
				mock.ExpectExec("INSERT INTO users \\(name, email, role, password\\) VALUES \\(\\?, \\?, \\?, \\?\\)").
					WithArgs("John Doe", "test@example.com", "customer", "hashedpassword").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedUserId:    1,
			expectedErrorCode: apperr.ErrNone,
		},
		{
			name: "Unique constraint violation",
			user: domain.User{
				Name:     "John Doe",
				Email:    "test@example.com",
				Role:     "customer",
				Password: "hashedpassword",
			},
			mockSetup: func() {
				mock.ExpectExec("INSERT INTO users \\(name, email, role, password\\) VALUES \\(\\?, \\?, \\?, \\?\\)").
					WithArgs("John Doe", "test@example.com", "customer", "hashedpassword").
					WillReturnError(sqlite3.Error{Code: sqlite3.ErrConstraint, ExtendedCode: sqlite3.ErrConstraintUnique})
			},
			expectedUserId:    0,
			expectedErrorCode: apperr.ErrConflict,
		},
		{
			name: "Database error",
			user: domain.User{
				Name:     "John Doe",
				Email:    "test@example.com",
				Role:     "customer",
				Password: "hashedpassword",
			},
			mockSetup: func() {
				mock.ExpectExec("INSERT INTO users \\(name, email, role, password\\) VALUES \\(\\?, \\?, \\?, \\?\\)").
					WithArgs("John Doe", "test@example.com", "customer", "hashedpassword").
					WillReturnError(sql.ErrConnDone)
			},
			expectedUserId:    0,
			expectedErrorCode: apperr.ErrInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			userId, err := repo.SaveUser(t.Context(), tt.user)
			require.Equal(t, tt.expectedUserId, userId, "Expected user to match")
			if err == nil && tt.expectedErrorCode == apperr.ErrNone {
				err = mock.ExpectationsWereMet()
				require.NoError(t, err, "Expected all sqlmock expectations to be met")
				return
			}
			appErr, ok := err.(*apperr.AppError)
			require.True(t, ok, "Expected error to be of type *apperr.AppError")
			if tt.expectedErrorCode == apperr.ErrNone {
				require.NoError(t, err, "Expected no error")
			} else {
				require.Error(t, err, "Expected an error")
				require.Equal(t, tt.expectedErrorCode, appErr.Code, "Expected error code to match")
			}
			err = mock.ExpectationsWereMet()
			require.NoError(t, err, "Expected all sqlmock expectations to be met")
		})
	}
}
