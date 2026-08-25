# DBExp Backend Service

This is the backend service for the DBExp project, responsible for handling business logic and communicating with the PostgreSQL database. It is implemented as a **gRPC service** using Go.

## Architecture & Code Structure

The backend utilizes a decoupled **Producer-Consumer** architecture pattern mapped strictly to the folder structure:
- **`service/producers/`**: Receives the incoming gRPC request, validates the input, enforces business rules, and delegates raw queries to the Consumer layer.
- **`service/consumers/`**: Handles all infrastructure concerns, primarily executing database queries against Cloud SQL via `database/sql` and `lib/pq`.

### Service Initialization
The entry point of the service is `main.go`. It follows a strict initialization order:
1. **Telemetry**: Initializes global GCP Cloud Logging first.
2. **Configuration**: Pulls database credentials securely from GCP Secret Manager (using the Telemetry logger to capture failures).
3. **Dependency Injection**: Wires the Consumer and Producer structs together and mounts them to the gRPC Server.

## Development & Local Setup

### 1. Prerequisites
You will need the Protocol Buffers compiler (`protoc`) and the Go gRPC plugins:

```bash
# Mac Installation
brew install protobuf

# Install Go Plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

*Ensure your `$(go env GOPATH)/bin` is in your system's `$PATH`.*

### 2. Compiling Proto Contracts
To generate the Go stubs from the `dbexp.proto` definition, run this from the `backend/` directory:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
protoc -I contracts --go_out=stubs --go_opt=paths=source_relative \
       --go-grpc_out=stubs --go-grpc_opt=paths=source_relative \
       dbexp.proto
```
This updates `stubs/dbexp.pb.go` and `stubs/dbexp_grpc.pb.go`.

### 3. Running the Service
The service requires two environment variables to successfully fetch its configuration from Secret Manager:
- `GCP_PROJECT_ID`: Your GCP Project ID.
- `DB_SECRET_NAME`: The name of the secret in Secret Manager (e.g., `dbexp-backend-config`).

To run locally:
```bash
export GCP_PROJECT_ID="your-project-id"
export DB_SECRET_NAME="dbexp-backend-config"

go mod tidy
go build -o server ./service
./server
```

*(Note: For infrastructure, deployment, or architecture decisions, see `infra/README.md` and `infra/SystemDesign.md` respectively).*
