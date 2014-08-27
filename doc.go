// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tcp implements TCP-level socket options.
//
// The package provides TCP-level socket options that allow
// manipulation of TCP connection facilities.
//
//
// Monitoring a TCP connection
//
// For now only Linux and FreeBSD kernels support the TCP information
// option. A custom net.Dial function that hooks up an underlying
// transport connection must be prepared before monitoring.
//
//	func dialWithTCPConnMonitor(network, address string) (net.Conn, error) {
//		d := net.Dialer{}
//		c, err := d.Dial(network, address)
//		if err != nil {
//			return nil, err
//		}
//		tc, err := tcp.NewConn(c)
//		if err != nil {
//			c.Close()
//			return nil, err
//		}
//		go tcpConnMonitor(tc) // launch a TCP connection monitor goroutine
//		return &tc.TCPConn, nil
//	}
//
// Also an application needs to construct a custom client such as a
// HTTP client containing a custom net.Dial function and get into a
// HTTP over TLS over TCP conversation.
//
//	tr := &http.Transport{
//		Dial:            dialWithTCPConnMonitor,
//		TLSClientConfig: &tls.Config{ServerName: host},
//	}
//	client := http.Client{Transport: tr}
//	resp, err := client.Get(url)
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
