// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"errors"
	"net"
	"time"
)

var (
	_ net.Conn = &Conn{}

	errInvalidArgument = errors.New("invalid argument")
)

// A Conn represents a network endpoint that uses TCP connection.
// It allows to set non-portable, platform-dependent TCP-level socket
// options.
type Conn struct {
	net.TCPConn
}

// ReadBufferLen returns the number of bytes that can be read from the
// underlying socket read buffer. It returns -1 when the platform
// doesn't support this feature.
func (c *Conn) ReadBufferLen() int {
	return c.readBufferLen()
}

// WriteBufferSpace returns how many bytes are unused in the
// underlying socket write buffer. It returns -1 when the platform
// doesn't support this feature.
func (c *Conn) WriteBufferSpace() int {
	return c.writeBufferSpace()
}

// SetKeepAliveIdleInterval sets the idle interval until the first
// probe is sent.
func (c *Conn) SetKeepAliveIdleInterval(d time.Duration) error {
	return c.setKeepAliveIdleInterval(d)
}

// SekKeepAliveProbeInterval sets the interval between keepalive
// probes.
func (c *Conn) SetKeepAliveProbeInterval(d time.Duration) error {
	return c.setKeepAliveProbeInterval(d)
}

// SetKeepAliveProbes sets the number of keepalive probes.
func (c *Conn) SetKeepAliveProbes(probes int) error {
	if probes < 1 {
		return errInvalidArgument
	}
	return c.setKeepAliveProbes(probes)
}

// Cork enables TCP_CORK option on Linux, TCP_NOPUSH option on Darwin,
// DragonFlyBSD, FreeBSD and OpenBSD.
func (c *Conn) Cork() error {
	if err := c.setCork(true); err != nil {
		return err
	}
	return nil
}

// Uncork disables TCP_CORK option on Linux, TCP_NOPUSH option on
// Darwin, DragonFly BSD, FreeBSD and OpenBSD.
func (c *Conn) Uncork() error {
	return c.setCork(false)
}

// Info returns information of current connection. For now this option
// is supported on FreeBSD and Linux.
func (c *Conn) Info() (*Info, error) {
	return c.info()
}

// NewConn returns a new Conn.
func NewConn(c net.Conn) (*Conn, error) {
	switch c := c.(type) {
	case *net.TCPConn:
		return &Conn{TCPConn: *c}, nil
	default:
		return nil, errInvalidArgument
	}
}
