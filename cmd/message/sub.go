package main

import (
	"context"
	"fmt"
	"strings"

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
		log.Info(val.String())
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

func DaemonSubBlock(ctx context.Context, topic *pubsub.Topic) error {
	sub, err := topic.Subscribe()
	if err != nil {
		return err
	}
	ts := tipset{}
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
			log.Infof("new block:%s", b.String())
			for _, r := range removed {
				log.Infof("remove block:%s", r.Headers())
			}
		}
	}
	return nil
}

var countMsg = 0

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

		//log.Info("waitting the messages")
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
		//log.Infof("%+v", msg.Message)
		countMsg++
		if countMsg%100 == 0 {
			log.Infof("msg received:%d, current:%+v", countMsg, msg)
		}
	}
	return nil
}
