package rpc

import (
	"context"
	"sync"
	"time"

	"github.com/gwaylib/errors"
	xclient "github.com/smallnest/rpcx/client"
)

type Client interface {
	Close() error
	IsClosed() bool

	Call(ctx context.Context, method string, args interface{}, reply interface{}) error
}

type client struct {
	dis xclient.ServiceDiscovery

	srvName string
	token   string

	lk sync.Mutex
	xc xclient.XClient
}

func NewClient(addr, srvName, token string) Client {
	d, _ := xclient.NewPeer2PeerDiscovery("tcp@"+addr, "")
	return &client{
		dis: d,

		srvName: srvName,
		token:   token,
	}
}

func (c *client) close() error {
	c.lk.Lock()
	defer c.lk.Unlock()
	if c.xc != nil {
		if err := c.xc.Close(); err != nil {
			return errors.As(err)
		}
		c.xc = nil
	}
	return nil
}

func (c *client) client() xclient.XClient {
	c.lk.Lock()
	defer c.lk.Unlock()
	if c.xc != nil {
		return c.xc
	}

	option := xclient.DefaultOption
	option.IdleTimeout = 10 * time.Second

	xclient := xclient.NewXClient(c.srvName, xclient.Failtry, xclient.RandomSelect, c.dis, option)
	xclient.Auth(c.token)
	c.xc = xclient

	return c.xc
}

func (c *client) Close() error {
	return c.close()
}

func (c *client) IsClosed() bool {
	c.lk.Lock()
	defer c.lk.Unlock()
	return c.xc == nil
}

func (c *client) Call(ctx context.Context, method string, args interface{}, reply interface{}) error {
	xclient := c.client()
	if err := xclient.Call(ctx, method, args, reply); err != nil {
		switch {
		case ErrInvalidToken.Equal(err), ErrEOF.Equal(err):
			c.close()
		}
		return errors.As(err)
	}
	return nil
}
