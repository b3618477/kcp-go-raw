package kcpraw

import (
	kcp "github.com/xtaci/kcp-go"
	"github.com/pkg/errors"
)

// DialWithOptions connects to the remote address "raddr" on the network "udp" with packet encryption
func DialWithOptions(raddr string, block kcp.BlockCrypt, dataShards, parityShards int) (*kcp.UDPSession, error) {
	conn, err := dialRAW(raddr)
	if err != nil {
		return nil, errors.Wrap(err, "net.DialRAW")
	}
	return kcp.NewConn(raddr, block, dataShards, parityShards, conn)
}

// ListenWithOptions listens for incoming KCP packets addressed to the local address laddr on the network "udp" with packet encryption,
// dataShards, parityShards defines Reed-Solomon Erasure Coding parameters
func ListenWithOptions(laddr string, block kcp.BlockCrypt, dataShards, parityShards int) (*kcp.Listener, error) {
	conn, err := listenRAW(laddr)
	if err != nil {
		return nil, errors.Wrap(err, "net.ListenUDP")
	}
	return kcp.ServeConn(block, dataShards, parityShards, conn)
}
