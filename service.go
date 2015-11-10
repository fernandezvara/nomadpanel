package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/fernandezvara/nomadpanel/api"
	"github.com/fernandezvara/nomadpanel/usage"
	nomadapi "github.com/hashicorp/nomad/api"
)

func cliStart(c *cli.Context) {
	config := nomadapi.DefaultConfig()
	config.Address = c.GlobalString("nomad-address")
	config.Region = c.GlobalString("nomad-region")
	config.WaitTime = time.Duration(c.GlobalInt("wait-time")) * time.Second
	client, err := nomadapi.NewClient(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log := logrus.New()
	level, err := logrus.ParseLevel(c.GlobalString("log-level"))
	if err != nil {
		fmt.Println("incorrect log-level")
		os.Exit(2)
	}
	log.Out = os.Stderr
	log.Level = level

	usage := usage.NewUsage(client, time.Duration(c.GlobalInt("wait-time"))*time.Second, log)
	usage.Loop()

	context := api.NewContext(c.GlobalString("api-addr"), serviceVersion, usage, log)
	log.Fatal(api.ListenAndServe(context))

}
