// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp_test

import (
	"net"
	"testing"
	"time"

	"github.com/mikioh/tcp"
)

func TestInfo(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	go func() {
		b := make([]byte, 256)
		for {
			c, err := ln.Accept()
			if err != nil {
				break
			}
			n, err := c.Read(b)
			if err != nil {
				break
			}
			c.Write(b[:n])
			go func() {
				time.Sleep(100 * time.Millisecond)
				c.Close()
			}()
		}
	}()

	c, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	tc, err := tcp.NewConn(c)
	if err != nil {
		t.Fatal(err)
	}
	defer tc.Close()
	if _, err := tc.Write([]byte("HELLO-R-U-THERE")); err != nil {
		t.Fatal(err)
	}
	b := make([]byte, 256)
	tc.Read(b)

	ti, err := tc.Info()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", ti)
	for i, opts := range [][]tcp.Option{ti.Options, ti.PeerOptions} {
		for _, opt := range opts {
			t.Logf("%v/%v: %v", i, opt.Kind(), opt)
		}
	}
}
