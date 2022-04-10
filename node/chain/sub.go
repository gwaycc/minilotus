package chain

import (
	"context"
	"math"
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
	timeoutCtx, timeoutEnd := context.WithCancel(ctx)
	defer timeoutEnd()

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
			block, err := types.DecodeBlockMsg(m.Data)
			if err != nil {
				log.Warn(errors.As(err, *m))
				continue
			}
			// TODO: verify the blocksig

			blsMessageData := map[string]*types.SignedMessage{}
			for _, cid := range block.BlsMessages {
				cidStr := cid.String()
				blob, _ := rpcSrv.mpool.DelMessageByCid(cidStr)
				if blob != nil {
					blsMessageData[cidStr] = blob
				}
			}
			secpkMessageData := map[string]*types.SignedMessage{}
			for _, cid := range block.SecpkMessages {
				cidStr := cid.String()
				blob, _ := rpcSrv.mpool.DelMessageByCid(cidStr)
				if blob != nil {
					secpkMessageData[cidStr] = blob
				}
			}
			b := &BlockMsg{
				BlockMsg:         block,
				BlsMessageData:   blsMessageData,
				SecpkMessageData: secpkMessageData,
			}
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
			timeoutEnd()
			return errors.New("data timeout")
		case <-ctx.Done():
			return ctx.Err()
		case <-alive:
			timeoutTimer.Reset(timeout)
		}
	}
	return nil
}

var countMsg = int64(0)

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
	exit := make(chan error, 1)
	go func() {
		for {
			//log.Info("waitting the messages")
			m, err := sub.Next(ctx)
			if err != nil {
				exit <- errors.As(err)
				return
			}
			msg, err := types.DecodeSignedMessage(m.Data)
			if err != nil {
				log.Warn(errors.As(err, *m))
				continue
			}
			rpcSrv.mpool.PutMessage(msg)

			countMsg = (countMsg + 1) % math.MaxInt64
			if countMsg%100 == 0 {
				log.Infof("msg received:%d, current size:%d", countMsg, rpcSrv.mpool.Len())
			}
			select {
			case <-ctx.Done():
				return
			default:
				// break
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-exit:
			return errors.As(err)
		}
	}
	return nil
}
