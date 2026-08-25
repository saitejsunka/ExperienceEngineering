package service

import (
	"dbexp/backend/service/producers"
	pb "dbexp/backend/stubs"
)

// DBExpServer implements the gRPC DBExpServiceServer interface.
type DBExpServer struct {
	pb.UnimplementedDBExpServiceServer
	healthProducer producers.HealthProducer
}

// NewDBExpServer initializes a new gRPC service instance with its dependencies.
func NewDBExpServer(hp producers.HealthProducer) *DBExpServer {
	return &DBExpServer{
		healthProducer: hp,
	}
}
