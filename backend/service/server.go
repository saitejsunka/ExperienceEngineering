package main

import (
	"dbexp/backend/service/producers/source"
	pb "dbexp/backend/stubs"
)

// DBExpServer implements the gRPC DBExpServiceServer interface.
type DBExpServer struct {
	pb.UnimplementedDBExpServiceServer
	healthProducer source.HealthProducer
}

// NewDBExpServer initializes a new gRPC service instance with its dependencies.
func NewDBExpServer(hp source.HealthProducer) *DBExpServer {
	return &DBExpServer{
		healthProducer: hp,
	}
}
