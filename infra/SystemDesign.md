# DBPlay System Design

## Overview
This document serves as the architectural source of truth for the DBPlay infrastructure. It tracks the system design, code organization, database selection rationale, structural trade-offs, and future optimization roadmaps.

## System Organization & Code Architecture

### GitOps Execution Model
- **Trade-off:** Manual infrastructure deployment (running `terraform apply` locally) is fast for solo developers but highly error-prone, risks state divergence, and requires local dependencies (`gcloud`, Terraform binaries).
- **Decision:** We strictly enforce a **GitOps pipeline** via GitHub Actions. All infrastructure changes are executed in a clean, ephemeral Ubuntu runner. 
- **Concurrency:** We utilize GitHub Actions `concurrency` queues to ensure that rapid sequential commits do not trigger parallel Terraform runs, which would cause Terraform State Lock crashes.

### Infrastructure as Code (IaC) Organization
The Terraform codebase is highly modularized to avoid a monolithic configuration:
- `main.tf`: Acts solely as the entry point and state backend configuration.
- `commons.tf`: Centralizes shared variables, GCP provider setup, and common tagging locals.
- `vpc.tf`: Instantiates the networking module.
- `db.tf`: Instantiates the database module.
- `modules/`: Contains reusable, self-contained components (Networking and Database) that abstract away the complex resource logic.

---

## Database Architecture

### What Database is Chosen?
- **PostgreSQL 15 (Google Cloud SQL)**

### How Many Databases?
We are deploying a **3-node cluster**:
1. **Primary Node** (`dbexp-write`) located in `us-west1`.
2. **Read Replica 1** (`dbexp-read-west`) located in `us-west1`.
3. **Read Replica 2** (`dbexp-read-east`) located in `us-east1`.

### Database Trade-offs & Rationale
- **Why not Firestore (NoSQL)?** Serverless NoSQL like Firestore is incredibly cost-optimized (near $0.00 for low traffic) and scales globally by default. However, it completely abstracts away networking, query planning, and routing constraints.
- **Why Cloud SQL (PostgreSQL)?** To maximize engineering learning outcomes, we pivoted to a traditional Relational Database Management System (RDBMS). This forces us to manually design the infrastructure layer for **CQRS (Command Query Responsibility Segregation)**. By explicitly separating the write node from the read nodes, the application code must become "infrastructure aware" and manually route mutations to the primary and queries to the replicas.
- **Cost Considerations:** To fit within the project budget constraints ($50/month), all three nodes utilize the absolute smallest instance size (`db-f1-micro`) with minimal SSD storage (10 GB). This brings the total database cost to ~$32/month.

### Schema Management
- **File:** `schema.sql`
- We rely on raw Data Definition Language (DDL) scripts for initialization (currently configuring `posts` and `comments` tables). This avoids the heavy abstraction of ORMs (Object-Relational Mappers) and ensures total control over indexing and foreign key constraints.

---

## Networking Architecture

### Global VPC with Dynamic Subnets
- **VPC Name:** `dbexp-vpc` (Configured with `GLOBAL` routing mode).
- Instead of hardcoding network topologies, the networking module dynamically provisions subnets based on a list of target regions (`us-east1`, `us-central1`, `us-west1`) using Terraform's `for_each` and `cidrsubnet` functions. 

### Security & Private IPs
- **Trade-off:** Exposing databases to the public internet via Public IPs makes it incredibly easy to connect from local development machines, but poses a massive security vulnerability.
- **Decision:** All Cloud SQL instances are configured **strictly with Private IPs**. 
- **Implementation:** We utilize VPC Peering (`servicenetworking.googleapis.com`) to bridge the Google Managed Services tenant project directly into our `dbexp-vpc`. 

---

## Future Roadmap & Optimizations

As the infrastructure scales, we plan to introduce the following optimizations:

1. **Application-Layer CQRS:** Implementing the backend logic to seamlessly route `INSERT`/`UPDATE` queries to the primary node and load-balance `SELECT` queries across the read replicas.
2. **Connection Pooling:** Introducing PgBouncer (either natively via Cloud SQL or hosted on a separate proxy) to manage connection exhaustion limits inherent to `db-f1-micro` instances.
3. **Bastion Host / IAP Tunneling:** Because the databases have no Public IPs, we will need to deploy a micro-VM (Bastion Host) or configure Identity-Aware Proxy (IAP) tunneling to allow local IDEs (like DBeaver or DataGrip) to securely inspect the production database schema.
4. **Automated Migrations:** Moving away from manual `schema.sql` execution and integrating a programmatic migration tool (like Flyway, Liquibase, or Go-Migrate) directly into the CI/CD pipeline.
