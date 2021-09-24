package main

import (
	"context"

	"github.com/filecoin-project/lotus/chain/types"
	"github.com/gwaylib/errors"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func Publish(ctx context.Context, topic *pubsub.Topic, signed *types.SignedMessage) error {
	msgb, err := signed.Serialize()
	if err != nil {
		return errors.As(err)
	}
	return topic.Publish(ctx, msgb)
}
