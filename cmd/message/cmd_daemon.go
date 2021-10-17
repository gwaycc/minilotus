package main

import (
	"os"
	"os/signal"

	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/gwaylib/errors"
	"github.com/gwaylib/log"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/urfave/cli/v2"
)

func init() {
	app.Register("daemon",
		&cli.Command{
			Name:  "daemon",
			Flags: []cli.Flag{},
			Action: func(cctx *cli.Context) error {
				ctx := cctx.Context

				// waiting exit.
				opts := []libp2p.Option{
					NetID(),
					libp2p.NoListenAddrs,
					libp2p.UserAgent("lotus-1.11.2"),
				}
				srcHost, err := libp2p.New(ctx, opts...)
				if err != nil {
					return errors.As(err)
				}
				defer srcHost.Close()

				netName := dtypes.NetworkName(cctx.String("network"))

				if err := ConnectBootstrap(ctx, srcHost, string(netName)); err != nil {
					return errors.As(err)
				}

				ps, err := pubsub.NewGossipSub(ctx, srcHost)
				if err != nil {
					return errors.As(err)
				}

				blockTopic, err := ps.Join(build.BlocksTopic(netName))
				if err != nil {
					return errors.As(err)
				}
				msgTopic, err := ps.Join(build.MessagesTopic(netName))
				if err != nil {
					return errors.As(err)
				}

				log.Infof("Already join the network: %s", netName)
				go DaemonSubBlock(ctx, blockTopic)
				go DaemonSubMsg(ctx, msgTopic)

				// waiting exit.
				log.Info("[ctrl+c to exit]")
				end := make(chan os.Signal, 2)
				signal.Notify(end, os.Interrupt, os.Kill)
				<-end

				return nil
			},
		},
	)
}
