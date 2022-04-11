package cli

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/gwaycc/minilotus/lib/rpc"
	"github.com/gwaycc/minilotus/node/chain"
	"github.com/gwaycc/minilotus/node/repo"
	"github.com/gwaylib/errors"
	"github.com/gwaylib/log"
	"github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/smallnest/rpcx/protocol"
	"github.com/urfave/cli/v2"
)

func init() {
	app.Register("daemon",
		&cli.Command{
			Name: "daemon",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "rpc-listen",
					Value: ":9882",
					Usage: "api listen address.",
				},
				&cli.StringFlag{
					Name:  "rpc-discovery",
					Value: "rpc-server-1",
					Usage: "rpc discovery name",
				},
			},
			Action: func(cctx *cli.Context) error {
				ctx := cctx.Context
				// implement the rpc
				rpcAddr := cctx.String("rpc-listen")
				//rpcHost := cctx.String("rpc-discovery") // TODO: make a service discovery
				r, err := repo.NewRepo(repo.ExpandPath(cctx.String("repo")))
				if err != nil {
					return errors.As(err)
				}
				token, err := r.ReadToken()
				if err != nil {
					return errors.As(err)
				}

				// rpc auth
				auth := func(ctx context.Context, req *protocol.Message, clientToken string) error {
					// TODO: parse token params
					if clientToken != token {
						return rpc.ErrInvalidToken.As(req.ServiceMethod, clientToken)
					}
					return nil
				}
				// listen the api address
				go func() {
					s := rpc.NewServer(auth, chain.RPC_SERVICE_NAME, chain.RpcSrvInstance())
					log.Infof("rpc listen at:%s", rpcAddr)
					if err := s.Serve("reuseport", rpcAddr); err != nil {
						log.Exit(2, errors.As(err))
					}
				}()

				// waiting exit.
				opts := []libp2p.Option{
					chain.NetID(),
					libp2p.NoListenAddrs,
					libp2p.Ping(true),
					libp2p.ConnectionManager(connmgr.NewConnManager(5, 50, 20*time.Second)),
					libp2p.UserAgent("minilotus-1.13.1"),
					libp2p.FallbackDefaults,
				}
				srcHost, err := libp2p.New(ctx, opts...)
				if err != nil {
					return errors.As(err)
				}
				defer srcHost.Close()

				chain.InitRpcSrv(srcHost)

				netName := dtypes.NetworkName(cctx.String("network"))
				ps, err := pubsub.NewGossipSub(ctx, srcHost)
				if err != nil {
					return errors.As(err)
				}
				blkTopic, err := ps.Join(build.BlocksTopic(netName))
				if err != nil {
					return errors.As(err)
				}
				var retryCtx context.Context
				var retryCancel context.CancelFunc
				go func() {
				connect:
					select {
					case <-ctx.Done():
					default:
					}
					if retryCancel != nil {
						retryCancel()
					}
					retryCtx, retryCancel = context.WithCancel(ctx)
					addrs, err := chain.GetConnectTrustNode(retryCtx, string(netName))
					if err != nil {
						log.Warn(errors.As(err))
						time.Sleep(1e9)
						goto connect
					}
					if resp, err := chain.ConnectTrustNode(retryCtx, srcHost, addrs); err != nil {
						log.Warn(errors.As(err))
						time.Sleep(1e9)
						goto connect
					} else {
						log.Debug(resp)
					}
					log.Infof("Join the network: %s", netName)
					go chain.DaemonSubMsg(retryCtx, ps, build.MessagesTopic(netName))
					if err := chain.DaemonSubBlock(retryCtx, blkTopic, 1*time.Minute); err != nil {
						log.Error(errors.As(err))
						time.Sleep(1e9)
						goto connect
					}
				}()

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
