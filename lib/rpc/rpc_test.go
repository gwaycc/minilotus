package rpc

import (
	"context"
	"testing"
	"time"

	"github.com/smallnest/rpcx/protocol"
)

type RpcTestArg struct {
	A int
	B int
}
type RpcTestRet struct {
	C int
}
type RpcTest struct {
}

func (r *RpcTest) Add(ctx context.Context, arg *RpcTestArg, ret *RpcTestRet) error {
	ret.C = arg.A + arg.B
	return nil
}

func TestRpcServer(t *testing.T) {
	pwd := "testing"
	var auth = func(ctx context.Context, req *protocol.Message, token string) error {
		if pwd != token {
			return ErrInvalidToken
		}
		return nil
	}

	impl := &RpcTest{}
	testAddr := "127.0.0.1:3080"
	srvName := "Math"
	s := NewServer(auth, srvName, impl)
	go func() {
		if err := s.Serve("reuseport", testAddr); err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(1e8)
	c := NewClient(testAddr, srvName, pwd)
	arg := &RpcTestArg{
		A: 1,
		B: 2,
	}
	ret := &RpcTestRet{}
	if err := c.Call(context.TODO(), "Add", arg, ret); err != nil {
		t.Fatal(err)
	}
	if ret.C != 3 {
		t.Fatalf("expect 3, but %d\n", ret.C)
	}
}
