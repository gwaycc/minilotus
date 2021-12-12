package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/gwaylib/errors"
	"github.com/gwaylib/log"
	"github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
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
					libp2p.Ping(true),
					libp2p.ConnectionManager(connmgr.NewConnManager(50, 200, 20*time.Second)),
					libp2p.UserAgent("minilotus-1.13.1"),
					libp2p.FallbackDefaults,
				}
				srcHost, err := libp2p.New(ctx, opts...)
				if err != nil {
					return errors.As(err)
				}
				defer srcHost.Close()

				netName := dtypes.NetworkName(cctx.String("network"))
				ps, err := pubsub.NewGossipSub(ctx, srcHost)
				if err != nil {
					return errors.As(err)
				}
				blkTopic, err := ps.Join(build.BlocksTopic(netName))
				if err != nil {
					return errors.As(err)
				}
				go func() {
				connect:
					select {
					case <-ctx.Done():
					default:
					}
					if err := ConnectBootstrap(ctx, srcHost, string(netName)); err != nil {
						log.Warn(errors.As(err))
						time.Sleep(1e9)
						goto connect
					}
					log.Infof("Join the network: %s", netName)
					if err := DaemonSubBlock(ctx, blkTopic, 1*time.Minute); err != nil {
						log.Error(errors.As(err))
						time.Sleep(1e9)
						goto connect
					}
				}()
				// go DaemonSubMsg(ctx, ps, build.MessagesTopic(netName))

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
