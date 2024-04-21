package main

import (
	"context"
	"homework-3/config"
	"homework-3/internal/pkg/db"
	"homework-3/internal/pkg/kafka"
	inmemorycache "homework-3/internal/pkg/repository/in_memory_cache"
	"homework-3/internal/pkg/repository/postrgesql"
	"homework-3/internal/pkg/repository/postrgesql/transaction_manager"
	"homework-3/internal/pkg/repository/redis"
	"homework-3/internal/pkg/router"
	"log"
	"net/http"
	"time"
)

const port = ":9000"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDb(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer database.GetPool(ctx).Close()

	trManager := transaction_manager.NewTransactionManager(database.GetPool(ctx))

	pvzRepo := postrgesql.NewPvz(database, *trManager)
	implemetation := router.Server{Repo: pvzRepo, AuthConfig: config.AuthConfig}

	inmemorycache.NewPvzCache(time.Minute)
	redis.NewRedisPvzCache(config.RedisOpt)

	//kafka
	kafka.InitKafka()
	go kafka.ReadFromKafka(kafka.KafPrCo.Consumer, kafka.Topic)

	http.Handle("/", router.CreateRouter(implemetation))
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
