package rpc

import (
	"os"

	"github.com/gwaylib/log/logger"
	"github.com/gwaylib/log/logger/adapter/stdio"
	"github.com/gwaylib/log/logger/proto"
	"github.com/smallnest/rpcx/log"
)

func init() {
	//log.SetDummyLogger()
	log.SetLogger(logger.New(&logger.DefaultContext, "rpc", proto.LevelWarn, stdio.New(os.Stdout, os.Stderr)))
}
