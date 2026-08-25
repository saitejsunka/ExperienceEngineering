package main

import (
	"context"
	"log"
	"net"
	"os"

	"dbexp/backend/service/configs"
	"dbexp/backend/service/consumers"
	"dbexp/backend/service/observability"
	"dbexp/backend/service/producers"
	pb "dbexp/backend/stubs"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	projectID := os.Getenv("GCP_PROJECT_ID")
	secretName := os.Getenv("DB_SECRET_NAME")

	// 0. Initialize Telemetry FIRST
	telemetry, err := observability.InitializeObservabilityAndTelemetry(ctx, projectID)
	if err != nil {
		log.Fatalf("failed to init telemetry: %v", err)
	}
	defer telemetry.Close()

	// 1. Load Configuration using Telemetry to log failure
	appConfig, err := configs.LoadConfigFromSecretManager(ctx, projectID, secretName)
	if err != nil {
		telemetry.LogError("Failed to load configuration from Secret Manager", err)
		return
	}

	// 2. Initialize the Consumers (Layer responsible for hitting the DB)
	healthConsumer := consumers.InitializeHealthConsumer(appConfig, telemetry)

	// 3. Initialize the Producers (Layer responsible for business logic & queries)
	healthProducer := producers.InitializeHealthProducer(healthConsumer)

	// 4. Initialize the gRPC Service implementation
	grpcHandler := NewDBExpServer(healthProducer)

	// 5. Start standard gRPC server on 8050
	port := os.Getenv("PORT")
	if port == "" {
		port = "8050"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		telemetry.LogError("Failed to listen on port "+port, err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterDBExpServiceServer(s, grpcHandler)

	// Enable Server Reflection for tools like grpcurl/Postman
	reflection.Register(s)

	log.Printf("DBExpService gRPC server listening on port %s...\n", port)
	telemetry.LogInfo("Starting DBExpService", nil)

	if err := s.Serve(lis); err != nil {
		telemetry.LogError("Failed to serve gRPC", err)
		return
	}
}
