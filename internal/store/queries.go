package store

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
)

type DbCommon interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type Queries struct {
	dq DbCommon
}

func (q *Queries) withTx(tx *sql.Tx) *Queries {
	return &Queries{dq: tx}
}

func (q *Queries) InsertPassword(ctx context.Context, password *Password) error {
	_, err := q.dq.ExecContext(ctx, "INSERT INTO password (key,value) VALUES (?,?)", password.Key, password.Value)
	if err != nil {
		return fmt.Errorf("failed to insert password with key: %s: %w", password.Key, err)
	}

	return nil
}

func (q *Queries) GetPasswordByKey(ctx context.Context, key string) (*Password, error) {
	var pass Password

	err := q.dq.QueryRowContext(ctx, "SELECT * FROM `password` WHERE `key`=?", key).Scan(&pass.Key, &pass.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to get password with key: %s: %w", key, err)
	}

	return &pass, nil
}

func (q *Queries) GetAllPasswords(ctx context.Context) ([]*Password, error) {
	rows, err := q.dq.QueryContext(ctx, "SELECT * FROM `password`")
	if err != nil {
		return nil, fmt.Errorf("failed to get all passwords: %w", err)
	}
	defer rows.Close()

	passwords := make([]*Password, 0)

	for rows.Next() {
		var pass Password
		if err := rows.Scan(&pass.Key, &pass.Value); err != nil {
			return nil, fmt.Errorf("unmarshlling password error: %w", err)
		}

		passwords = append(passwords, &pass)
	}

	return passwords, nil
}

func (q *Queries) BatchInsertPassword(ctx context.Context, passwords []*Password, batchSize int) error {
	tx, err := q.dq.(*sql.DB).BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to batch insert passwords: %w", err)
	}

	qu := q.withTx(tx)

	var wg sync.WaitGroup

	errTx := make(chan error)
	done := make(chan struct{})

	for i := 0; i < len(passwords); i += batchSize {
		passes := make([]*Password, 0)
		for j := i; j <= i+batchSize-1 && j < len(passwords); j++ {
			passes = append(passes, passwords[j])
		}
		wg.Add(1)

		go func(errTx chan error, passes []*Password) {
			defer wg.Done()

			for _, pass := range passes {
				if err := qu.InsertPassword(ctx, pass); err != nil {
					errTx <- fmt.Errorf("batch insert failed: %w", err)
				}
			}
		}(errTx, passes)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case errTxVal := <-errTx:
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("%w, rollback failed: %w", errTxVal, err)
		}

		return errTxVal
	case <-done:
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("tx commit failed: %w", err)
		}

		return nil
	}
}
