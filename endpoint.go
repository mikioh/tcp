// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tcp implements TCP-level socket options.
//
// The package provides TCP-level socket options that allow
// manipulation of TCP connection facilities.
package tcp

import (
	"io"
	"net"
	"sync"
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
// RFC 6994  Shared Use of Experimental TCP Options
//	http://tools.ietf.org/html/rfc6994

var _ net.Conn = &Conn{}

// A Conn represents a network endpoint that uses TCP connection.
// It allows to set non-portable, platfrom-dependent  TCP-level socket
// options.
type Conn struct {
	opt
}

type opt struct {
	*net.TCPConn
	flowctl struct { // flow control
		sync.RWMutex
		corking bool // corking or no-pushing
	}
	congctl struct { // congestion control
	}
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

// SetBufferedWrite enables TCP_CORK option on Linux, TCP_NOPUSH
// option on FreeBSD.
func (c *Conn) SetBufferedWrite() error {
	if !c.opt.ok() {
		return syscall.EINVAL
	}
	c.opt.flowctl.Lock()
	defer c.opt.flowctl.Unlock()
	if c.opt.flowctl.corking {
		return nil
	}
	if err := c.opt.setCorkedBuffer(true); err != nil {
		return err
	}
	c.opt.flowctl.corking = true
	return nil
}

// Flush disables TCP_CORK option on Linux, TCP_NOPUSH option on
// FreeBSD.
func (c *Conn) Flush() error {
	if !c.opt.ok() {
		return syscall.EINVAL
	}
	c.opt.flowctl.Lock()
	c.opt.setCorkedBuffer(false)
	c.opt.flowctl.corking = false
	c.opt.flowctl.Unlock()
	return nil
}

// Info returns connection on the connection.
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
