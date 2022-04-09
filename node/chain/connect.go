package chain

import (
	"context"
	"fmt"
	"strings"

	"github.com/filecoin-project/lotus/lib/addrutil"
	"github.com/gwaylib/errors"
	"github.com/gwaylib/log"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
)

const (
	TOPIC_MAINNET  = "testnetnet"
	TOPIC_CALIBNET = "calibrationnet"
)

// see : https://github.com/filecoin-project/lotus/tree/master/build/bootstrap
const (
	mainTrustNode = `
/dns4/bootstrap-0.mainnet.filops.net/tcp/1347/p2p/12D3KooWCVe8MmsEMes2FzgTpt9fXtmCY7wrq91GRiaC8PHSCCBj
/dns4/bootstrap-1.mainnet.filops.net/tcp/1347/p2p/12D3KooWCwevHg1yLCvktf2nvLu7L9894mcrJR4MsBCcm4syShVc
/dns4/bootstrap-2.mainnet.filops.net/tcp/1347/p2p/12D3KooWEWVwHGn2yR36gKLozmb4YjDJGerotAPGxmdWZx2nxMC4
/dns4/bootstrap-3.mainnet.filops.net/tcp/1347/p2p/12D3KooWKhgq8c7NQ9iGjbyK7v7phXvG6492HQfiDaGHLHLQjk7R
/dns4/bootstrap-4.mainnet.filops.net/tcp/1347/p2p/12D3KooWL6PsFNPhYftrJzGgF5U18hFoaVhfGk7xwzD8yVrHJ3Uc
/dns4/bootstrap-5.mainnet.filops.net/tcp/1347/p2p/12D3KooWLFynvDQiUpXoHroV1YxKHhPJgysQGH2k3ZGwtWzR4dFH
/dns4/bootstrap-6.mainnet.filops.net/tcp/1347/p2p/12D3KooWP5MwCiqdMETF9ub1P3MbCvQCcfconnYHbWg6sUJcDRQQ
/dns4/bootstrap-7.mainnet.filops.net/tcp/1347/p2p/12D3KooWRs3aY1p3juFjPy8gPN95PEQChm2QKGUCAdcDCC4EBMKf
/dns4/bootstrap-8.mainnet.filops.net/tcp/1347/p2p/12D3KooWScFR7385LTyR4zU1bYdzSiiAb5rnNABfVahPvVSzyTkR
/dns4/lotus-bootstrap.ipfsforce.com/tcp/41778/p2p/12D3KooWGhufNmZHF3sv48aQeS13ng5XVJZ9E6qy2Ms4VzqeUsHk
/dns4/bootstrap-0.starpool.in/tcp/12757/p2p/12D3KooWDqaZkm3oSczUm3dvAJ5aL2rdSeQ5VQbnHRTQNEFShhmc
/dns4/bootstrap-1.starpool.in/tcp/12757/p2p/12D3KooWSkxqRYoFwtoHJ8cVcoeSpAkfrr4f3wzBUGxhNLYr8Dyb
/dns4/node.glif.io/tcp/1235/p2p/12D3KooWBF8cpp65hp2u9LK5mh19x67ftAam84z9LsfaquTDSBpt
`
	calibTrustNode = `
/dns4/bootstrap-0.calibration.fildev.network/tcp/1347/p2p/12D3KooWJkikQQkxS58spo76BYzFt4fotaT5NpV2zngvrqm4u5ow
/dns4/bootstrap-1.calibration.fildev.network/tcp/1347/p2p/12D3KooWLce5FDHR4EX4CrYavphA5xS3uDsX6aoowXh5tzDUxJav
/dns4/bootstrap-2.calibration.fildev.network/tcp/1347/p2p/12D3KooWA9hFfQG9GjP6bHeuQQbMD3FDtZLdW1NayxKXUT26PQZu
/dns4/bootstrap-3.calibration.fildev.network/tcp/1347/p2p/12D3KooWMHDi3LVTFG8Szqogt7RkNXvonbQYqSazxBx41A5aeuVz
`
)

func GetConnectTrustNode(ctx context.Context, kind string) ([]peer.AddrInfo, error) {
	spi := ""
	switch kind {
	case TOPIC_MAINNET:
		spi = mainTrustNode
	case TOPIC_CALIBNET:
		spi = calibTrustNode
	default:
		return nil, errors.New("Unknow net kind").As(kind)
	}
	pis, err := addrutil.ParseAddresses(ctx, strings.Split(strings.TrimSpace(spi), "\n"))
	if err != nil {
		return nil, errors.As(err)
	}
	return pis, nil
}

func ConnectTrustNode(ctx context.Context, src host.Host, pis []peer.AddrInfo) ([]string, error) {
	result := make([]string, len(pis))
	done := make(chan string, len(pis))
	for i, p := range pis {
		go func(index int, pi peer.AddrInfo) {
			resp := ""
			defer func() {
				result[index] = resp
				done <- resp
			}()
			resp += fmt.Sprintf("connect %s: ", pi.ID.Pretty())
			err := src.Connect(ctx, pi)
			if err != nil {
				resp += fmt.Sprintf("failure:%s", err.Error())
				return
			}
			resp += "success"
		}(i, p)
	}
	noPeers := true
	for _, _ = range pis {
		resp := <-done
		if strings.HasSuffix(resp, "success") {
			noPeers = false
		}
		log.Debug(resp)
	}
	if noPeers {
		return result, errors.New("no available peer")
	}
	return result, nil
}
