package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"

	"github.com/gwaylib/errors"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/urfave/cli/v2"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/wallet"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
)

var app = &cli.App{
	Name:    "Lotus Send",
	Version: "0.0.1",
	Usage:   "Send Filecoin message to the net, lotus-send [to address] [0.1/0.1fil/0.1afil]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "wallet-key",
			Value: "./wallet.key",
			Usage: "the private key of filecoin wallet",
		},
		&cli.StringFlag{
			Name:  "netname",
			Value: "testnetnet",
			Usage: "the netkind of filecoin",
		},
		&cli.BoolFlag{
			Name:  "subscribe-block",
			Usage: "the netkind of filecoin",
		},
		&cli.Uint64Flag{
			Name:  "nonce",
			Value: 0,
			Usage: "message nonce of wallet address, you can get it from the statistics of messages in the public filecoin browsers",
		},
		&cli.StringFlag{
			Name:  "basefee",
			Value: "1afil", // 1 attoFIL
			Usage: "you can get it from the public filecoin browsers",
		},
	},
	Action: func(cctx *cli.Context) error {
		ctx := context.TODO()

		// loading wallet private key
		kiHex, err := ioutil.ReadFile("private.key")
		if err != nil {
			return errors.As(err)
		}
		kiBytes, err := hex.DecodeString(string(kiHex))
		if err != nil {
			return errors.As(err)
		}
		ki := types.KeyInfo{}
		if err := json.Unmarshal(kiBytes, &ki); err != nil {
			return errors.As(err)
		}
		key, err := wallet.NewKey(ctx, ki)
		if err != nil {
			return errors.As(err)
		}

		if cctx.Args().Len() != 2 {
			return errors.New("expects target and amount")
		}
		to, err := address.NewFromString(cctx.Args().Get(0))
		if err != nil {
			return errors.As(err)
		}

		value, err := types.ParseFIL(cctx.Args().Get(1))
		if err != nil {
			return errors.As(err)
		}

		nonce := cctx.Uint64("nonce")
		method := abi.MethodNum(0)

		// read from private key
		from := key.Address

		bFee, err := types.ParseFIL(cctx.String("basefee"))
		if err != nil {
			return errors.As(err)
		}

		// TODO: compute the gas
		baseFee := abi.TokenAmount(bFee)
		gasLimit := int64(1)
		gasFeeCap := abi.NewTokenAmount(0)
		gasPremium := abi.NewTokenAmount(0)
		_ = baseFee

		msg := &types.Message{
			Version: 0,
			To:      to,
			From:    from,
			Nonce:   nonce,
			Method:  method,
			Value:   abi.TokenAmount(value),

			GasLimit:   gasLimit,
			GasFeeCap:  gasFeeCap,
			GasPremium: gasPremium,
		}
		signed, err := Sign(ctx, ki, msg)
		if err != nil {
			return errors.As(err)
		}

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

		netName := dtypes.NetworkName(cctx.String("netname"))

		if err := ConnectBootstrap(ctx, srcHost, string(netName)); err != nil {
			return errors.As(err)
		}

		ps, err := pubsub.NewGossipSub(ctx, srcHost)
		if err != nil {
			return errors.As(err)
		}

		if cctx.Bool("subscribe-block") {
			// send message
			subTitle := build.BlocksTopic(netName)
			subTopic, err := ps.Join(subTitle)
			if err != nil {
				return errors.As(err)
			}
			go DaemonSub(ctx, subTopic)
		}

		// send message
		sendTitle := build.MessagesTopic(netName)
		sendTopic, err := ps.Join(sendTitle)
		if err != nil {
			return errors.As(err)
		}
		if err := Publish(ctx, sendTopic, signed); err != nil {
			return errors.As(err)
		}
		fmt.Printf("message has sent: %+v\n", *signed)

		// end
		fmt.Println("[ctrl+c to exit]")
		end := make(chan os.Signal, 2)
		signal.Notify(end, os.Interrupt, os.Kill)
		<-end

		return nil
	},
}

func main() {
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
