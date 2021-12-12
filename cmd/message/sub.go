package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gwaylib/errors"
	"github.com/gwaylib/log"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/filecoin-project/lotus/chain/types"
)

type BlockMsg struct {
	*types.BlockMsg
}

func (b *BlockMsg) SamePart() string {
	// all blocks in same tipset
	return fmt.Sprintf("Parents:%+v,ParentsWeight:%+v,Height:%+v,ParentStateRoot:%+v,PerentMessagereceipts:%+v", b.Header.Parents, b.Header.ParentWeight, b.Header.Height, b.Header.ParentStateRoot, b.Header.ParentMessageReceipts)
}
func (b *BlockMsg) Compare(to *BlockMsg) int {
	bStr := b.SamePart()
	toStr := to.SamePart()
	return strings.Compare(bStr, toStr)
}

func (b *BlockMsg) Headers() string {
	return fmt.Sprintf(
		"IsValided:%t,Miner:%+v,Height:%d,Timestamp:%d,ParentBaseFee:%s",
		b.Header.IsValidated(),
		b.Header.Miner,
		b.Header.Height,
		b.Header.Timestamp,
		b.Header.ParentBaseFee.String(),
	)
}
func (b *BlockMsg) String() string {
	return b.Headers() + "|" + b.SamePart()
}

type tipset map[string]*BlockMsg

func (t tipset) Dump() {
	for _, val := range t {
		log.Debug(val.String())
	}
}
func (t tipset) Put(b *BlockMsg) ([]*BlockMsg, error) {
	key := fmt.Sprintf("%d%x", b.Header.BlockSig.Type, b.Header.BlockSig.Data)
	sameNum := 0
	diffNum := 0
	removed := []*BlockMsg{}
	for key, val := range t {
		if val.Header.Height < b.Header.Height {
			// TODO: need make sure the sign is verified
			t[key] = b

			delete(t, key)
			removed = append(removed, val)
			continue
		}
		switch val.Compare(b) {
		case 0:
			sameNum++
		default:
			diffNum++
		}
	}
	if sameNum >= diffNum {
		t[key] = b

		// clean the different
		for key, val := range t {
			if val.Compare(b) != 0 {
				delete(t, key)
				removed = append(removed, val)
			}
		}
		return removed, nil
	}
	return nil, errors.New("fork").As(b.String())
}

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
		ts := tipset{}
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
