package postrgesql

import (
	"context"
	"database/sql"
	"errors"

	"homework-3/internal/pkg/db"
	"homework-3/internal/pkg/repository"
	"homework-3/internal/pkg/repository/postrgesql/transaction_manager"
)

type PvzRepo struct {
	db          db.DBops
	transaction transaction_manager.TransactionManager
}

func NewPvz(database db.DBops, trManage transaction_manager.TransactionManager) *PvzRepo {
	return &PvzRepo{db: database, transaction: trManage}
}

func (r *PvzRepo) Add(ctx context.Context, pvz *repository.Pvz) (int64, error) {
	var id int64
	err := r.transaction.RunInTransaction(ctx, func(ctxTX context.Context) error {
		err := r.db.ExecQueryRow(ctxTX, `INSERT INTO pvz(name,address,contact) VALUES ($1,$2,$3) RETURNING id;`, pvz.Name, pvz.Address, pvz.Contact).Scan(&id)
		return err
	})
	return id, err
}

func (r *PvzRepo) GetByID(ctx context.Context, id int64) (*repository.Pvz, error) {
	var a repository.Pvz

	err := r.transaction.RunInTransaction(ctx, func(ctxTX context.Context) error {
		err := r.db.Get(ctxTX, &a, "SELECT id,name,address,contact FROM pvz where id=$1", id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return repository.ErrObjectNotFound
			}
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (r *PvzRepo) Update(ctx context.Context, pvz *repository.Pvz) error {
	return r.transaction.RunInTransaction(ctx, func(ctxTX context.Context) error {
		_, err := r.db.Exec(ctxTX, `UPDATE pvz SET name = $2, address = $3, contact = $4 WHERE id = $1;`, pvz.ID, pvz.Name, pvz.Address, pvz.Contact)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *PvzRepo) DeleteByID(ctx context.Context, id int64) error {
	return r.transaction.RunInTransaction(ctx, func(ctxTX context.Context) error {
		_, err := r.db.Exec(ctxTX, `DELETE FROM pvz WHERE id = $1`, id)
		if err != nil {
			return err
		}
		return nil
	})
}
