package store

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"time"
)

//go:embed schema.sql
var schema string

type Password struct {
	Key   string
	Value string
}

type Database struct {
	queries     Queries
	execTimeout time.Duration
}

func NewDatabase(connStr string, timeout time.Duration) *Database {
	d, err := sql.Open("sqlite3", connStr)
	if err != nil {
		panic(err.Error())
	}

	return &Database{queries: d, execTimeout: timeout}
}

func (d *Database) Close() error {
	if err := d.queries.(*sql.DB).Close(); err != nil {
		return fmt.Errorf("error while closing db: %w", err)
	}

	return nil
}

func (d *Database) InitSchema() error {
	ctx, cancel := context.WithTimeout(context.Background(), d.execTimeout)
	defer cancel()

	_, err := d.queries.ExecContext(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to initlize database tables: %w", err)
	}

	return nil
}

func (d *Database) setQueries(q Queries) {
	d.queries = q
}
