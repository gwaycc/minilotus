package cli

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gwaycc/minilotus/node/chain"
	"github.com/gwaycc/minilotus/node/repo"
	"github.com/urfave/cli/v2"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/wallet"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/gwaylib/errors"
)

func init() {
	subCmds := []*cli.Command{
		&cli.Command{
			Name:  "send",
			Usage: "Send a Filecoin message to the net, send [json message]",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "wallet-key",
					Value: "./private.key",
					Usage: "the private key of filecoin wallet",
				},
				&cli.BoolFlag{
					Name:  "wallet-encrypted",
					Value: true,
					Usage: "the private key is encrypted",
				},
				&cli.Uint64Flag{
					Name:  "nonce",
					Value: 0,
					Usage: "TODO",
				},
				&cli.StringFlag{
					Name:  "basefee",
					Value: "", // 1 attoFIL
					Usage: "1afil for unit, unset will do auto fill",
				},
				&cli.StringFlag{
					Name:  "data",
					Value: "", // json params string format
					Usage: "TODO",
				},
			},
			Action: func(cctx *cli.Context) error {
				ctx := context.TODO()

				walletKey := cctx.String("wallet-key")
				// loading wallet private key
				kiHex, err := ioutil.ReadFile(walletKey)
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

				// TODO: decrypt private key
				key, err := wallet.NewKey(ki)
				if err != nil {
					return errors.As(err)
				}

				if cctx.Args().Len() != 1 {
					return errors.New("expects json string input")
				}
				msgStr := cctx.Args().Get(0)
				msg := &types.Message{}
				if err := json.Unmarshal([]byte(msgStr), msg); err != nil {
					return errors.As(err)
				}

				nonce := cctx.Uint64("nonce")

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

				// TODO: fix gas
				msg.From = key.Address
				msg.Nonce = nonce
				msg.GasLimit = gasLimit
				msg.GasFeeCap = gasFeeCap
				msg.GasPremium = gasPremium
				signed, err := chain.Sign(ctx, ki, msg)
				if err != nil {
					return errors.As(err)
				}

				netName := dtypes.NetworkName(cctx.String("network"))

				// send message
				topic := build.MessagesTopic(netName)

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

				if _, err := c.Publish(ctx, topic, signed); err != nil {
					return errors.As(err)
				}
				fmt.Printf("message has sent: %+v\n", *signed)
				return nil
			},
		},
	}
	app.Register("message",
		&cli.Command{
			Name:        "msg",
			Subcommands: subCmds,
		},
	)

}
