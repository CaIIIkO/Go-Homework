package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"homework-3/config"
	"homework-3/internal/pkg/db"
	"homework-3/internal/pkg/kafka"
	"homework-3/internal/pkg/metrics"
	inmemorycache "homework-3/internal/pkg/repository/in_memory_cache"
	"homework-3/internal/pkg/repository/postrgesql"
	"homework-3/internal/pkg/repository/postrgesql/transaction_manager"
	"homework-3/internal/pkg/repository/redis"
	"homework-3/internal/pkg/server"
	pb "homework-3/internal/pkg/server/pb"
	"homework-3/internal/pkg/tracer"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
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

	inmemorycache.NewPvzCache(time.Minute)
	redis.NewRedisPvzCache(config.RedisOpt)

	kafka.InitKafka()
	go kafka.ReadFromKafka(kafka.KafPrCo.Consumer, kafka.Topic)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := server.Server{Repo: pvzRepo}

	shutdown, err := tracer.InitProvider()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()
	srv.Tracer = otel.Tracer("test-tracer")

	grpcServer := grpc.NewServer()
	pb.RegisterPVZServiceServer(grpcServer, &srv)

	metrics.InitMetrics()

	go http.ListenAndServe(":9091", promhttp.Handler())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
