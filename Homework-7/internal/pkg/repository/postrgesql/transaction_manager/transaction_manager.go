package transaction_manager

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/multierr"
)

type contextKey string

const key contextKey = "transaction"

type TransactionManager struct {
	pool *pgxpool.Pool
}

type TransactionManagerI interface {
	RunInTransaction(ctx context.Context, f func(ctxTX context.Context) error) error
}

func NewTransactionManager(pool *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{
		pool: pool,
	}
}

func (t *TransactionManager) RunInTransaction(ctx context.Context, f func(ctxTX context.Context) error) error {
	tx, err := t.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite})
	if err != nil {
		return err
	}

	if err := f(context.WithValue(ctx, key, tx)); err != nil {
		errRollback := tx.Rollback(ctx)
		return multierr.Combine(err, errRollback)
	}

	if err := tx.Commit(ctx); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	return nil
}
