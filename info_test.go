// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build freebsd linux,amd64 linux,arm

package tcp_test

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net"
	"net/http"
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
	return tc, nil
}

func tcpConnMonitor(c *tcp.Conn) {
	dir := func(i int) string {
		if i == 0 {
			return "requesting"
		}
		return "requested"
	}
	tt.Logf("%v -> %v", c.LocalAddr(), c.RemoteAddr())
	for {
		ti, err := c.Info()
		if err != nil {
			break
		}
		tt.Logf("---- State: %v ----", ti.State)
		for i, opts := range [][]tcp.Option{ti.Options, ti.PeerOptions} {
			for _, opt := range opts {
				tt.Logf("%v %v: %v", dir(i), opt.Kind(), opt)
			}
		}
		tt.Logf("MSS: sender: %v, receiver: %v", ti.SenderMSS, ti.ReceiverMSS)
		tt.Logf("Time: last data sent: %v, last data received: %v, last ack received: %v", ti.LastDataSent, ti.LastDataReceived, ti.LastAckReceived)
		if ti.CC != nil {
			tt.Logf("CC: rto: %v, ato: %v, rtt: %v, rtt stdev: %v, sender ssthresh: %v, receiver ssthresh: %v, window: %v", ti.CC.RTO, ti.CC.ATO, ti.CC.RTT, ti.CC.RTTStdDev, ti.CC.SenderSSThreshold, ti.CC.ReceiverSSThreshold, ti.CC.SenderWindow)
		}
		tt.Logf("SysInfo: %+v", ti.SysInfo)
		time.Sleep(30 * time.Millisecond)
	}
}
