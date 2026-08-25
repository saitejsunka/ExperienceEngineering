package main

import (
	"context"
	"log"
	"net"

	"dbexp/backend/service"
	"dbexp/backend/service/configs"
	"dbexp/backend/service/consumers"
	"dbexp/backend/service/observability"
	"dbexp/backend/service/producers"
	pb "dbexp/backend/stubs"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 0. Load Configuration
	appConfig, err := configs.LoadConfig("service/configs/database.config")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 0. Initialize Telemetry
	ctx := context.Background()
	telemetry, err := observability.InitTelemetry(ctx, appConfig.GCPProjectID)
	if err != nil {
		log.Fatalf("failed to init telemetry: %v", err)
	}
	defer telemetry.Close()

	// 1. Initialize the Consumers (Layer responsible for hitting the DB)
	healthConsumer := consumers.NewMockHealthConsumer(appConfig, telemetry)

	// 2. Initialize the Producers (Layer responsible for business logic & queries)
	healthProducer := producers.NewHealthProducer(healthConsumer)

	// 3. Initialize the gRPC Service implementation
	grpcHandler := service.NewDBExpServer(healthProducer)

	// 4. Start standard gRPC server
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterDBExpServiceServer(s, grpcHandler)

	// Enable Server Reflection for tools like grpcurl/Postman
	reflection.Register(s)

	log.Println("DBExpService gRPC server listening on port 8080...")
	telemetry.LogInfo("Starting DBExpService", nil)
	
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
