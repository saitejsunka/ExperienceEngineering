# DBPlay Infrastructure

This directory contains the Infrastructure as Code (IaC) for the DBPlay application, managed via Terraform. 

## Structure & Execution

As per the architectural guidelines, the infrastructure is broken down into highly modular, granular components:

- `main.tf`: The core entry point defining terraform requirements and backends.
- `commons.tf`: Stores shared configurations like variables, providers, and common locals (tags) used across all resources.
- `vpc.tf`: Declares the networking module for the DBExp VPC. Note: Terraform automatically merges all `.tf` files in this directory, so `vpc.tf` and `commons.tf` do not need to be explicitly "included" in `main.tf`.
- `SystemDesign.md`: The detailed architectural documentation, tracking the design and rationale behind every component.

### Modules Added
- `networking`: Creates the foundational Virtual Private Cloud (VPC), currently named `dbexp-vpc`.

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

## Troubleshooting
### Error: `google-github-actions/auth failed...`
If your GitHub Action fails with the error:
> `the GitHub Action workflow must specify exactly one of "workload_identity_provider" or "credentials_json"! If you are specifying input values via GitHub secrets, ensure the secret is being injected into the environment.`

**Cause:** This happens because the `${{ secrets.GCP_CREDENTIALS }}` variable is evaluating to an empty string, meaning GitHub cannot find the secret in your repository.
**Fix:**
1. Go to your GitHub Repository -> **Settings**.
2. Scroll down on the left sidebar to **Secrets and variables** -> **Actions**.
3. Click **New repository secret**.
4. Name it `GCP_CREDENTIALS`.
5. Paste the entire JSON content of your Google Cloud Service Account Key into the "Secret" field and click **Add secret**.
6. (Also make sure to add `GCP_PROJECT_ID` while you are there).
7. Go back to the **Actions** tab and re-run the failed job.
