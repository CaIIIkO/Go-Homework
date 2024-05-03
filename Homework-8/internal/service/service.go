package service

import (
	"context"
	"time"

	"homework-3/internal/pkg/repository"
	inmemorycache "homework-3/internal/pkg/repository/in_memory_cache"
	"homework-3/internal/pkg/repository/postrgesql"
	"homework-3/internal/pkg/repository/redis"
)

const ttl time.Duration = 12 * time.Hour

type Service struct {
	Repo postrgesql.PvzRepo
}

func (s *Service) Add(ctx context.Context, pvz *repository.Pvz) (int64, error) {
	id, err := s.Repo.Add(ctx, pvz)
	if err == nil {
		inmemorycache.Cache.SetPvz(id, *pvz, ttl)
	}
	return id, err
}

func (s *Service) GetByID(ctx context.Context, id int64) (*repository.Pvz, error) {
	var pvz *repository.Pvz
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

	pvz, err = s.Repo.GetByID(ctx, id)
	if err == nil {
		//Заменил inmemorycahe на redis
		//inmemorycache.Cache.SetPvz(id, a, ttl)
		redis.RedisCache.SetPvz(ctx, id, *pvz, ttl)
	}
	if err != nil {
		return nil, err
	}
	return pvz, nil
}

func (s *Service) Update(ctx context.Context, pvz *repository.Pvz) error {
	_, ok := inmemorycache.Cache.GetPvz(pvz.ID)
	if ok {
		inmemorycache.Cache.SetPvz(pvz.ID, *pvz, ttl)
	}
	return s.Repo.Update(ctx, pvz)
}

func (s *Service) DeleteByID(ctx context.Context, id int64) error {
	_, ok := inmemorycache.Cache.GetPvz(id)
	if ok {
		inmemorycache.Cache.DeletePvz(id)
	}
	return s.Repo.DeleteByID(ctx, id)
}
