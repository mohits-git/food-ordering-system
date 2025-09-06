package sqlite

import (
	_ "embed"
	"log"
  "database/sql"
)

//go:embed schema.sql
var schemaSQL string

func Migrate(db *sql.DB) error {
	_, err := db.Exec(schemaSQL)
	if err != nil {
		log.Fatal("Failed to execute migrations:", err)
	}
	return err
}
