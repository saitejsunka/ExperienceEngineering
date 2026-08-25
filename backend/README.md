# DBExp Backend Service

This is the backend service for the DBExp project, responsible for handling business logic and communicating with the PostgreSQL database. It is implemented as a **gRPC service** running on Google Cloud Run.

## Architecture

The backend utilizes a **Producer-Consumer** architecture pattern:
- **Producer (gRPC Handlers):** Receives the incoming gRPC request, validates the input, and constructs a query request. It then passes this request to the Consumer.
- **Consumer (Database Workers):** Receives the query request from the Producer, executes the actual database calls against Cloud SQL, and returns the result back to the Producer. 

Currently, the Consumer layer **mocks** the database calls (returning successful responses) until the database schema and connections are fully integrated.

## gRPC APIs

The service is defined by the `proto/dbexp.proto` contract. It exposes the following Remote Procedure Calls (RPCs):

### 1. `CheckHealth`
A standard health check endpoint to verify the service is running.
- **Returns:** `{ status_code: 200, message: "..." }`

### 2. `UpsertPost`
Creates a new post or updates an existing one based on the `post_id`.
- **Request:**
  ```protobuf
  {
    string post_id = 1;
    string title = 2;
    string content = 3;
    string author = 4;
  }
  ```
- **Returns:** `{ status_code: 200, message: "..." }`

### 3. `ReadPost`
Fetches a post by its `post_id`.
- **Request:**
  ```protobuf
  {
    string post_id = 1;
  }
  ```
- **Returns:** `{ status_code: 200, post_id: "...", title: "...", content: "...", author: "..." }`

*(Note: While gRPC natively handles errors via gRPC status codes, we explicitly include HTTP-like `status_code` fields in the response payloads for granular internal tracking).*

## Development & Deployment

- **Containerization:** The backend is compiled and packaged into a minimal, multi-stage Docker container.
- **CI/CD:** A GitHub Actions workflow automatically builds the Docker image, pushes it to Google Artifact Registry, and deploys the new revision.
- **Platform:** Google Cloud Run (Configured to scale between 1 and 2 instances maximum).

## Generating gRPC Stubs (Golang)

To write the actual backend code, you need to "compile" the `.proto` contract into Go code (called "stubs"). These stubs provide the interfaces and structs that your Go code will use.

### 1. Prerequisites & Installations

Since you are on a Mac, you can install the Protocol Buffers compiler (`protoc`) using Homebrew:

```bash
brew install protobuf
```

Next, initialize a Go module in the `backend` directory so Go can track your dependencies:

```bash
go mod init dbexp/backend
```

Then, you need to install the Go-specific plugins for `protoc`. These plugins generate standard Go code and gRPC-specific Go code. Make sure you have Go installed, then run:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

*Note: Ensure your `$(go env GOPATH)/bin` directory is in your system's `$PATH` so the `protoc` compiler can find these plugins.*

### 2. Compiling the Proto Files

Once the tools are installed, open your terminal, navigate to the `backend` directory, and run the following command:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/dbexp.proto
```

**What does this do?**
- `--go_out=.` generates the core Protocol Buffer Go structures (like `CheckHealthRequest`).
- `--go-grpc_out=.` generates the gRPC service code (like the `DBExpServiceClient` and `DBExpServiceServer` interfaces).
- `paths=source_relative` tells the compiler to place the generated `.pb.go` files right next to the original `.proto` file in the `proto/` directory.

After running this command, you will see two new files generated inside the `proto/` folder:
- `dbexp.pb.go`
- `dbexp_grpc.pb.go`
