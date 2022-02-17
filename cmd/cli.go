package cmd

import (
	"github.com/urfave/cli/v2"
)

func NewCli() *cli.App {
	app := &cli.App{
		Commands: []*cli.Command{},
	}

	return app
}
