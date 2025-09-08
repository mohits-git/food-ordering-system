package sqlite

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/mattn/go-sqlite3"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
)

func HandleSQLiteError(err error) error {
	if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "no rows in result set") {
		return apperr.NewAppError(apperr.ErrNotFound, "record not found", err)
	}
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		switch sqliteErr.Code {
		case sqlite3.ErrEmpty, sqlite3.ErrNotFound:
			return apperr.NewAppError(apperr.ErrNotFound, "record not found", err)
		case sqlite3.ErrConstraint:
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				return apperr.NewAppError(apperr.ErrConflict, "unique constraint violation", err)
			case sqlite3.ErrConstraintForeignKey:
				return apperr.NewAppError(apperr.ErrConflict, "foreign key constraint violation", err)
			default:
				return apperr.NewAppError(apperr.ErrConflict, "constraint violation", err)
			}
		}
	}
	return apperr.NewAppError(apperr.ErrInternal, "database error", err)
}
