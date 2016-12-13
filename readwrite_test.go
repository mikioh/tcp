// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin freebsd linux netbsd

package tcp_test

import (
	"net"
	"runtime"
	"sync"
	"testing"

	"github.com/mikioh/tcp"
	"github.com/mikioh/tcpinfo"
)

func server(t *testing.T, ln net.Listener) {
	c, err := ln.Accept()
	if err != nil {
		return
	}
	go server(t, ln)
	defer c.Close()
	var b [1]byte
	if _, err := c.Read(b[:]); err != nil {
		t.Error(err)
		return
	}
	if _, err := c.Write(b[:]); err != nil {
		t.Error(err)
		return
	}
}

func TestConcurrentReadWriteAndInfo(t *testing.T) {
	switch runtime.GOOS {
	case "darwin":
		vers := kernelVersion()
		if vers == nil || vers[0] < 15 {
			t.Skipf("not supported on %v, %s/%s", vers, runtime.GOOS, runtime.GOARCH)
		}
	case "freebsd", "linux", "netbsd":
	default:
		t.Skipf("not supported on %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	go server(t, ln)

	const N = 10
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			var d net.Dialer
			c, err := d.Dial(ln.Addr().Network(), ln.Addr().String())
			if err != nil {
				t.Error(err)
				return
			}
			defer c.Close()
			tc, err := tcp.NewConn(c)
			if err != nil {
				t.Error(err)
				return
			}
			var wwg sync.WaitGroup
			wwg.Add(1)
			go func() {
				defer wwg.Done()
				var b [1]byte
				if _, err := c.Write(b[:]); err != nil {
					t.Error(err)
					return
				}
				if _, err := c.Read(b[:]); err != nil {
					t.Error(err)
					return
				}
			}()
			wwg.Add(1)
			go func() {
				defer wwg.Done()
				var o tcpinfo.Info
				b := make([]byte, 256)
				if _, err := tc.Option(o.Level(), o.Name(), b); err != nil {
					t.Error(err)
					return
				}
			}()
			wwg.Wait()
		}()
	}
	wg.Wait()
}
