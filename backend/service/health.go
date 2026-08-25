package main

import (
	"context"

	pb "dbexp/backend/stubs"
)

// CheckHealth implements the CheckHealth gRPC method.
func (s *DBExpServer) CheckHealth(ctx context.Context, req *pb.CheckHealthRequest) (*pb.CheckHealthResponse, error) {
	statusCode, message, err := s.healthProducer.CheckHealth(ctx)
	if err != nil {
		// Even if error is present, we still return the structured response
		return &pb.CheckHealthResponse{
			StatusCode: statusCode,
			Message:    message,
		}, nil
	}

	return &pb.CheckHealthResponse{
		StatusCode: statusCode,
		Message:    message,
	}, nil
}
