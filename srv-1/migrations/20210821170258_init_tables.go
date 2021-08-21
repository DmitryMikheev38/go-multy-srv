package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upInitTables, downInitTables)
}

func upInitTables(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return nil
}

func downInitTables(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
