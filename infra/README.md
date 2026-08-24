# DBPlay Infrastructure

This directory contains the Infrastructure as Code (IaC) for the DBPlay application, managed strictly via Terraform and executed through GitHub Actions (GitOps).

## Prerequisites & Setup

Before executing this infrastructure pipeline, the following setup must be completed in your Google Cloud Platform (GCP) project and your GitHub Repository.

### 1. Service Account Roles
You must create a Service Account in GCP for GitHub Actions to use. This Service Account requires the following roles to successfully provision the infrastructure:
- **Editor** (Broad permission for ease of learning/development)
- *Alternatively, for strict least-privilege:*
  - Compute Network Admin (For VPC and Subnets)
  - Cloud SQL Admin (For Database Instances)
  - Storage Admin (For the Terraform State Bucket)
  - Service Usage Admin (To enable required APIs)

### 2. Enabled GCP APIs
The GitHub Actions pipeline will attempt to automatically enable these APIs, but they are required for the infrastructure to function:
- `cloudresourcemanager.googleapis.com` (To manage project APIs)
- `compute.googleapis.com` (For VPC and networking)
- `servicenetworking.googleapis.com` (For VPC Peering and Private IPs)
- `sqladmin.googleapis.com` (For Cloud SQL databases)

### 3. GitHub Secrets (Environment Variables)
Your GitHub repository must have the following **Repository Secrets** configured (`Settings` -> `Secrets and variables` -> `Actions`):
- `GCP_PROJECT_ID`: The exact string of your GCP Project ID (e.g., `portfolio-506303`).
- `GCP_CREDENTIALS`: The raw JSON content of your GCP Service Account Key.

---

## Execution via GitHub Actions (GitOps)

**Important:** Do not execute Terraform commands locally. This infrastructure is strictly governed by CI/CD to prevent state divergence and local configuration issues.

### How it Works
1. **Workflow File:** `.github/workflows/terraform.yml`
2. **Concurrency Queue:** The workflow is configured with a `concurrency` block. If multiple commits are pushed rapidly, GitHub will queue them up and run them one by one. This prevents Terraform State Lock crashes and corrupted infrastructure.
3. **Remote State:** Terraform tracks the status of your infrastructure using a "State File". The pipeline automatically creates a Google Cloud Storage (GCS) Bucket (e.g., `gs://dbplay-<project-id>-tf-state`) and configures Terraform to store the state there.
4. **Triggering Execution:**
   - **Feature Branches:** Pushing to any non-main branch will trigger a `terraform plan`. This allows you to safely review what Terraform *intends* to do in the GitHub Actions console without modifying live infrastructure.
   - **Main Branch:** Pushing directly to `main` (or merging a PR) triggers `terraform apply -auto-approve`, which executes the changes on GCP.

## Troubleshooting

### Error: `google-github-actions/auth failed...`
**Cause:** The `GCP_CREDENTIALS` secret is missing or empty.
**Fix:** Ensure the secret is created in the GitHub Repo settings and contains the valid JSON key.

### Error: `Error 403: Required 'compute.networks.create' permission...`
**Cause:** The Service Account in GCP lacks sufficient permissions.
**Fix:** Grant the Service Account the **Editor** role in GCP IAM.

### Error: `Error code 9: Failed to delete connection; Producer services are still using this connection`
**Cause:** GCP caching bug when destroying a VPC that previously held Cloud SQL instances.
**Fix:** Manually delete the VPC and Subnets in Google Cloud Shell (`gcloud compute networks delete dbexp-vpc`), then re-run the GitHub Action.
