package chain

import (
	"fmt"
	"strings"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/gwaylib/errors"
)

type BlockMsg struct {
	*types.BlockMsg
}

func (b *BlockMsg) GetSamePart() string {
	// all blocks in same tipset
	return fmt.Sprintf("Parents:%+v,ParentsWeight:%+v,Height:%+v,ParentStateRoot:%+v,PerentMessagereceipts:%+v", b.Header.Parents, b.Header.ParentWeight, b.Header.Height, b.Header.ParentStateRoot, b.Header.ParentMessageReceipts)
}
func (b *BlockMsg) Compare(to *BlockMsg) int {
	bStr := b.GetSamePart()
	toStr := to.GetSamePart()
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
	return b.Headers() + "|" + b.GetSamePart()
}

type Tipset map[string]*BlockMsg

func (t Tipset) Dump() {
	for _, val := range t {
		fmt.Println(val.String())
	}
}
func (t Tipset) Height() (int64, error) {
	for _, b := range t {
		return int64(b.Header.Height), nil
	}
	return 0, errors.ErrNoData
}
func (t Tipset) ParentBaseFee() (abi.TokenAmount, error) {
	for _, b := range t {
		return b.Header.ParentBaseFee, nil
	}
	return abi.NewTokenAmount(0), errors.ErrNoData
}
func (t Tipset) DumpString() []string {
	result := []string{}
	for _, val := range t {
		result = append(result, val.String())
	}
	return result
}

func (t Tipset) Put(b *BlockMsg) ([]*BlockMsg, error) {
	key := fmt.Sprintf("%d%x", b.Header.BlockSig.Type, b.Header.BlockSig.Data)
	sameNum := 0
	diffNum := 0
	removed := []*BlockMsg{}
	for key, val := range t {
		if val.Header.Height < b.Header.Height {
			// TODO: make sure the sign is verified
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
