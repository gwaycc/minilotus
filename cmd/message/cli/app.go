package cli

import (
	"os"

	"github.com/gwaycc/minilotus/cmd"
	"github.com/gwaycc/minilotus/version"
	"github.com/urfave/cli/v2"
)

var app = &cmd.App{
	&cli.App{
		Name:    "Mini Lotus",
		Version: version.Version(),
		Usage:   "Mini Lotus",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "network",
				Value: "testnetnet",
				Usage: "the netkind of filecoin, support: 'testnetnet', 'calibrationnet'; the mainnet is called 'testnetnet'",
			},
			&cli.StringFlag{
				Name:  "repo",
				Value: "./repo",
				Usage: "repository of data",
			},
		},
	},
}

func Run() error {
	return app.Run(os.Args)
}
