# DBPlay Infrastructure Pipeline

This directory contains the Infrastructure as Code (IaC) for the DBPlay application, managed strictly via Terraform and executed through GitHub Actions (GitOps).

## Deployment & GitOps Execution

**Important:** Do not execute Terraform commands locally. This infrastructure is strictly governed by CI/CD to prevent state divergence and local configuration issues.

### How the Pipeline Works
1. **Workflow File:** `.github/workflows/terraform.yml`
2. **Concurrency Queue:** The workflow uses a `concurrency` block. If multiple commits are pushed rapidly, GitHub will queue them up and run them one by one, preventing Terraform State Lock crashes.
3. **Remote State:** Terraform tracks the status of your infrastructure in a Google Cloud Storage (GCS) Bucket.
4. **Triggering Execution:**
   - **Feature Branches:** Pushing to any non-main branch will trigger a `terraform plan`. This allows you to safely review intended changes in the GitHub Actions console.
   - **Main Branch:** Pushing directly to `main` triggers `terraform apply -auto-approve`, executing the changes on GCP.

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
