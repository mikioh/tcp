// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp_test

import (
	"net"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/mikioh/tcp"
)

func TestCorkAndUncork(t *testing.T) {
	switch runtime.GOOS {
	case "darwin", "dragonfly", "freebsd", "linux", "openbsd", "solaris":
	default:
		t.Skipf("%s/%s", runtime.GOOS, runtime.GOARCH)
	}

	const N = 1280
	const M = N / 10

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := ln.(*net.TCPListener).SetDeadline(time.Now().Add(200 * time.Millisecond)); err != nil {
			t.Error(err)
			return
		}
		c, err := ln.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer c.Close()
		b := make([]byte, N)
		n, err := c.Read(b)
		if err != nil {
			t.Error(err)
			return
		}
		if n != N {
			t.Errorf("got %v; want %v", n, N)
			return
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
	if err := tc.Cork(); err != nil {
		t.Fatal(err)
	}
	b := make([]byte, N)
	for i := 0; i+M <= N; i += M {
		if _, err := tc.Write(b[i : i+M]); err != nil {
			t.Fatal(err)
		}
		time.Sleep(5 * time.Millisecond)
	}
	if err := tc.Uncork(); err != nil {
		t.Fatal(err)
	}

	wg.Wait()
}
