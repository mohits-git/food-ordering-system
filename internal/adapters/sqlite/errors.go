package sqlite

import (
	"github.com/mattn/go-sqlite3"
	"github.com/mohits-git/food-ordering-system/internal/utils"
)

func HandleSQLiteError(err error) error {
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		switch sqliteErr.Code {
		case sqlite3.ErrNotFound:
			return utils.NewAppError(utils.ErrNotFound, "record not found", err)
		case sqlite3.ErrConstraint:
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				return utils.NewAppError(utils.ErrConflict, "unique constraint violation", err)
			case sqlite3.ErrConstraintForeignKey:
				return utils.NewAppError(utils.ErrInvalid, "foreign key constraint violation", err)
			default:
				return utils.NewAppError(utils.ErrConflict, "constraint violation", err)
			}
		}
	}
	return utils.NewAppError(utils.ErrInternal, "database error", err)
}
