// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp_test

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"runtime"
	"testing"
	"time"

	"github.com/mikioh/tcp"
)

var (
	host = "www.google.com"
	url  = "https://www.google.com/robots.txt"
	tt   *testing.T
)

func TestInfoWithGoogle(t *testing.T) {
	switch runtime.GOOS {
	case "freebsd", "linux":
	default:
		t.Skipf("%s/%s", runtime.GOOS, runtime.GOARCH)
	}

	tt = t
	tr := &http.Transport{
		Dial:            dialWithTCPConnMonitor,
		TLSClientConfig: &tls.Config{ServerName: host},
	}
	client := http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
}

func dialWithTCPConnMonitor(network, address string) (net.Conn, error) {
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
	go tcpConnMonitor(tc)
	return &tc.TCPConn, nil
}

func tcpConnMonitor(c *tcp.Conn) {
	tt.Logf("%v -> %v", c.LocalAddr(), c.RemoteAddr())
	for {
		ti, err := c.Info()
		if err != nil {
			tt.Error(err)
			return
		}
		text, err := json.Marshal(ti)
		if err != nil {
			tt.Error(err)
			return
		}
		tt.Log(string(text))
		time.Sleep(5 * time.Millisecond)
	}
}
