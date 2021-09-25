package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gwaycc/lotus-easy/cmd"
	"github.com/urfave/cli/v2"
)

var app = &cmd.App{
	&cli.App{
		Name:    "Easy Lotus",
		Version: cmd.Version(),
		Usage:   "Easy Lotus",
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

	// end
	fmt.Println("[ctrl+c to exit]")
	end := make(chan os.Signal, 2)
	signal.Notify(end, os.Interrupt, os.Kill)
	<-end

}
