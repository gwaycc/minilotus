package main

import (
	"os"

	"github.com/gwaycc/minilotus/cmd"
	"github.com/urfave/cli/v2"
)

var app = &cmd.App{
	&cli.App{
		Name:    "Mini Lotus",
		Version: cmd.Version(),
		Usage:   "Mini Lotus",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "network",
				Value: "testnetnet",
				Usage: "the netkind of filecoin, support: 'testnetnet', 'calibrationnet'; the mainnet is called 'testnetnet'",
			},
		},
	},
}

func main() {
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
