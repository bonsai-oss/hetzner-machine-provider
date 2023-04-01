package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/fatih/color"
	"github.com/hetznercloud/hcloud-go/hcloud"

	"hcloud-machine-provider/internal/actions"
)

var version = "dev"

type application struct {
	hcloudToken string
	jobID       string

	execScriptPath string
	execStageName  string

	hcloudClient *hcloud.Client
}

func (a *application) prepare(_ *kingpin.ParseContext) error {
	color.Green("🚀 Preparing environment")
	return actions.Prepare(a.hcloudClient, a.jobID)
}

func (a *application) cleanup(_ *kingpin.ParseContext) error {
	color.Green("🧼 Cleaning up resources")
	return actions.Cleanup(a.hcloudClient, a.jobID)
}

func (a *application) exec(_ *kingpin.ParseContext) error {
	return actions.Exec(a.execScriptPath, a.execStageName)
}

func (a *application) prepareClient(_ *kingpin.ParseContext) error {
	a.hcloudClient = hcloud.NewClient(hcloud.WithToken(a.hcloudToken), hcloud.WithApplication("hmp", version))
	return nil
}

func main() {
	var app application

	kingpinApp := kingpin.New("hmp", "hcloud-machine-provider")
	kingpinApp.HelpFlag.Short('h')
	kingpinApp.Version(version)

	prepareCmd := kingpinApp.Command("prepare", "prepare the environment").PreAction(app.prepareClient).Action(app.prepare)
	prepareCmd.Flag("hcloud-token", "hcloud token").Envar("HCLOUD_TOKEN").Required().StringVar(&app.hcloudToken)
	prepareCmd.Flag("job-id", "job id").Envar("CI_JOB_ID").Envar("CUSTOM_ENV_CI_JOB_ID").Required().StringVar(&app.jobID)

	cleanupCmd := kingpinApp.Command("cleanup", "cleanup the environment").PreAction(app.prepareClient).Action(app.cleanup)
	cleanupCmd.Flag("hcloud-token", "hcloud token").Envar("HCLOUD_TOKEN").Required().StringVar(&app.hcloudToken)
	cleanupCmd.Flag("job-id", "job id").Envar("CI_JOB_ID").Envar("CUSTOM_ENV_CI_JOB_ID").Required().StringVar(&app.jobID)

	execCmd := kingpinApp.Command("exec", "execute a command").Action(app.exec)
	execCmd.Arg("scriptPath", "script to execute").Required().StringVar(&app.execScriptPath)
	execCmd.Arg("stageName", "stage name").Required().StringVar(&app.execStageName)

	_, err := kingpinApp.Parse(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
