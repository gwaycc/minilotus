package main

import (
	"context"

	"github.com/gwaylib/errors"
	"github.com/gwaylib/log"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
)

const (
	RPC_SERVICE_NAME = "RPCService"
)

type RpcService struct {
	host host.Host

	ts Tipset
}

var RpcSrv = &RpcService{
	ts: Tipset{},
	// TODO: more restfull init
}

type CurrentHeightArg struct{}
type CurrentHeightRet struct {
	Info Tipset
}

func (r *RpcService) CurrentHeight(ctx context.Context, arg *CurrentHeightArg, ret *CurrentHeightRet) error {
	log.Debug("current height called:%+v", r.ts)
	ret.Info = r.ts
	return nil
}

type PeersArg struct{}
type PeersRet struct {
	Peers peer.IDSlice
}

func (r *RpcService) Peers(ctx context.Context, arg *PeersArg, ret *PeersRet) error {
	ret.Peers = r.host.Peerstore().Peers()
	return nil
}

type ConnectArg struct {
	Addr peer.AddrInfo
}
type ConnectRet struct{}

func (r *RpcService) Connect(ctx context.Context, arg *ConnectArg, ret *ConnectRet) error {
	if err := r.host.Connect(ctx, arg.Addr); err != nil {
		return errors.As(err)
	}
	return nil
}
