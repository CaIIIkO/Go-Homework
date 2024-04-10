package postrgesql

import (
	"context"
	"database/sql"
	"errors"
	"homework-3/internal/pkg/db"
	"homework-3/internal/pkg/repository"
)

type PvzRepo struct {
	db *db.Database
}

func NewPvz(database *db.Database) *PvzRepo {
	return &PvzRepo{db: database}
}

func (r *PvzRepo) Add(ctx context.Context, pvz *repository.Pvz) (int64, error) {
	var id int64
	err := r.db.ExecQueryRow(ctx, `INSERT INTO pvz(name,address,contact) VALUES ($1,$2,$3) RETURNING id;`, pvz.Name, pvz.Address, pvz.Contact).Scan(&id)
	return id, err
}

func (r *PvzRepo) GetByID(ctx context.Context, id int64) (*repository.Pvz, error) {
	var a repository.Pvz
	err := r.db.Get(ctx, &a, "SELECT id,name,address,contact FROM pvz where id=$1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (r *PvzRepo) Update(ctx context.Context, pvz *repository.Pvz) error {
	_, err := r.db.Exec(ctx, `UPDATE pvz SET name = $2, address = $3, contact = $4 WHERE id = $1;`, pvz.ID, pvz.Name, pvz.Address, pvz.Contact)
	return err
}

func (r *PvzRepo) DeleteByID(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM pvz WHERE id = $1`, id)
	return err
}
