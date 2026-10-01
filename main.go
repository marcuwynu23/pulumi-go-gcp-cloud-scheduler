package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/cloudscheduler"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		project := config.Get(ctx, "gcp:project")
		if project == "" {
			return fmt.Errorf("gcp:project is required")
		}

		region := config.Get(ctx, "gcp:region")
		if region == "" {
			region = "us-central1"
		}

		jobName := config.Get(ctx, "scheduler:jobName")
		if jobName == "" {
			jobName = "my-scheduler-job"
		}

		jobDescription := config.Get(ctx, "scheduler:jobDescription")
		if jobDescription == "" {
			jobDescription = "Triggers a Cloud Run service periodically"
		}

		schedule := config.Get(ctx, "scheduler:schedule")
		if schedule == "" {
			schedule = "0 */3 * * *"
		}

		timeZone := config.Get(ctx, "scheduler:timeZone")
		if timeZone == "" {
			timeZone = "UTC"
		}

		attemptDeadline := config.Get(ctx, "scheduler:attemptDeadline")
		if attemptDeadline == "" {
			attemptDeadline = "180s"
		}

		httpMethod := config.Get(ctx, "scheduler:httpMethod")
		if httpMethod == "" {
			httpMethod = "POST"
		}

		targetUri := config.Get(ctx, "scheduler:targetUri")
		if targetUri == "" {
			targetUri = "https://your-cloud-run-service-uc.a.run.app"
		}

		body := config.Get(ctx, "scheduler:body")
		if body == "" {
			body = `{"schedule": "trigger"}`
		}

		headersStr := config.Get(ctx, "scheduler:headers")
		headers := map[string]string{"Content-Type": "application/json"}
		if headersStr != "" {
			if err := json.Unmarshal([]byte(headersStr), &headers); err != nil {
				return fmt.Errorf("invalid headers JSON: %w", err)
			}
		}

		oidcSaEmail := config.Get(ctx, "scheduler:oidcSaEmail")
		oidcAudience := config.Get(ctx, "scheduler:oidcAudience")

		// Enable Cloud Scheduler API
		_, err := projects.NewService(ctx, "scheduler-api", &projects.ServiceArgs{
			Service:          pulumi.String("cloudscheduler.googleapis.com"),
			DisableOnDestroy: pulumi.Bool(false),
			Project:          pulumi.String(project),
		})
		if err != nil {
			return err
		}

		// Create Cloud Scheduler job
		bodyEncoded := base64.StdEncoding.EncodeToString([]byte(body))

		httpTarget := &cloudscheduler.JobHttpTargetArgs{
			HttpMethod: pulumi.String(httpMethod),
			Uri:        pulumi.String(targetUri),
			Body:       pulumi.String(bodyEncoded),
			Headers:    pulumi.ToStringMap(headers),
		}

		if oidcSaEmail != "" {
			oidcToken := &cloudscheduler.JobHttpTargetOidcTokenArgs{
				ServiceAccountEmail: pulumi.String(oidcSaEmail),
			}
			if oidcAudience != "" {
				oidcToken.Audience = pulumi.String(oidcAudience)
			}
			httpTarget.OidcToken = oidcToken
		}

		job, err := cloudscheduler.NewJob(ctx, "scheduler-job", &cloudscheduler.JobArgs{
			Name:            pulumi.String(jobName),
			Description:     pulumi.String(jobDescription),
			Schedule:        pulumi.String(schedule),
			TimeZone:        pulumi.String(timeZone),
			AttemptDeadline: pulumi.String(attemptDeadline),
			HttpTarget:      httpTarget,
			Project:         pulumi.String(project),
			Region:          pulumi.String(region),
		})
		if err != nil {
			return err
		}

		// Export outputs
		ctx.Export("jobName", job.Name)
		ctx.Export("jobSchedule", job.Schedule)
		ctx.Export("jobTargetUri", job.HttpTarget.ApplyT(func(ht *cloudscheduler.JobHttpTarget) string {
			if ht != nil {
				return ht.Uri
			}
			return ""
		}).(pulumi.StringOutput))

		return nil
	})
}
