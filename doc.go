// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tcp implements TCP-level socket options.
//
// The package provides TCP-level socket options that allow
// manipulation of TCP connection facilities.
//
// The Transmission Control Protocol (TCP) is defined in RFC 793.
// TCP Selective Acknowledgment Options is defined in RFC 2018.
// Management Information Base for the Transmission Control Protocol
// (TCP) is defined in RFC 4022.
// TCP Congestion Control is defined in RFC 5681.
// Computing TCP's Retransmission Timer is described in RFC 6298.
// TCP Options and Maximum Segment Size (MSS) is defined in RFC 6691.
// Shared Use of Experimental TCP Options is defined in RFC 6994.
// TCP Extensions for High Performance is defined in RFC 7323.
//
//
// Monitoring a TCP connection
//
// For now only Darwin, FreeBSD, Linux and NetBSD kernels support the
// TCP information option. A custom net.Dial function that hooks up an
// underlying transport connection must be prepared before monitoring.
//
//	tr := &http.Transport{
//		Dial: func(network, host string) (net.Conn, error) {
//			d := net.Dialer{DualStack: true}
//			c, err := d.Dial(network, host)
//			if err != nil {
//				return nil, err
//			}
//			tc, err := tcp.NewConn(c)
//			if err != nil {
//				c.Close()
//				return nil, err
//			}
//			go tcpConnMonitor(tc)
//			return tc.Conn, nil
//		},
//		TLSClientConfig: &tls.Config{ServerName: "golang.org"},
//	}
//	client := http.Client{Transport: tr}
//	resp, err := client.Get("https://golang.org")
//	if err != nil {
//		// error handling
//	}
//
// When the underlying transport connection is established, your
// monitor goroutine can start monitoring the TCP connection by using
// the Info method of tcp.Conn.
//
//	 func tcpConnMonitor(c *tcp.Conn) {
//		for {
//			ti, err := c.Info() // fetch TCP connection information
//			if err != nil {
//				// error handling
//			}
//			text, err := json.Marshal(ti)
//			if err != nil {
//				// error handling
//			}
//			fmt.Println(text)
package tcp
