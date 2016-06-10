// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin freebsd linux netbsd

package tcp_test

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/mikioh/tcp"
)

func kernelVersion() []int {
	s, err := syscall.Sysctl("kern.osrelease")
	if err != nil {
		return nil
	}
	ss := strings.Split(s, ".")
	vers := make([]int, len(ss))
	for i, s := range ss {
		vers[i], _ = strconv.Atoi(s)
	}
	return vers
}

var infoTests = []struct {
	host, url string
}{
	{
		host: "golang.org",
		url:  "https://golang.org/robots.txt",
	},
	{
		host: "github.com",
		url:  "https://github.com/robots.txt",
	},
}

func TestInfo(t *testing.T) {
	switch runtime.GOOS {
	case "darwin":
		vers := kernelVersion()
		if vers == nil || vers[0] < 15 {
			t.Skip("not supported on this version: %v, %s/%s", vers, runtime.GOOS, runtime.GOARCH)
		}
	case "freebsd", "linux", "netbsd":
	default:
		t.Skipf("%s/%s", runtime.GOOS, runtime.GOARCH)
	}

	for _, tt := range infoTests {
		tr := &http.Transport{
			Dial: func(network, address string) (net.Conn, error) {
				d := net.Dialer{DualStack: true}
				c, err := d.Dial(network, address)
				if err != nil {
					return nil, err
				}
				tc, err := tcp.NewConn(c)
				if err != nil {
					c.Close()
					return nil, err
				}
				go tcpConnMonitor(t, tc)
				return tc.Conn, nil
			},
			TLSClientConfig: &tls.Config{ServerName: tt.host},
		}
		client := http.Client{Transport: tr}
		resp, err := client.Get(tt.url)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		time.Sleep(100 * time.Millisecond)
	}
}

func tcpConnMonitor(t *testing.T, c *tcp.Conn) {
	t.Logf("%s %v->%v", c.LocalAddr().Network(), c.LocalAddr(), c.RemoteAddr())
	for {
		ti, err := c.Info()
		if err != nil {
			return
		}
		b, err := json.MarshalIndent(ti, "", "\t")
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(string(b))
		time.Sleep(100 * time.Millisecond)
	}
}
