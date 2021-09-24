package main

import (
	"crypto/rand"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
)

func genLibp2pKey() (crypto.PrivKey, error) {
	pk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, err
	}
	return pk, nil
}

func NetID() libp2p.Option {
	// TODO: read from keystore
	pk, err := genLibp2pKey()
	if err != nil {
		panic(err)
	}
	return libp2p.Identity(pk)
}
