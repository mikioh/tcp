// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp_test

import (
	"net"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/mikioh/tcp"
	"github.com/mikioh/tcpopt"
)

func TestKeepAliveOptions(t *testing.T) {
	opts := []tcpopt.Option{
		tcpopt.KeepAlive(true),
		tcpopt.KeepAliveIdleInterval(10 * time.Second), // solaris requires 10 seconds as the lowest value
		tcpopt.KeepAliveProbeInterval(time.Second),
	}
	switch runtime.GOOS {
	case "darwin", "dragonfly", "freebsd", "linux", "netbsd", "solaris":
		opts = append(opts, tcpopt.KeepAliveProbeCount(1))
	case "windows":
	default:
		t.Skipf("not supported on %s/%s", runtime.GOOS, runtime.GOARCH)
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

	for _, o := range opts {
		var b [4]byte
		if runtime.GOOS != "windows" {
			if _, err := tc.Option(o.Level(), o.Name(), b[:]); err != nil {
				t.Fatal(err)
			}
		}
		if err := tc.SetOption(o); err != nil {
			t.Fatal(err)
		}
		if runtime.GOOS != "windows" {
			oo, err := tc.Option(o.Level(), o.Name(), b[:])
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(oo, o) {
				t.Fatalf("got %#v; want %#v", oo, o)
			}
		}
	}
}
