package postrgesql

import (
	"context"
	"database/sql"
	"errors"
	"homework-3/internal/pkg/db"
	"homework-3/internal/pkg/repository"
	inmemorycache "homework-3/internal/pkg/repository/in_memory_cache"
	"homework-3/internal/pkg/repository/postrgesql/transaction_manager"
	"homework-3/internal/pkg/repository/redis"
	"time"
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
		if err == nil {
			inmemorycache.Cache.SetPvz(id, *pvz, 12*time.Hour)
		}
		return err
	})
	return id, err
}

func (r *PvzRepo) GetByID(ctx context.Context, id int64) (*repository.Pvz, error) {
	var a repository.Pvz
	//Заменил inmemorycahe на redis
	// cachedPvz, ok := inmemorycache.Cache.GetPvz(id)
	// if ok {
	// 	a = cachedPvz
	// 	return &a, nil
	// }
	pvz, err := redis.RedisCache.GetPvz(ctx, id)
	if err == nil {
		return pvz, nil
	}

	err = r.transaction.RunInTransaction(ctx, func(ctxTX context.Context) error {
		err := r.db.Get(ctxTX, &a, "SELECT id,name,address,contact FROM pvz where id=$1", id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return repository.ErrObjectNotFound
			}
			return err
		}

		//Заменил inmemorycahe на redis
		//inmemorycache.Cache.SetPvz(id, a, 12*time.Hour)
		redis.RedisCache.SetPvz(ctx, id, a, 12*time.Hour)

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

		_, ok := inmemorycache.Cache.GetPvz(pvz.ID)
		if ok {
			inmemorycache.Cache.SetPvz(pvz.ID, *pvz, 12*time.Hour)
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
		_, ok := inmemorycache.Cache.GetPvz(id)
		if ok {
			inmemorycache.Cache.DeletePvz(id)
		}
		return nil
	})
}
