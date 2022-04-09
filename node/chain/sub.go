package chain

import (
	"context"
	"time"

	"github.com/gwaylib/errors"
	"github.com/gwaylib/log"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/filecoin-project/lotus/chain/types"
)

func DaemonSubBlock(ctx context.Context, topic *pubsub.Topic, timeout time.Duration) error {
	defer topic.Close()

	sub, err := topic.Subscribe()
	if err != nil {
		return err
	}
	timeoutTimer := time.NewTimer(timeout)
	timeoutCtx, timeoutFn := context.WithCancel(ctx)
	defer timeoutFn()

	alive := make(chan error, 1)
	go func() {
		ts := rpcSrv.ts
		for {
			log.Info("waitting the blocks")
			m, err := sub.Next(timeoutCtx)
			if err != nil {
				alive <- errors.As(err)
				return
			}
			alive <- nil
			blocks, err := types.DecodeBlockMsg(m.Data)
			if err != nil {
				log.Warn(errors.As(err, *m))
				continue
			}
			// TODO: verify the blocksig
			b := &BlockMsg{blocks}
			removed, err := ts.Put(b)
			if err != nil {
				log.Warn(errors.As(err))
			} else {
				log.Debugf("new block:%s", b.String())
				for _, r := range removed {
					log.Debugf("remove block:%s", r.Headers())
				}
			}
		}
	}()
	for {
		select {
		case <-timeoutTimer.C:
			timeoutFn()
			return errors.New("data timeout")
		case <-ctx.Done():
			return ctx.Err()
		case <-alive:
			timeoutTimer.Reset(timeout)
		}
	}
	return nil
}

var countMsg = 0

func DaemonSubMsg(ctx context.Context, ps *pubsub.PubSub, tc string) error {
	topic, err := ps.Join(tc)
	if err != nil {
		return errors.As(err)
	}
	defer topic.Close()

	sub, err := topic.Subscribe()
	if err != nil {
		return err
	}
	for {
		//log.Info("waitting the messages")
		m, err := sub.Next(ctx)
		if err != nil {
			return errors.As(err)
		}
		msg, err := types.DecodeSignedMessage(m.Data)
		if err != nil {
			log.Warn(errors.As(err, *m))
			continue
		}
		//log.Infof("%+v", msg.Message)
		countMsg++
		if countMsg%100 == 0 {
			log.Infof("msg received:%d, current:%+v", countMsg, msg)
		}
	}
	return nil
}
