package main

import (
	"context"

	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/wallet"
	"github.com/filecoin-project/lotus/lib/sigs"
	"github.com/gwaylib/errors"
)

func Sign(ctx context.Context, ki types.KeyInfo, msg *types.Message) (*types.SignedMessage, error) {
	mb, err := msg.ToStorageBlock()
	if err != nil {
		return nil, errors.As(err)
	}

	sig, err := sigs.Sign(wallet.ActSigType(ki.Type), ki.PrivateKey, mb.Cid().Bytes())
	if err != nil {
		return nil, errors.As(err)
	}

	return &types.SignedMessage{
		Message:   *msg,
		Signature: *sig,
	}, nil
}
func Verify() error {
	return nil
}
