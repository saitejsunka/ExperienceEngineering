# DBPlay System Design

## Overview
This document tracks the system design, architectural decisions, and techniques implemented for the DBPlay infrastructure.

## Goals
- High Availability & Scalability.
- Modular, low-maintenance Infrastructure as Code (IaC).
- Automated CI/CD execution to remove manual provisioning risk.

## Architecture Decisions

### 1. Infrastructure Execution & Automation (GitOps)
**Problem:** Running Terraform locally risks state divergence, requires local machine configuration (like installing `gcloud`), and lacks execution auditing.
- **Technique (Optimized - Chosen):** GitOps via GitHub Actions. We have completely decoupled execution from local machines. Terraform variables (like GCP Project ID and Authentication credentials) are injected directly into the CI/CD pipeline using GitHub Secrets (`TF_VAR_gcp_project_id`). This ensures a unified, auditable, and automated infrastructure deployment whenever changes are pushed to the `main` branch.

### 2. Infrastructure as Code (IaC) Organization
**Problem:** Keeping Terraform code clean and easily readable.
- **Technique (Highly Optimized - Chosen):** Terraform Modules. The `main.tf` file acts only as an orchestrator and entry point for CI/CD, dynamically accepting variables and calling separate sub-folders for each resource under a `modules/` directory.

### 3. CI/CD Testing Strategy
**Problem:** Validating that infrastructure automation scripts (`terraform.yml`) work correctly before pushing to production.
- **Technique (Optimized - Chosen):** Push-Based Verification. We test the pipeline by pushing infrastructure changes to feature branches. The pipeline is configured to detect any push and run `terraform plan` safely on non-main branches. This allows engineers to verify the execution dry-run in the GitHub Actions console without needing to raise Pull Requests or modify live infrastructure.

### 4. Networking Infrastructure (Multi-Region)
**Problem:** Establishing a secure, high-availability, low-latency foundation for all DBPlay resources serving the entire United States.
- **Technique (Optimized - Chosen):** Global VPC with Dynamic Regional Subnets. We configured the `dbexp-vpc` with `routing_mode = "GLOBAL"`. Instead of hardcoding subnets, the networking module takes a list of target regions (`us-east1`, `us-central1`, `us-west1`) and dynamically provisions a subnet in each region using Terraform's `for_each` and `cidrsubnet` functions. This ensures the app can scale across the US instantly while maintaining strict, centralized network boundaries.

### 5. Database Architecture (CQRS Pattern vs Cost Optimization)
**Problem:** How to structure the database to support heavy global reads while maintaining consistency on writes.
- **Technique (Explicit Control - Chosen):** Cloud SQL PostgreSQL (Primary-Replica Architecture).
  - **Rationale:** While Serverless NoSQL (Firestore) would reduce costs to near $0.00, it abstracts away database networking and manual CQRS (Command Query Responsibility Segregation) routing. To maximize learning outcomes and simulate enterprise backend patterns, we explicitly chose to deploy a traditional SQL setup: 1 Primary (Write) node and 2 Replica (Read) nodes across `us-west1` and `us-east1`.
  - **Cost Management:** By selecting the absolute smallest instance size (`db-f1-micro`) and minimum SSD storage (10 GB) for each of the 3 databases, the total database cost is maintained at ~$32/month. This safely fits inside the $50/month project constraint.

## Current State
- **File Structure Reorganization:** Separated concerns by splitting the root infrastructure into `main.tf` (core terraform config), `commons.tf` (shared variables, locals, providers), and `vpc.tf` (networking module call). This allows for cleaner, domain-specific `.tf` files instead of a bloated `main.tf`.
- Created GitHub Actions Workflow (`.github/workflows/terraform.yml`) as the sole executor.
- Implemented the `networking` module and initialized the `dbexp-vpc` network via `vpc.tf`.
