package main

import (
	"context"
	"homework-3/config"
	"homework-3/internal/pkg/db"
	"homework-3/internal/pkg/kafka"
	"homework-3/internal/pkg/repository/postrgesql"
	"homework-3/internal/pkg/router"
	"log"
	"net/http"
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

	pvzRepo := postrgesql.NewPvz(database)
	implemetation := router.Server{Repo: pvzRepo, AuthConfig: config.AuthConfig}

	//kafka
	kafka.InitKafka()
	go kafka.ReadFromKafka(kafka.KafPrCo.Consumer, kafka.Topic)

	http.Handle("/", router.CreateRouter(implemetation))
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
