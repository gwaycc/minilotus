package cli

import (
	"context"
	"fmt"

	"github.com/gwaycc/minilotus/node/chain"
	"github.com/gwaycc/minilotus/node/repo"
	"github.com/urfave/cli/v2"

	"github.com/gwaylib/errors"
)

func init() {
	subCmds := []*cli.Command{
		&cli.Command{
			Name:  "gas-info",
			Usage: "get the current block gas infomation",
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

				c := chain.NewRpcClient(rpcApi, token)
				ret, err := c.CurrentGasInfo(ctx)
				if err != nil {
					return errors.As(err)
				}
				fmt.Printf("%+v\n", *ret)
				return nil
			},
		},

		&cli.Command{
			Name:  "state",
			Usage: "get the current tipset",
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

				c := chain.NewRpcClient(rpcApi, token)
				ret, err := c.CurrentTipset(ctx)
				if err != nil {
					return errors.As(err)
				}
				fmt.Printf("mpool len:%d\n", ret.MpoolLen)
				ret.Tipset.Dump()

				return nil
			},
		},
	}
	app.Register("block",
		&cli.Command{
			Name: "block",
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
