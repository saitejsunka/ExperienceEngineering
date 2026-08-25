package service

import (
	"context"

	pb "dbexp/backend/stubs"
)

// CheckHealth implements the CheckHealth gRPC method.
func (s *DBExpServer) CheckHealth(ctx context.Context, req *pb.CheckHealthRequest) (*pb.CheckHealthResponse, error) {
	statusCode, message := s.healthProducer.CheckHealth(ctx)

	return &pb.CheckHealthResponse{
		StatusCode: statusCode,
		Message:    message,
	}, nil
}
