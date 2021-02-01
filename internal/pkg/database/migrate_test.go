package database

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrateDBError(t *testing.T) {
	err := migrate(closedDb, []migration{
		{
			name: "0_test_migration",
			fn: func(Queryer) error {
				return nil
			},
		},
	})

	assert.Error(t, err)
}

func TestMigrateMigrationError(t *testing.T) {
	td, tDB := StartTestPostgreSQL()

	err := migrate(tDB, []migration{
		{
			name: "0_test_migration",
			fn: func(Queryer) error {
				return fmt.Errorf("testing error functionality")
			},
		},
	})

	assert.Error(t, err)
	td()
}

func TestMigrateMigrationAppliedCheckError(t *testing.T) {
	td, tDB := StartTestPostgreSQL()

	err := applyMigration(tDB, "0_test_migration", func(Queryer) error {
		return nil
	})

	assert.Error(t, err)
	td()
}

func TestMigrateMigrationApplied(t *testing.T) {
	td, tDB := StartTestPostgreSQL()
	assert.NoError(t, createMigrationTable(tDB))

	err := applyMigration(tDB, "0_test_migration", func(Queryer) error {
		return nil
	})
	assert.NoError(t, err)

	err = applyMigration(tDB, "0_test_migration", func(Queryer) error {
		return nil
	})
	assert.NoError(t, err)
	td()
}

func TestMigrateAddMigrationLogError(t *testing.T) {
	tx := db.MustBegin()

	assert.NoError(t, addMigrationLog(tx, "0_test_migration"))
	assert.Error(t, addMigrationLog(tx, "0_test_migration"))

	assert.NoError(t, tx.Rollback())
}

func TestMigrationError(t *testing.T) {
	assert.Error(t, ApplyMigrations(closedDb))
}
