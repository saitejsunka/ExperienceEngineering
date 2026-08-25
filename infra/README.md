# DBPlay Infrastructure Pipeline

This directory contains the Infrastructure as Code (IaC) for the DBPlay application, managed strictly via Terraform and executed through GitHub Actions (GitOps).

## Deployment & GitOps Execution

**Important:** Do not execute Terraform commands locally. This infrastructure is strictly governed by CI/CD to prevent state divergence and local configuration issues.

### How the GitOps Pipelines Work
We utilize a strict separation of concerns for CI/CD deployments:

1. **Static Infrastructure (`.github/workflows/terraform.yml`)**
   - Triggers on changes to `infra/`.
   - Provisions all fixed GCP resources (VPC, Cloud SQL, Secret Manager, Artifact Registry).
   - Uses remote state stored in a Google Cloud Storage (GCS) Bucket.

2. **Dynamic Application (`.github/workflows/deploy.yml`)**
   - Triggers on changes to `backend/`.
   - Builds the Docker container, authenticates with Artifact Registry, and pushes the image.
   - Executes `gcloud run deploy` to safely deploy the new container to Cloud Run with Direct VPC Egress enabled.

## Prerequisites & GitHub Secrets

Before the GitHub Actions pipeline can run, the repository must have the following **Repository Secrets** configured (`Settings` -> `Secrets and variables` -> `Actions`):
- `GCP_PROJECT_ID`: The exact string of your GCP Project ID.
- `GCP_CREDENTIALS`: The raw JSON content of a GCP Service Account Key.

### Required Service Account Roles
The injected Service Account must possess the **Editor** role (or the equivalent granular networking, cloud sql, and secret manager admin roles) to successfully provision the databases and network.

## Troubleshooting CI/CD Failures

### Error: `google-github-actions/auth failed...`
**Cause:** The `GCP_CREDENTIALS` secret is missing or empty.
**Fix:** Ensure the secret is created in the GitHub Repo settings and contains the valid JSON key.

### Error: `Error 403: Required 'compute.networks.create' permission...`
**Cause:** The Service Account lacks sufficient permissions.
**Fix:** Grant the Service Account the **Editor** role in GCP IAM.

### Error: `Error code 9: Failed to delete connection; Producer services are still using this connection`
**Cause:** GCP caching bug when destroying a VPC that previously held Cloud SQL instances.
**Fix:** Manually delete the VPC and Subnets in Google Cloud Shell (`gcloud compute networks delete dbexp-vpc`), then re-run the Action.

*(Note: For architectural decisions or application code setup, refer to `SystemDesign.md` and `../backend/README.md` respectively).*
