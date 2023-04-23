package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/fatih/color"
	"github.com/hetznercloud/hcloud-go/hcloud"

	"github.com/bonsai-oss/hetzner-machine-provider/internal/actions"
	"github.com/bonsai-oss/hetzner-machine-provider/internal/helper"
)

var version = "dev"

type application struct {
	hcloudToken string
	jobID       string

	execScriptPath string
	execStageName  string

	hcloudClient *hcloud.Client

	vmParams actions.VMParams

	resourceNamePrefix string
	prepareOptions     actions.PrepareOptions
}

func (a *application) prepare(_ *kingpin.ParseContext) error {
	color.Green("🚀 Preparing environment")
	a.prepareOptions.JobID = a.jobID
	return actions.Prepare(a.hcloudClient, a.prepareOptions, a.vmParams)
}

func (a *application) cleanup(_ *kingpin.ParseContext) error {
	color.Green("🧼 Cleaning up resources")
	return actions.Cleanup(a.hcloudClient, a.jobID)
}

func (a *application) exec(_ *kingpin.ParseContext) error {
	return actions.Exec(a.execScriptPath, a.execStageName)
}

func (a *application) configure(_ *kingpin.ParseContext) error {
	// see https://docs.gitlab.com/runner/executors/custom.html#config for more information
	data := map[string]any{
		"driver": map[string]any{
			"name":    "hetzner-machine-provider",
			"version": version,
		},
		"hostname": "🔭 hetzner-machine-provider @ " + os.Getenv("HOSTNAME"),
	}

	return json.NewEncoder(os.Stdout).Encode(data)
}

func (a *application) prepareClient(_ *kingpin.ParseContext) error {
	a.hcloudClient = hcloud.NewClient(hcloud.WithToken(a.hcloudToken), hcloud.WithApplication("hmp", version))
	return nil
}

func main() {
	var app application

	kingpinApp := kingpin.New("hmp", "hetzner-machine-provider")
	kingpinApp.HelpFlag.Short('h')
	kingpinApp.Version(version)
	kingpinApp.Flag("resource-name-prefix", "cloud resource name prefix").Envar("CUSTOM_ENV_HMP_RESOURCE_NAME_PREFIX").Default("hmp-job-").StringVar(&app.resourceNamePrefix)

	// set resource name prefix before any command is executed
	kingpinApp.PreAction(func(_ *kingpin.ParseContext) error {
		validationError := helper.SetResourceNamePrefix(app.resourceNamePrefix)
		if validationError != nil {
			fmt.Fprintf(os.Stderr, "❌ %s\n", validationError)
		}
		return validationError
	})

	prepareCmd := kingpinApp.Command("prepare", "prepare the environment").PreAction(app.prepareClient).Action(app.prepare)
	prepareCmd.Flag("hcloud-token", "hcloud token").Envar("HCLOUD_TOKEN").Required().StringVar(&app.hcloudToken)
	prepareCmd.Flag("job-id", "job id").Envar("CI_JOB_ID").Envar("CUSTOM_ENV_CI_JOB_ID").Required().StringVar(&app.jobID)
	prepareCmd.Flag("prepare.server-wait-deadline", "deadline for server to become reachable").Envar("CUSTOM_ENV_HMP_SERVER_WAIT_DEADLINE").Default("5m").DurationVar(&app.prepareOptions.WaitDeadline)
	prepareCmd.Flag("prepare.additional-authorized-keys", "specify additional authorized keys separated by '\\n'").Envar("CUSTOM_ENV_HMP_ADDITIONAL_AUTHORIZED_KEYS").StringVar(&app.prepareOptions.AdditionalAuthorizedKeys)
	prepareCmd.Flag("vm.image", "vm image").Envar("CUSTOM_ENV_CI_JOB_IMAGE").Default("ubuntu-22.04").StringVar(&app.vmParams.Image)
	prepareCmd.Flag("vm.type", "vm type").Envar("CUSTOM_ENV_HCLOUD_SERVER_TYPE").Default("ccx12").StringVar(&app.vmParams.Type)
	prepareCmd.Flag("vm.location", "vm location").Envar("CUSTOM_ENV_HCLOUD_SERVER_LOCATION").Default("fsn1").StringVar(&app.vmParams.Location)

	cleanupCmd := kingpinApp.Command("cleanup", "cleanup the environment").PreAction(app.prepareClient).Action(app.cleanup)
	cleanupCmd.Flag("hcloud-token", "hcloud token").Envar("HCLOUD_TOKEN").Required().StringVar(&app.hcloudToken)
	cleanupCmd.Flag("job-id", "job id").Envar("CI_JOB_ID").Envar("CUSTOM_ENV_CI_JOB_ID").Required().StringVar(&app.jobID)

	execCmd := kingpinApp.Command("exec", "execute a command").Action(app.exec)
	execCmd.Arg("scriptPath", "script to execute").Required().StringVar(&app.execScriptPath)
	execCmd.Arg("stageName", "stage name").Required().StringVar(&app.execStageName)

	kingpinApp.Command("configure", "configure the environment").Action(app.configure)

	_, err := kingpinApp.Parse(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
