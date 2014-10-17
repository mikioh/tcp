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

func TestReadBuffer(t *testing.T) {
	switch runtime.GOOS {
	case "darwin", "dragonfly", "freebsd", "linux", "netbsd", "openbsd":
	default:
		t.Skipf("%s/%s", runtime.GOOS, runtime.GOARCH)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	m := []byte("HELLO-R-U-THERE")
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
		if err := c.(*net.TCPConn).SetReadBuffer(65535); err != nil {
			t.Error(err)
			return
		}
		tc, err := tcp.NewConn(c)
		if err != nil {
			t.Error(err)
			return
		}
		time.Sleep(100 * time.Millisecond)
		n := tc.ReadBufferLen()
		if n != len(m) {
			t.Errorf("got %v; want %v", n, len(m))
			return
		}
		t.Logf("%v bytes buffered to be read", n)
	}()

	c, err := net.Dial(ln.Addr().Network(), ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	if _, err := c.Write(m); err != nil {
		t.Fatal(err)
	}

	wg.Wait()
}

func TestWriteBuffer(t *testing.T) {
	switch runtime.GOOS {
	case "freebsd", "linux", "netbsd":
	default:
		t.Skipf("%s/%s", runtime.GOOS, runtime.GOARCH)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	m := []byte("HELLO-R-U-THERE")
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		d := net.Dialer{Timeout: 200 * time.Millisecond}
		c, err := d.Dial(ln.Addr().Network(), ln.Addr().String())
		if err != nil {
			t.Error(err)
			return
		}
		defer c.Close()
		if err := c.(*net.TCPConn).SetWriteBuffer(65535); err != nil {
			t.Error(err)
			return
		}
		tc, err := tcp.NewConn(c)
		if err != nil {
			t.Error(err)
			return
		}
		defer tc.Close()
		if _, err := c.Write(m); err != nil {
			t.Error(err)
			return
		}
		n := tc.WriteBufferSpace()
		if n <= 0 {
			t.Errorf("got %v; want >0", n)
			return
		}
		t.Logf("%v bytes write available", n)
	}()

	c, err := ln.Accept()
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	wg.Wait()
}
