package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "nomadpanel"
	app.Usage = "nomadpanel"
	app.Version = serviceVersion
	app.Author = "antoniofernandezvara@gmail.com"
	app.Email = "antoniofernandezvara@gmail.com"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "log-level",
			Value:  "info",
			Usage:  "Log level (options: debug, info, warn, error, fatal, panic)",
			EnvVar: "LOG_LEVEL",
		},
		cli.StringFlag{
			Name:   "nomad-address",
			Usage:  "Nomad Address",
			EnvVar: "NOMAD_ADDRESS",
			Value:  "http://127.0.0.1:4646",
		},
		cli.StringFlag{
			Name:   "nomad-region",
			Usage:  "Nomad Region",
			EnvVar: "NOMAD_REGION",
			Value:  "global",
		},
		cli.IntFlag{
			Name:   "wait-time",
			Usage:  "WaitTime",
			EnvVar: "WAIT_TIME",
			Value:  3600,
		},
		cli.StringFlag{
			Name:   "api-addr",
			Usage:  "Listener address",
			EnvVar: "API_ADDR",
			Value:  "0.0.0.0:8000",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "start",
			Usage:  "service start",
			Action: cliStart,
		},
		{
			Name:   "version",
			Usage:  "service version",
			Action: cliVersion,
		},
	}

	app.Run(os.Args)
}
