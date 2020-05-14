package main

import (
	"gitlab.connectwisedev.com/platform/platform-tasking-service/utils/migrate-targets-and-target-type/commands"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Continuum Juno Automation TaskTargets converting"
	app.Usage = "This piece of software will convert task targets set to task targets map"
	app.Version = "0.0.1"

	// global flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Usage:  "Path to the config file",
			Value:  "./config.json",
			EnvVar: "K_CONFIG",
		},
	}

	app.Action = commands.UpdateTasksTable
	app.Run(os.Args)
}
