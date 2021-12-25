package main

import (
	"context"

	"github.com/gwaycc/minilotus/lib/rpc"
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
