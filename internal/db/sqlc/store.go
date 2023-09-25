package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxStore interface {
	Querier
	CreateTenantTx(ctx context.Context, input CreateTenantInput) (CreateTenantOutput, error)
}

type Store struct {
	*Queries
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) TxStore {
	return &Store{
		pool:    pool,
		Queries: New(pool),
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})

	if err != nil {
		return fmt.Errorf("error begin tx : %w", err)
	}

	q := New(tx)

	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("rollback error : %w", err)
		}
		return fmt.Errorf("error : %w", err)
	}
	return tx.Commit(ctx)

}
