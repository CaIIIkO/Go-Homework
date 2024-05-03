package server

import (
	"context"
	"errors"
	"time"

	"homework-3/internal/pkg/kafka"
	"homework-3/internal/pkg/metrics"
	"homework-3/internal/pkg/repository"
	pb "homework-3/internal/pkg/server/pb"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedPVZServiceServer
	Tracer trace.Tracer
	Repo   Service
}

type Service interface {
	Add(ctx context.Context, pvz *repository.Pvz) (int64, error)
	GetByID(ctx context.Context, id int64) (*repository.Pvz, error)
	Update(ctx context.Context, pvz *repository.Pvz) error
	DeleteByID(ctx context.Context, id int64) error
}

func (s *Server) AddPvz(ctx context.Context, req *pb.AddPvzRequest) (*pb.AddPvzResponse, error) {
	defer metrics.CustomizedCounterMetricAddPvz.Add(1)

	event := kafka.Event{
		Timestamp: time.Now(),
		Method:    "AddPvz",
		Request:   req.String(),
	}
	if err := kafka.WriteToKafka(event, kafka.KafPrCo.Producer, kafka.Topic); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to produce message: %v", err)
	}

	pvzRepo := &repository.Pvz{
		Name:    req.Name,
		Address: req.Address,
		Contact: req.Contact,
	}

	id, err := s.Repo.Add(ctx, pvzRepo)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add PVZ: %v", err)
	}

	resp := &pb.AddPvzResponse{
		Id: id,
	}

	return resp, nil
}

func (s *Server) GetPvzByID(ctx context.Context, req *pb.GetPvzByIDRequest) (*pb.GetPvzByIDResponse, error) {
	defer metrics.CustomizedCounterMetricGetById.Add(1)
	event := kafka.Event{
		Timestamp: time.Now(),
		Method:    "GetPvzById",
		Request:   req.String(),
	}
	if err := kafka.WriteToKafka(event, kafka.KafPrCo.Producer, kafka.Topic); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to produce message: %v", err)
	}

	id := req.Id

	pvz, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, status.Errorf(codes.NotFound, "PVZ with ID %v not found", id)
		}
		return nil, status.Errorf(codes.Internal, "failed to get PVZ: %v", err)
	}

	resp := &pb.GetPvzByIDResponse{
		Id:      pvz.ID,
		Name:    pvz.Name,
		Address: pvz.Address,
		Contact: pvz.Contact,
	}

	return resp, nil
}

func (s *Server) UpdatePvz(ctx context.Context, req *pb.UpdatePvzRequest) (*pb.UpdatePvzResponse, error) {
	defer metrics.CustomizedCounterMetricUpdateById.Add(1)
	event := kafka.Event{
		Timestamp: time.Now(),
		Method:    "UpdatePvz",
		Request:   req.String(),
	}
	if err := kafka.WriteToKafka(event, kafka.KafPrCo.Producer, kafka.Topic); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to produce message: %v", err)
	}

	pvzRepo := &repository.Pvz{
		ID:      req.Id,
		Name:    req.Name,
		Address: req.Address,
		Contact: req.Contact,
	}

	if err := s.Repo.Update(ctx, pvzRepo); err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, status.Errorf(codes.NotFound, "PVZ with ID %v not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "failed to update PVZ: %v", err)
	}

	return &pb.UpdatePvzResponse{}, nil
}

func (s *Server) DeletePvzByID(ctx context.Context, req *pb.DeletePvzByIDRequest) (*pb.DeletePvzByIDResponse, error) {
	defer metrics.CustomizedCounterMetricDeleteById.Add(1)
	event := kafka.Event{
		Timestamp: time.Now(),
		Method:    "DeletePvzByID",
		Request:   req.String(),
	}
	if err := kafka.WriteToKafka(event, kafka.KafPrCo.Producer, kafka.Topic); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to produce message: %v", err)
	}
	err := s.Repo.DeleteByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, status.Errorf(codes.NotFound, "PVZ with ID %v not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "failed to delete PVZ: %v", err)
	}

	return &pb.DeletePvzByIDResponse{}, nil
}
