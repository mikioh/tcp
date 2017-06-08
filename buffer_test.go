// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp_test

import (
	"net"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/mikioh/tcp"
	"github.com/mikioh/tcpopt"
)

func TestBuffered(t *testing.T) {
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
		if err := c.(*net.TCPConn).SetReadBuffer(1<<16 - 1); err != nil {
			t.Error(err)
			return
		}
		tc, err := tcp.NewConn(c)
		if err != nil {
			t.Error(err)
			return
		}
		time.Sleep(100 * time.Millisecond)
		n := tc.Buffered()
		if n != len(m) {
			switch runtime.GOOS {
			case "darwin", "dragonfly", "freebsd", "linux", "netbsd", "openbsd":
				t.Errorf("got %d; want %d", n, len(m))
			default:
				t.Logf("not supported on %s/%s", runtime.GOOS, runtime.GOARCH)

			}
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

func TestAvailable(t *testing.T) {
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
		if err := c.(*net.TCPConn).SetWriteBuffer(1<<16 - 1); err != nil {
			t.Error(err)
			return
		}
		tc, err := tcp.NewConn(c)
		if err != nil {
			t.Error(err)
			return
		}
		if _, err := c.Write(m); err != nil {
			t.Error(err)
			return
		}
		n := tc.Available()
		if n <= 0 {
			switch runtime.GOOS {
			case "darwin", "freebsd", "linux", "netbsd":
				t.Errorf("got %d; want >0", n)
			default:
				t.Logf("not supported on %s/%s", runtime.GOOS, runtime.GOARCH)
			}
			return
		}
		t.Logf("%d bytes write available", n)
	}()

	c, err := ln.Accept()
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	wg.Wait()
}

func TestCorkAndUncork(t *testing.T) {
	switch runtime.GOOS {
	case "darwin", "freebsd", "linux", "openbsd":
	case "dragonfly":
		t.Log("you may need to adjust the net.inet.tcp.disable_nopush kernel state")
	default:
		t.Skipf("not supported on %s/%s", runtime.GOOS, runtime.GOARCH)
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
			t.Errorf("got %d; want %d", n, N)
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
	if err := tc.SetOption(tcpopt.Cork(true)); err != nil {
		t.Fatal(err)
	}
	b := make([]byte, N)
	for i := 0; i+M <= N; i += M {
		if _, err := tc.Write(b[i : i+M]); err != nil {
			t.Fatal(err)
		}
		time.Sleep(5 * time.Millisecond)
	}
	if err := tc.SetOption(tcpopt.Cork(false)); err != nil {
		t.Fatal(err)
	}

	wg.Wait()
}

func TestBufferCap(t *testing.T) {
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

	for _, o := range []tcpopt.Option{
		tcpopt.SendBuffer(1<<16 - 1),
		tcpopt.ReceiveBuffer(1<<16 - 1),
	} {
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
			if _, err := tc.Option(o.Level(), o.Name(), b[:]); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestBufferTrim(t *testing.T) {
	switch runtime.GOOS {
	case "darwin", "linux":
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

	o := tcpopt.NotSentLowWMK(1)
	var b [4]byte
	if _, err := tc.Option(o.Level(), o.Name(), b[:]); err != nil {
		t.Fatal(err)
	}
	if err := tc.SetOption(o); err != nil {
		t.Fatal(err)
	}
	oo, err := tc.Option(o.Level(), o.Name(), b[:])
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(oo, o) {
		t.Fatalf("got %#v; want %#v", oo, o)
	}
}
