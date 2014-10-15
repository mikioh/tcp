// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp_test

import (
	"net"
	"runtime"
	"testing"
	"time"

	"github.com/mikioh/tcp"
)

func TestKeepAlive(t *testing.T) {
	switch runtime.GOOS {
	case "darwin", "dragonfly", "freebsd", "linux", "netbsd", "solaris", "windows":
	default:
		t.Skipf("%s/%s", runtime.GOOS, runtime.GOARCH)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				break
			}
			defer c.Close()
		}
	}()

	c, err := net.Dial(ln.Addr().Network(), ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	tc, err := tcp.NewConn(c)
	if err != nil {
		t.Fatal(err)
	}
	defer tc.Close()
	if err := tc.SetKeepAlive(true); err != nil {
		t.Error(err)
	}
	if err := tc.SetKeepAliveIdleInterval(time.Second); err != nil {
		t.Error(err)
	}
	if err := tc.SetKeepAliveProbeInterval(time.Second); err != nil {
		t.Error(err)
	}
	if err := tc.SetKeepAliveProbes(1); err != nil {
		t.Error(err)
	}
}
