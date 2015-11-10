package main

import "github.com/codegangsta/cli"

const (
	serviceVersion = "0.0.1"
)

func cliVersion(c *cli.Context) {
	print(serviceVersion)
}
