package database

import (
	"context"
	"database/sql"
)

// Shared interface over both sqlx.Db and sqlx.Tx to enable
// querying on both without considering the underlying type
type Queryer interface {
	Exec(string, ...interface{}) (sql.Result, error)
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)

	Get(interface{}, string, ...interface{}) error
	GetContext(context.Context, interface{}, string, ...interface{}) error
}
