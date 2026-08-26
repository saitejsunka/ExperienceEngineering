# DBPlay System Design Architecture

This document serves as the architectural source of truth for the DBPlay infrastructure. It tracks the system design, database selection rationale, structural trade-offs, and future optimization roadmaps.

## 1. Application Layer Architecture

### Containerization (Cloud Run)
The application is deployed as a stateless Docker container to **Google Cloud Run**. This provides:
- **Autoscaling:** Scales up automatically based on traffic volume, and scales down to zero when idle to minimize costs.
- **Direct VPC Egress:** The Cloud Run service is injected directly into our private `dbexp-vpc` network, allowing it to communicate securely with our private Cloud SQL database without needing Serverless VPC Access Connectors or public IPs.

### Producer-Consumer Decoupling
To enforce strict separation of concerns, the backend codebase isolates business rules from infrastructure:
- **Producers**: Enforce application logic (e.g., verifying a post author matches the session user) and dictate *what* data is needed.
- **Consumers**: Act as pure infrastructure adapters. They execute SQL against the database, unaware of the broader business context.

### Bootstrapping & Secret Manager Integration
- **Trade-off:** Hardcoding credentials or using raw environment variables is an immediate security risk in production.
- **Decision:** The application strictly retrieves database credentials (IP, Username, Password) from **GCP Secret Manager** via Cloud Run's native secret integration as environment variables.
- **Implementation:** The service leverages a two-phase initialization. It initializes Telemetry *first*, allowing any failure in the Secret Manager network fetch to be securely logged to GCP Cloud Logging before the application exits.

## 2. Database Architecture

### What Database is Chosen?
- **PostgreSQL 15 (Google Cloud SQL)**

### Cluster Topology
We are deploying a **3-node cluster**:
1. **Primary Node** (`dbexp-write`) located in `us-west1`.
2. **Read Replica 1** (`dbexp-read-west`) located in `us-west1`.
3. **Read Replica 2** (`dbexp-read-east`) located in `us-east1`.

### Rationale & Trade-offs
- **Why Cloud SQL over Firestore?** Serverless NoSQL (Firestore) is cost-optimized and scales globally by default. However, to maximize engineering learning outcomes, we pivoted to a traditional Relational Database (RDBMS). This forces us to manually design the infrastructure layer for **CQRS (Command Query Responsibility Segregation)**, making the application code infrastructure-aware.
- **Cost Considerations:** To fit within budget constraints ($50/month), all three nodes utilize the smallest instance size (`db-f1-micro`) with minimal SSD storage (10 GB).

## 3. Networking & Security Architecture

### Global VPC with Dynamic Subnets
- **VPC Name:** `dbexp-vpc` (Configured with `GLOBAL` routing mode).
- Instead of hardcoding network topologies, the networking module dynamically provisions subnets based on target regions (`us-east1`, `us-central1`, `us-west1`) using Terraform.

### Strict Private IPs
- **Trade-off:** Exposing databases to the public internet via Public IPs makes it easy to connect from local development machines, but poses a massive security vulnerability.
- **Decision:** All Cloud SQL instances are configured **strictly with Private IPs**. We utilize VPC Peering (`servicenetworking.googleapis.com`) to bridge the Google Managed Services tenant project directly into our `dbexp-vpc`.

## 4. Multi-Region Database Connectivity (Read/Write Splitting)

Deploying Cloud Run across multiple regions (`us-west1`, `us-east1`) requires two distinct connectivity links to function properly:

### A. The Network Link (VPC Egress)
- **Automatic Routing:** We use **Direct VPC Egress** (`--vpc-egress all-traffic`) to attach each Cloud Run instance to its local regional subnet (e.g., `us-east1` Cloud Run attaches to `dbexp-vpc-us-east1-subnet`).
- **Global Backbone:** Because `dbexp-vpc` is a Global VPC, an instance in `us-east1` can instantly route traffic to a private IP in `us-west1` over Google's internal fiber network. The network routing is fully automatic.

### B. The Application Link (CQRS Pattern)
- **Manual Implementation:** Network routing does not imply query intelligence. By default, a backend application will send all queries (Reads and Writes) to a single database URL, completely ignoring read replicas.
- **CQRS Requirement:** To utilize regional read replicas, the Go application must be explicitly programmed to maintain **two separate connection pools**:
  1. `dbWrite`: Points exclusively to the `us-west1` Primary Node for all `INSERT`/`UPDATE`/`DELETE` operations.
  2. `dbRead`: Points to the geographically closest Read Replica (e.g., `dbexp-read-east` for the `us-east1` instance) for all `SELECT` operations.
- **Deployment Strategy:** The CI/CD pipeline injects region-specific IPs into the containers as environment variables (e.g., `DB_WRITE_HOST` and `DB_READ_HOST`) during the `gcloud run deploy` step.

## 5. Future Roadmap & Optimizations

1. **Application-Layer CQRS:** Implementing the backend logic to seamlessly route `INSERT`/`UPDATE` queries to the primary node and load-balance `SELECT` queries across the read replicas.
2. **Connection Pooling:** Introducing PgBouncer to manage connection exhaustion limits inherent to `db-f1-micro` instances.
3. **Bastion Host / IAP Tunneling:** Deploying a micro-VM (Bastion Host) or configuring Identity-Aware Proxy (IAP) tunneling to allow local IDEs to securely inspect the private database schema.
4. **Automated Migrations:** Integrating a programmatic migration tool (like Flyway or Go-Migrate) directly into the CI/CD pipeline, replacing manual `schema.sql` execution.

*(Note: For application setup instructions, see `../backend/README.md`. For CI/CD deployment instructions, see `README.md`).*
