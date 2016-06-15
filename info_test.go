// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin freebsd linux netbsd

package tcp_test

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"runtime"
	"testing"
	"time"

	"github.com/mikioh/tcp"
)

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
			t.Skipf("%v, %s/%s", vers, runtime.GOOS, runtime.GOARCH)
		}
	case "freebsd", "linux", "netbsd":
	default:
		t.Skipf("%s/%s", runtime.GOOS, runtime.GOARCH)
	}

	var logs []chan string
	for i := 0; i < len(infoTests); i++ {
		logs = append(logs, make(chan string, 100))
	}

	for i, tt := range infoTests {
		sig := make(chan struct{})
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
				go monitor(tc, logs[i], sig)
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
		sig <- struct{}{}
	}

	for _, log := range logs {
		for r := range log {
			t.Log(r)
		}
	}
}

func monitor(c *tcp.Conn, log chan<- string, sig <-chan struct{}) {
	defer close(log)
	log <- fmt.Sprintf("%s %v->%v", c.LocalAddr().Network(), c.LocalAddr(), c.RemoteAddr())
	for {
		ti, err := c.Info()
		if err != nil {
			return
		}
		b, err := json.MarshalIndent(ti, "", "\t")
		if err != nil {
			continue
		}
		select {
		case <-sig:
			return
		default:
			log <- string(b)
			time.Sleep(100 * time.Millisecond)
		}
	}
}
