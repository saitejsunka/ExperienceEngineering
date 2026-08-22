# DBPlay Infrastructure

This directory contains the Infrastructure as Code (IaC) for the DBPlay application, managed via Terraform. 

## Structure & Execution

As per the architectural guidelines, the infrastructure is broken down into highly modular, granular components:

- `main.tf`: The primary orchestrator. It receives injected environment variables (such as the GCP Project ID) from GitHub Actions and calls individual modules for each resource.
- `SystemDesign.md`: The detailed architectural documentation, tracking the design and rationale behind every component.

## Execution via GitHub Actions (GitOps)
**Important:** Do not execute Terraform commands locally. 

This infrastructure is strictly governed by CI/CD through GitHub Actions. 
- **Workflow File:** `.github/workflows/terraform.yml`
- **Environment Variables:** Credentials and configuration (like `GCP_PROJECT_ID` and `GCP_CREDENTIALS`) must be stored as GitHub Repository Secrets.
- **Process:** Whenever code in this `infra/` folder is pushed to any branch, GitHub Actions will trigger automatically. Pushes to non-main branches will run a `terraform plan` for review. Pushes to the `main` branch will automatically run `terraform apply`.

## Testing the CI/CD Pipeline
To test whether the `terraform.yml` configuration works:
1. **GitHub Secrets:** Ensure that `GCP_PROJECT_ID` and `GCP_CREDENTIALS` are configured as secrets in your GitHub repository settings.
2. **Push to any Branch:** Commit your changes and push them. This triggers the GitHub Action to run `terraform plan` (if on a feature branch) or `terraform apply` (if on `main`).
3. **Observe the Actions Tab:** Go to the "Actions" tab in GitHub to review the output of the workflow. You will see exactly what Terraform plans to create without actually provisioning resources (unless you pushed directly to `main`).
