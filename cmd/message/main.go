package main

import (
	"github.com/gwaycc/minilotus/cmd/message/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		panic(err)
	}
}
