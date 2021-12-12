module github.com/gwaycc/minilotus

go 1.16

require (
	github.com/filecoin-project/go-address v0.0.5
	github.com/filecoin-project/go-jsonrpc v0.1.4-0.20210217175800-45ea43ac2bec
	github.com/filecoin-project/go-state-types v0.1.1-0.20210810190654-139e0e79e69e
	github.com/filecoin-project/lotus v1.11.2
	github.com/google/uuid v1.2.0
	github.com/gwaylib/errors v0.0.0-20190905023356-162e59439c92
	github.com/gwaylib/log v0.0.0-20210507100943-24bc495476d8
	github.com/kr/text v0.2.0 // indirect
	github.com/libp2p/go-libp2p v0.14.2
	github.com/libp2p/go-libp2p-core v0.8.6
	github.com/libp2p/go-libp2p-pubsub v0.5.4
	github.com/mitchellh/go-homedir v1.1.0
	github.com/multiformats/go-multiaddr v0.3.3
	github.com/urfave/cli/v2 v2.2.0
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
)

replace github.com/libp2p/go-libp2p-yamux => github.com/libp2p/go-libp2p-yamux v0.5.1

replace github.com/filecoin-project/go-jsonrpc => github.com/lib10/go-jsonrpc v0.0.0-20210806021800-80a4ef41e98d

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi

replace google.golang.org/grpc => google.golang.org/grpc v1.29.1
