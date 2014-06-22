// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"io"
	"net"
	"syscall"
	"time"
)

// References:
//
// RFC  793  TRANSMISSION CONTROL PROTOCOL
//	http://tools.ietf.org/html/rfc793
// RFC 2780  IANA Allocation Guidelines For Values In the Internet Protocol and Related Headers
//	http://tools.ietf.org/html/rfc2780
// RFC 4022  Management Information Base for the Transmission Control Protocol (TCP)
//	https://tools.ietf.org/html/rfc4022
// RFC 5681 TCP Congestion Control
//	http://tools.ietf.org/html/rfc5681
// RFC 6298 Computing TCP's Retransmission Timer
//	http://tools.ietf.org/html/rfc6298
// RFC 6994  Shared Use of Experimental TCP Options
//	http://tools.ietf.org/html/rfc6994
// 1323bis   TCP Extensions for High Performance
//	http://tools.ietf.org/html/draft-ietf-tcpm-1323bis-21

var _ net.Conn = &Conn{}

// A Conn represents a network endpoint that uses TCP connection.
// It allows to set non-portable, platform-dependent TCP-level socket
// options.
type Conn struct {
	opt
}

type opt struct {
	*net.TCPConn
}

func (c *opt) ok() bool { return c != nil && c.TCPConn != nil }

// Read implements the Read method of net.Conn interface.
func (c *Conn) Read(b []byte) (int, error) {
	if !c.opt.ok() {
		return 0, syscall.EINVAL
	}
	return c.TCPConn.Read(b)
}

// Write implements the Write method of net.Conn interface.
func (c *Conn) Write(b []byte) (int, error) {
	if !c.opt.ok() {
		return 0, syscall.EINVAL
	}
	return c.TCPConn.Write(b)
}

// LocalAddr implements the LocalAddr method of net.Conn interface.
func (c *Conn) LocalAddr() net.Addr {
	if !c.opt.ok() {
		return nil
	}
	return c.TCPConn.LocalAddr()
}

// RemoteAddr implements the RemoteAddr method of net.Conn interface.
func (c *Conn) RemoteAddr() net.Addr {
	if !c.opt.ok() {
		return nil
	}
	return c.TCPConn.RemoteAddr()
}

// SetDeadline implements the SetDeadline method of net.Conn
// interface.
func (c *Conn) SetDeadline(t time.Time) error {
	if !c.opt.ok() {
		return syscall.EINVAL
	}
	return c.TCPConn.SetDeadline(t)
}

// SetReadDeadline implements the SetReadDeadline method of net.Conn
// interface.
func (c *Conn) SetReadDeadline(t time.Time) error {
	if !c.opt.ok() {
		return syscall.EINVAL
	}
	return c.TCPConn.SetReadDeadline(t)
}

// SetWriteDeadline implements the SetWriteDeadline method of net.Conn
// interface.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	if !c.opt.ok() {
		return syscall.EINVAL
	}
	return c.TCPConn.SetWriteDeadline(t)
}

// Close implements the Close method of net.Conn interface.
func (c *Conn) Close() error {
	if !c.opt.ok() {
		return syscall.EINVAL
	}
	return c.TCPConn.Close()
}

// ReadFrom implements the ReadFrom method of io.ReaderFrom interface.
func (c *Conn) ReadFrom(r io.Reader) (int64, error) {
	if !c.opt.ok() {
		return 0, syscall.EINVAL
	}
	return c.TCPConn.ReadFrom(r)
}

// CloseRead implemets the CloseRead method of net.TCPConn.
func (c *Conn) CloseRead() error {
	if !c.opt.ok() {
		return syscall.EINVAL
	}
	return c.TCPConn.CloseRead()
}

// CloseWrite implemets the CloseWrite method of net.TCPConn.
func (c *Conn) CloseWrite() error {
	if !c.opt.ok() {
		return syscall.EINVAL
	}
	return c.TCPConn.CloseWrite()
}

// SetLinger implemets the SetLinger method of net.TCPConn.
func (c *Conn) SetLinger(sec int) error {
	if !c.opt.ok() {
		return syscall.EINVAL
	}
	return c.TCPConn.SetLinger(sec)
}

// SetKeepAlive implements the SetKeepAlive method of net.TCPConn.
func (c *Conn) SetKeepAlive(on bool) error {
	if !c.opt.ok() {
		return syscall.EINVAL
	}
	return c.TCPConn.SetKeepAlive(on)
}

// SetKeepAlivePeriod implemets the SetKeepAlivePeriod method of
// net.TCPConn.
func (c *Conn) SetKeepAlivePeriod(d time.Duration) error {
	if !c.opt.ok() {
		return syscall.EINVAL
	}
	return c.TCPConn.SetKeepAlivePeriod(d)
}

// SetMaxKeepAliveProbes sets the maximum number of keep alive probes.
func (c *Conn) SetMaxKeepAliveProbes(probes int) error {
	if !c.opt.ok() {
		return syscall.EINVAL
	}
	if probes < 1 {
		return syscall.EINVAL
	}
	return c.opt.setMaxKeepAliveProbes(probes)
}

// SetReadBuffer implements the SetReadBuffer method of net.TCPConn.
func (c *Conn) SetReadBuffer(bytes int) error {
	if !c.opt.ok() {
		return syscall.EINVAL
	}
	return c.TCPConn.SetReadBuffer(bytes)
}

// SetWriteBuffer implements the SetWriteBuffer method of net.TCPConn.
func (c *Conn) SetWriteBuffer(bytes int) error {
	if !c.opt.ok() {
		return syscall.EINVAL
	}
	return c.TCPConn.SetWriteBuffer(bytes)
}

// SetNoDelay implements the SetNoDealy method of net.TCPConn.
func (c *Conn) SetNoDelay(on bool) error {
	if !c.opt.ok() {
		return syscall.EINVAL
	}
	return c.TCPConn.SetNoDelay(on)
}

// Cork enables TCP_CORK option on Linux, TCP_NOPUSH option on Darwin,
// DragonFlyBSD, FreeBSD and OpenBSD.
func (c *Conn) Cork() error {
	if !c.opt.ok() {
		return syscall.EINVAL
	}
	if err := c.opt.setCork(true); err != nil {
		return err
	}
	return nil
}

// Uncork disables TCP_CORK option on Linux, TCP_NOPUSH option on
// Darwin, DragonFly BSD, FreeBSD and OpenBSD.
func (c *Conn) Uncork() error {
	if !c.opt.ok() {
		return syscall.EINVAL
	}
	c.opt.setCork(false)
	return nil
}

// Info returns information of current connection. For now this option
// is supported on FreeBSD and Linux.
func (c *Conn) Info() (*Info, error) {
	if !c.opt.ok() {
		return nil, syscall.EINVAL
	}
	return c.opt.info()
}

// NewConn returns a new Conn.
func NewConn(c net.Conn) (*Conn, error) {
	switch c := c.(type) {
	case *net.TCPConn:
		return &Conn{opt: opt{TCPConn: c}}, nil
	default:
		return nil, syscall.EINVAL
	}
}
