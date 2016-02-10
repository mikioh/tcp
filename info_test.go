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

var testingT *testing.T

var infoTests = []struct {
	host, url string
}{
	{
		host: "www.google.com",
		url:  "https://www.google.com/robots.txt",
	},
	{
		host: "github.com",
		url:  "https://github.com/robots.txt",
	},
}

func TestInfo(t *testing.T) {
	switch runtime.GOOS {
	case "darwin", "freebsd", "linux":
	default:
		t.Skipf("%s/%s", runtime.GOOS, runtime.GOARCH)
	}

	testingT = t
	for _, tt := range infoTests {
		tr := &http.Transport{
			Dial:            dialWithTCPConnMonitor,
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
	testingT.Logf("%s %v->%v", c.LocalAddr().Network(), c.LocalAddr(), c.RemoteAddr())
	for {
		ti, err := c.Info()
		if err != nil {
			testingT.Error(err)
			return
		}
		text, err := json.Marshal(ti)
		if err != nil {
			testingT.Error(err)
			return
		}
		testingT.Log(string(text))
		time.Sleep(100 * time.Millisecond)
	}
}
