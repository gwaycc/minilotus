package main

import (
	"context"

	"github.com/gwaycc/minilotus/lib/rpc"
	"github.com/libp2p/go-libp2p-core/peer"
)

type RpcClient struct {
	c rpc.Client
}

func NewRpcClient(host, token string) *RpcClient {
	return &RpcClient{
		c: rpc.NewClient(host, RPC_SERVICE_NAME, token),
	}
}

func (r *RpcClient) CurrentHeight(ctx context.Context) (*CurrentHeightRet, error) {
	arg := &CurrentHeightArg{}
	ret := &CurrentHeightRet{}
	return ret, r.c.Call(ctx, "CurrentHeight", arg, ret)
}

func (r *RpcClient) Peers(ctx context.Context) (*PeersRet, error) {
	arg := &PeersArg{}
	ret := &PeersRet{}
	return ret, r.c.Call(ctx, "Peers", arg, ret)
}

func (r *RpcClient) Connect(ctx context.Context, addr peer.AddrInfo) (*ConnectRet, error) {
	arg := &ConnectArg{Addr: addr}
	ret := &ConnectRet{}
	return ret, r.c.Call(ctx, "Connect", arg, ret)
}
