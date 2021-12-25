package rpc

import "github.com/gwaylib/errors"

var (
	ErrInvalidToken = errors.New("Invalid token")
	ErrEOF          = errors.New("Unexpected EOF")
)
