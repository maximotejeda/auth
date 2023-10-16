package sqles

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	_ "modernc.org/sqlite"
)

const (
	DefaultDriver   = "sqlite"
	DefaultConnMode = "wrc"
	DefaultPragma   = "_pragma=foreign_keys(1)"
)

var (
	//go:embed schema.sql
	schema string
)

type DB struct {
	*sql.DB
}

func Connect(ctx context.Context, driver, dns string) (*DB, error) {
	db, err := sql.Open(driver, dns)
	if err != nil {
		return nil, fmt.Errorf("openning database driver %q by %q", driver, dns)
	}
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return nil, fmt.Errorf("applying schema: %w", err)
	}
	
	return &DB{db}, nil
}
