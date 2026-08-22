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

## Remote State Configuration (CRITICAL)
Terraform tracks the status of your infrastructure using a "State File". Because GitHub Actions destroys its runner after every execution, **we must store this file in a Google Cloud Storage (GCS) Bucket**. 

Before you push any infrastructure code, you must:
1. Go to your [Google Cloud Storage Console](https://console.cloud.google.com/storage).
2. Click **Create Bucket**.
3. Name it something globally unique (e.g., `dbplay-<your-initials>-tf-state`).
4. Click **Create**.
5. Open `infra/main.tf` in this repository and replace `REPLACE_WITH_YOUR_GLOBALLY_UNIQUE_BUCKET_NAME` with the exact bucket name you just created.

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

### Error: `Error 403: Required 'compute.networks.create' permission...`
If your Terraform Apply fails with a 403 Forbidden error, it means the Google Cloud Service Account you provided in your `GCP_CREDENTIALS` secret does not have sufficient permissions to create resources.

**Fix:**
1. Go to the [Google Cloud Console](https://console.cloud.google.com/).
2. Navigate to **IAM & Admin** > **IAM**.
3. Find the Service Account associated with your GitHub Actions.
4. Click the **Edit (pencil icon)** next to it.
5. Click **Add Another Role**.
6. Since this Service Account is your Terraform executor, it needs broad access to provision resources. Grant it the **Editor** role (or if you want to be extremely strict, the **Compute Network Admin** role just for this VPC).
7. Click **Save** and re-run your GitHub Action.
