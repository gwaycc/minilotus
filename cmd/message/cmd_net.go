package main

import (
	"context"
	"fmt"

	"github.com/gwaycc/minilotus/node/repo"
	"github.com/urfave/cli/v2"

	"github.com/gwaylib/errors"
)

func init() {
	subCmds := []*cli.Command{
		&cli.Command{
			Name:  "peers",
			Usage: "get the peers",
			Flags: []cli.Flag{},
			Action: func(cctx *cli.Context) error {
				ctx := context.TODO()
				rpcApi := cctx.String("rpc-api")
				r, err := repo.NewRepo(repo.ExpandPath(cctx.String("repo")))
				if err != nil {
					return errors.As(err)
				}
				token, err := r.ReadToken()
				if err != nil {
					return errors.As(err)
				}

				c := NewRpcClient(rpcApi, token)

				ret, err := c.Peers(ctx)
				if err != nil {
					return errors.As(err)
				}
				fmt.Printf("len:%d\n", len(ret.Peers))
				for _, p := range ret.Peers {
					fmt.Println(p)
				}
				return nil
			},
		},
	}
	app.Register("net",
		&cli.Command{
			Name: "net",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "rpc-api",
					Value: "127.0.0.1:9882",
					Usage: "the rpc server api",
				},
			},
			Subcommands: subCmds,
		},
	)
}
