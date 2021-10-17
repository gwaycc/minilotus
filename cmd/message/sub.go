package main

import (
	"context"

	"github.com/gwaylib/errors"
	"github.com/gwaylib/log"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/filecoin-project/lotus/chain/types"
)

func DaemonSubBlock(ctx context.Context, topic *pubsub.Topic) error {
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

		log.Info("waitting the blocks")
		m, err := sub.Next(ctx)
		if err != nil {
			log.Warn(errors.As(err))
			continue
		}
		msg, err := types.DecodeSignedMessage(m.Data)
		if err != nil {
			log.Warn(errors.As(err, *m))
			continue
		}
		log.Infof("%+v", msg.Message)
	}
	return nil
}

func DaemonSubMsg(ctx context.Context, topic *pubsub.Topic) error {
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

		log.Info("waitting the messages")
		m, err := sub.Next(ctx)
		if err != nil {
			log.Warn(errors.As(err))
			continue
		}
		msg, err := types.DecodeSignedMessage(m.Data)
		if err != nil {
			log.Warn(errors.As(err, *m))
			continue
		}
		log.Infof("%+v", msg.Message)
	}
	return nil
}
