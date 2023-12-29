package store

import "context"

type Store interface {
	InsertPassword(ctx context.Context, password *Password) error
	BatchInsertPassword(ctx context.Context, passwords []*Password, batchSize int) error
	GetPasswordByKey(ctx context.Context, key string) (*Password, error)
	GetAllPasswords(ctx context.Context) ([]*Password, error)
}
