// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp_test

import (
	"math/rand"
	"net"
	"os"
	"runtime"
	"testing"

	"github.com/mikioh/tcp"
)

func TestOptionWithVariousBufferLengths(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	c, err := net.Dial(ln.Addr().Network(), ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	tc, err := tcp.NewConn(c)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 256; i++ {
		level, name := rand.Int(), rand.Int()
		b := make([]byte, i)
		tc.Option(level, name, b)
	}
}

func TestOriginalDst(t *testing.T) {
	switch runtime.GOOS {
	case "linux":
	case "darwin", "dragonfly", "freebsd", "openbsd":
		if os.Getuid() != 0 {
			t.Skip("must be root")
		}
	default:
		t.Skipf("not supported on %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	for _, address := range []string{"127.0.0.1:0"} {
		ln, err := net.Listen("tcp", address)
		if err != nil {
			t.Fatal(err)
		}
		defer ln.Close()
		c, err := net.Dial(ln.Addr().Network(), ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}
		defer c.Close()

		tc, err := tcp.NewConn(c)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := tc.OriginalDst(); err != nil {
			t.Log(err)
		}
	}
}
