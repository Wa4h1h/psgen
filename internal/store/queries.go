package store

import (
	"context"
	"fmt"
	"sync"
)

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

func (d *Database) BatchInsertPassword(ctx context.Context, passwords []*Password, batchSize int) error {
	tx, err := d.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to batch insert passwords: %w", err)
	}

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
				if err := d.InsertPassword(ctx, pass); err != nil {
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
