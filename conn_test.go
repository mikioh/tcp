// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp_test

import (
	"net"
	"testing"

	"github.com/mikioh/tcp"
	"golang.org/x/net/nettest"
)

func TestConn(t *testing.T) {
	nettest.TestConn(t, func() (net.Conn, net.Conn, func(), error) {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return nil, nil, nil, err
		}
		defer ln.Close()
		type pasv struct {
			net.Conn
			error
		}
		ch := make(chan pasv)
		go func() {
			var p pasv
			p.Conn, p.error = ln.Accept()
			ch <- p
		}()
		c, err := net.Dial(ln.Addr().Network(), ln.Addr().String())
		if err != nil {
			return nil, nil, nil, err
		}
		tc, err := tcp.NewConn(c)
		if err != nil {
			c.Close()
			return nil, nil, nil, err
		}
		p := <-ch
		if p.error != nil {
			tc.Close()
			return nil, nil, nil, err
		}
		return tc, p.Conn, func() { tc.Close(); p.Conn.Close() }, nil
	})
}
