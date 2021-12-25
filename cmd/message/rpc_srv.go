package main

import (
	"context"
)

const (
	RPC_SERVICE_NAME = "RPCService"
)

type RpcService struct {
	ts Tipset
}

var RpcSrv = &RpcService{
	ts: Tipset{},
}

type CurrentHeightArg struct{}
type CurrentHeightRet struct {
	Info []string
}

func (r *RpcService) CurrentHeight(ctx context.Context, arg *CurrentHeightArg, ret *CurrentHeightRet) error {
	ret.Info = r.ts.DumpString()
	return nil
}
