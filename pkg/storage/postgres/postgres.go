package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Store
type Store struct {
	Db *pgxpool.Pool
}

// pool constructor
func NewDb(ctx context.Context, constr string) (*pgxpool.Pool, error) {

	for {
		_, err := pgxpool.Connect(ctx, constr)
		if err == nil {
			break
		}
	}
	db, err := pgxpool.Connect(ctx, constr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Store constructor
func New(ctx context.Context, constr string) (*Store, error) {

	db, err := NewDb(ctx, constr)
	if err != nil {
		return nil, err
	}

	s := &Store{
		Db: db,
	}

	return s, nil
}
