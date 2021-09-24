package main

import (
	"context"
	"fmt"

	"github.com/gwaylib/errors"
	"github.com/gwaylib/log"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/filecoin-project/lotus/chain/types"
)

func DaemonSub(ctx context.Context, topic *pubsub.Topic) error {
	sub, err := topic.Subscribe()
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			//continue
		}

		fmt.Println("next message")
		m, err := sub.Next(ctx)
		if err != nil {
			log.Warn(errors.As(err))
			continue
		}
		fmt.Println("next message done")
		msg, err := types.DecodeSignedMessage(m.Data)
		if err != nil {
			log.Warn(errors.As(err, *m))
			continue
		}
		log.Debug(msg.Message)
	}
	return nil
}
