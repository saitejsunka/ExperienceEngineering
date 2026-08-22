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

## Current State
- Set up foundational Terraform entry point (`main.tf`) ready to consume variables from CI/CD.
- Created GitHub Actions Workflow (`.github/workflows/terraform.yml`) as the sole executor.

### 3. CI/CD Testing Strategy
**Problem:** Validating that infrastructure automation scripts (`terraform.yml`) work correctly before pushing to production.
- **Technique (Optimized - Chosen):** Push-Based Verification. We test the pipeline by pushing infrastructure changes to feature branches. The pipeline is configured to detect any push and run `terraform plan` safely on non-main branches. This allows engineers to verify the execution dry-run in the GitHub Actions console without needing to raise Pull Requests or modify live infrastructure.
