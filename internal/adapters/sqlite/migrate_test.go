package sqlite

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestMigrate(t *testing.T) {
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
