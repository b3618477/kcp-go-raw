package kcpraw

import (
	"fmt"
	"net"

	"github.com/ccsexyz/rawcon"
	"github.com/pkg/errors"
	kcp "github.com/xtaci/kcp-go"
)

var raw rawcon.Raw

func checkAddr(addr string) (err error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return
	} else if len(host) == 0 {
		err = fmt.Errorf("You must set the addr to ip:port")
	} else if host == "0.0.0.0" {
		err = fmt.Errorf("You can't set host to 0.0.0.0")
	}
	return
}

// DialWithOptions connects to the remote address "raddr" on the network "udp" with packet encryption
func DialWithOptions(raddr string, block kcp.BlockCrypt, dataShards, parityShards int) (*kcp.UDPSession, error) {
	err := checkAddr(raddr)
	if err != nil {
		return nil, errors.Wrap(err, "checkAddr")
	}
	conn, err := raw.DialRAW(raddr)
	if err != nil {
		return nil, errors.Wrap(err, "net.DialRAW")
	}
	return kcp.NewConn(raddr, block, dataShards, parityShards, conn)
}

// ListenWithOptions listens for incoming KCP packets addressed to the local address laddr on the network "udp" with packet encryption,
// dataShards, parityShards defines Reed-Solomon Erasure Coding parameters
func ListenWithOptions(laddr string, block kcp.BlockCrypt, dataShards, parityShards int) (*kcp.Listener, error) {
	err := checkAddr(laddr)
	if err != nil {
		return nil, errors.Wrap(err, "checkAddr")
	}
	conn, err := raw.ListenRAW(laddr)
	if err != nil {
		return nil, errors.Wrap(err, "net.ListenRAW")
	}
	return kcp.ServeConn(block, dataShards, parityShards, conn)
}

// SetNoHTTP determines whether to do http obfuscating
func SetNoHTTP(v bool) {
	raw.NoHTTP = v
}

// SetHost set http host
func SetHost(v string) {
	raw.Host = v
}

// SetDSCP set tos number
func SetDSCP(v int) {
	raw.DSCP = v
}

// SetIgnRST if v is true, the tcp rst packet will be ignored
func SetIgnRST(v bool) {
	raw.IgnRST = v
}
