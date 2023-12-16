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
	DB          *sql.DB
	execTimeout time.Duration
}

func NewDatabase(connStr string, timeout time.Duration) *Database {
	d, err := sql.Open("sqlite3", connStr)
	if err != nil {
		panic(err.Error())
	}

	return &Database{DB: d, execTimeout: timeout}
}

func (d *Database) Close() error {
	if err := d.DB.Close(); err != nil {
		return fmt.Errorf("error while closing db: %w", err)
	}

	return nil
}

func (d *Database) InitSchema() error {
	ctx, cancel := context.WithTimeout(context.Background(), d.execTimeout)
	defer cancel()

	_, err := d.DB.ExecContext(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to initlize database tables: %w", err)
	}

	return nil
}

func (d *Database) InsertPassword(ctx context.Context, password *Password) error {
	_, err := d.DB.ExecContext(ctx, "INSERT INTO password (key,value) VALUES (?,?)", password.Key, password.Value)
	if err != nil {
		return fmt.Errorf("failed to insert password with key: %s: %w", password.Key, err)
	}

	return nil
}

func (d *Database) GetPasswordByKey(ctx context.Context, key string) (*Password, error) {
	var pass Password

	err := d.DB.QueryRowContext(ctx, "SELECT * FROM `password` WHERE `key`=?", key).Scan(&pass.Key, &pass.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to get password with key: %s: %w", key, err)
	}

	return &pass, nil
}

func (d *Database) BatchInsertPassword(ctx context.Context, passwords []*Password) error {
	tx, err := d.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to batch insert passwords: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			panic(err.Error())
		}
	}()

	if err := tx.Commit(); err != nil {
		panic(err.Error())
	}

	return nil
}
