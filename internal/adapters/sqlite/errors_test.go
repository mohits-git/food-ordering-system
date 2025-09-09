package sqlite

import (
	"testing"

	"github.com/mattn/go-sqlite3"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleSQLiteError(t *testing.T) {
	tests := []struct {
		name         string
		input        error
		expectedCode apperr.AppErrorCode
	}{
		{
			name:         "NotFound error",
			input:        sqlite3.Error{Code: sqlite3.ErrNotFound},
			expectedCode: apperr.ErrNotFound,
		},
		{
			name:         "Unique constraint violation",
			input:        sqlite3.Error{Code: sqlite3.ErrConstraint, ExtendedCode: sqlite3.ErrConstraintUnique},
			expectedCode: apperr.ErrConflict,
		},
		{
			name:         "Foreign key constraint violation",
			input:        sqlite3.Error{Code: sqlite3.ErrConstraint, ExtendedCode: sqlite3.ErrConstraintForeignKey},
			expectedCode: apperr.ErrConflict,
		},
		{
			name:         "Other constraint violation",
			input:        sqlite3.Error{Code: sqlite3.ErrConstraint, ExtendedCode: 9999},
			expectedCode: apperr.ErrConflict,
		},
		{
			name:         "Non-SQLite error",
			input:        apperr.NewAppError(apperr.ErrInternal, "some other error", nil),
			expectedCode: apperr.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := HandleSQLiteError(tt.input)
			appErr, ok := err.(*apperr.AppError)
			require.Truef(t, ok, "expected AppError type, got %T", err)
			assert.Equalf(t, tt.expectedCode, appErr.Code, "expected code %v, got %v", tt.expectedCode, appErr.Code)
		})
	}
}
