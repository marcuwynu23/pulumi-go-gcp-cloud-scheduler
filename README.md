# Pulumi Go GCP Cloud Scheduler

This Pulumi project provisions a Google Cloud Scheduler job that triggers an HTTP endpoint on a cron schedule.

## Resources Created

- **Cloud Scheduler Job**: HTTP-targeted job with configurable schedule, headers, and OIDC authentication
- **API Enablement**: `cloudscheduler.googleapis.com`

## Prerequisites

- Pulumi CLI installed
- Go 1.21+
- GCP credentials configured
- A GCP project with billing enabled

## Configuration

Copy `Pulumi.dev.yaml.example` to `Pulumi.dev.yaml` and update values:

```yaml
config:
  gcp:project: your-project-id
  gcp:region: us-central1
  scheduler:jobName: my-scheduler-job
  scheduler:jobDescription: Triggers a Cloud Run service periodically
  scheduler:schedule: "0 */3 * * *"
  scheduler:timeZone: UTC
  scheduler:attemptDeadline: 180s
  scheduler:httpMethod: POST
  scheduler:targetUri: https://your-cloud-run-service-uc.a.run.app
  scheduler:body: '{"schedule": "trigger"}'
  scheduler:headers: '{"Content-Type": "application/json"}'
  scheduler:oidcSaEmail: ""
  scheduler:oidcAudience: ""
```

### Config Options

| Key | Description | Default |
|-----|-------------|---------|
| `gcp:project` | GCP Project ID | **Required** |
| `gcp:region` | GCP Region (Free Tier: us-central1, us-east1, us-west1) | `us-central1` |
| `scheduler:jobName` | Job name | `my-scheduler-job` |
| `scheduler:jobDescription` | Job description | `Triggers a Cloud Run service periodically` |
| `scheduler:schedule` | Cron schedule | `0 */3 * * *` |
| `scheduler:timeZone` | Time zone | `UTC` |
| `scheduler:attemptDeadline` | Execution deadline | `180s` |
| `scheduler:httpMethod` | HTTP method | `POST` |
| `scheduler:targetUri` | Target HTTP endpoint | `https://your-cloud-run-service-uc.a.run.app` |
| `scheduler:body` | HTTP request body | `{"schedule": "trigger"}` |
| `scheduler:headers` | HTTP headers (JSON) | `{"Content-Type": "application/json"}` |
| `scheduler:oidcSaEmail` | OIDC service account email | `""` |
| `scheduler:oidcAudience` | OIDC audience | `""` |

## Deployment

```bash
# Install dependencies
go mod tidy

# Preview changes
pulumi preview

# Deploy
pulumi up
```

## Outputs

| Output | Description |
|--------|-------------|
| `jobName` | The name of the Cloud Scheduler job |
| `jobSchedule` | The cron schedule of the job |
| `jobTargetUri` | The target URI of the job |

## OIDC Authentication

To authenticate with Cloud Run or other OIDC-protected endpoints:
1. Create a service account with appropriate roles
2. Set `scheduler:oidcSaEmail` to the service account email
3. Set `scheduler:oidcAudience` to the target service URL (optional, defaults to targetUri)

## Free Tier Notes

Cloud Scheduler free tier includes:
- 3 jobs/month
- 100,000 job executions/month

See [Cloud Scheduler Pricing](https://cloud.google.com/scheduler/pricing) for details.

## Cleanup

```bash
pulumi destroy
```