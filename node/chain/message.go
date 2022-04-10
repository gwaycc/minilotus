package chain

import (
	"fmt"
	"sync"

	"github.com/filecoin-project/lotus/chain/types"
	"github.com/gwaylib/errors"
)

// TODO: move to redis or etcd
type Mpool struct {
	poolLk sync.Mutex

	cidPool  map[string]*types.SignedMessage // TODO: gc?
	addrPool map[string]*types.SignedMessage
}

func (m *Mpool) Len() int {
	m.poolLk.Lock()
	defer m.poolLk.Unlock()
	return len(m.cidPool)
}

func (m *Mpool) GetMessageByCid(cid string) (*types.SignedMessage, error) {
	m.poolLk.Lock()
	defer m.poolLk.Unlock()
	smsg, ok := m.cidPool[cid]
	if !ok {
		return nil, errors.ErrNoData.As(cid)
	}
	return smsg, nil
}
func (m *Mpool) GetMessageByNonce(addr string, nonce int64) (*types.SignedMessage, error) {
	m.poolLk.Lock()
	defer m.poolLk.Unlock()
	key := fmt.Sprintf("%s_%d", addr, nonce)
	smsg, ok := m.addrPool[key]
	if !ok {
		return nil, errors.ErrNoData.As(addr, nonce)
	}
	return smsg, nil
}

func (m *Mpool) PutMessage(smsg *types.SignedMessage) error {
	m.poolLk.Lock()
	defer m.poolLk.Unlock()

	cid := smsg.Cid().String()
	addrKey := fmt.Sprintf("%s_%d", smsg.Message.From.String(), smsg.Message.Nonce)
	m.cidPool[cid] = smsg
	m.addrPool[addrKey] = smsg
	return nil
}

func (m *Mpool) DelMessageByCid(cid string) (*types.SignedMessage, error) {
	m.poolLk.Lock()
	defer m.poolLk.Unlock()

	smsg, ok := m.cidPool[cid]
	if !ok {
		return nil, errors.ErrNoData.As(cid)
	}
	addrKey := fmt.Sprintf("%s_%d", smsg.Message.From.String(), smsg.Message.Nonce)
	delete(m.cidPool, cid)
	delete(m.addrPool, addrKey)
	return smsg, nil
}
func (m *Mpool) DelMessageByNonce(addr string, nonce int64) (*types.SignedMessage, error) {
	m.poolLk.Lock()
	defer m.poolLk.Unlock()
	addrKey := fmt.Sprintf("%s_%d", addr, nonce)
	smsg, ok := m.addrPool[addrKey]
	if !ok {
		return nil, errors.ErrNoData.As(addr, nonce)
	}
	cid := smsg.Cid().String()
	delete(m.cidPool, cid)
	delete(m.addrPool, addrKey)
	return smsg, nil
}
