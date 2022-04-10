package chain

import (
	"context"

	"github.com/filecoin-project/lotus/chain/types"
	"github.com/gwaycc/minilotus/lib/rpc"
	"github.com/libp2p/go-libp2p-core/peer"
)

// interface contract
type RpcClient interface {
	Peers(context.Context) (*PeersRet, error)
	Connect(context.Context, peer.AddrInfo) (*ConnectRet, error)
	Publish(context.Context, string, *types.SignedMessage) (*PublishRet, error)

	CurrentTipset(context.Context) (*CurrentTipsetRet, error)
	CurrentHeight(context.Context) (*CurrentHeightRet, error)
	// TODO: fix parent base fee to current fee
	CurrentGasInfo(context.Context) (*CurrentGasInfoRet, error)
}

type rpcClient struct {
	c rpc.Client
}

func NewRpcClient(host, token string) RpcClient {
	return &rpcClient{
		c: rpc.NewClient(host, RPC_SERVICE_NAME, token),
	}
}

func (r *rpcClient) Peers(ctx context.Context) (*PeersRet, error) {
	arg := &PeersArg{}
	ret := &PeersRet{}
	return ret, r.c.Call(ctx, "Peers", arg, ret)
}

func (r *rpcClient) Connect(ctx context.Context, addr peer.AddrInfo) (*ConnectRet, error) {
	arg := &ConnectArg{Addr: addr}
	ret := &ConnectRet{}
	return ret, r.c.Call(ctx, "Connect", arg, ret)
}

func (r *rpcClient) Publish(ctx context.Context, topic string, msg *types.SignedMessage) (*PublishRet, error) {
	arg := &PublishArg{
		Topic:     topic,
		SignedMsg: msg,
	}
	ret := &PublishRet{}
	return ret, r.c.Call(ctx, "Publish", arg, ret)
}

func (r *rpcClient) CurrentTipset(ctx context.Context) (*CurrentTipsetRet, error) {
	arg := &CurrentTipsetArg{}
	ret := &CurrentTipsetRet{}
	return ret, r.c.Call(ctx, "CurrentTipset", arg, ret)
}
func (r *rpcClient) CurrentHeight(ctx context.Context) (*CurrentHeightRet, error) {
	arg := &CurrentHeightArg{}
	ret := &CurrentHeightRet{}
	return ret, r.c.Call(ctx, "CurrentHeight", arg, ret)
}
func (r *rpcClient) CurrentGasInfo(ctx context.Context) (*CurrentGasInfoRet, error) {
	arg := &CurrentGasInfoArg{}
	ret := &CurrentGasInfoRet{}
	return ret, r.c.Call(ctx, "CurrentGasInfo", arg, ret)
}
