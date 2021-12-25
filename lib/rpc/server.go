package rpc

import (
	"context"

	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/server"
)

type Server struct {
	*server.Server
}

type AuthFn func(ctx context.Context, req *protocol.Message, token string) error

// Example:
// s := NewServer(fn, "StructName", StructImpl)
// s.Serve("reuseport",addr)
func NewServer(auth AuthFn, srvName string, srvImpl interface{}) *Server {
	s := server.NewServer()
	s.RegisterName(srvName, srvImpl, "")
	s.AuthFunc = auth
	return &Server{Server: s}
}
