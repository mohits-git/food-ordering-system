package sqlite

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func Test_sqlite_Migrate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// the starting query of the schema.sql query
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = Migrate(db)
	if err != nil {
		t.Errorf("unexpected error during migration: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_sqlite_Migrate_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// the starting query of the schema.sql query
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").
		WillReturnError(sqlite3.ErrInternal)

	err = Migrate(db)
	require.Errorf(t, err, "unexpected error during migration: %s", err)

	err = mock.ExpectationsWereMet()
	require.NoErrorf(t, err, "there were unfulfilled expectations: %s", err)
}
