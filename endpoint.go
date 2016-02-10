// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"net"
	"time"
)

var _ net.Conn = &Conn{}

// A Conn represents a network endpoint that uses TCP connection.
// It allows to set non-portable, platform-dependent TCP-level socket
// options.
type Conn struct {
	net.TCPConn
}

// KeepAliveOptions represents keepalive options.
type KeepAliveOptions struct {
	// IdleInterval is the idle interval until the first probe is
	// sent.
	// See TCP_KEEPIDLE or TCP_KEEPALIVE for further information.
	IdleInterval time.Duration

	// ProbeInterval is the interval between keepalive probes.
	// See TCP_KEEPINTVL for further information.
	ProbeInterval time.Duration

	// ProbeCount is the number of keepalive probes should be
	// repeated when the peer is not responding.
	// See TCP_KEEPCNT for further information.
	ProbeCount int
}

// SetKeepAliveOptions sets keepalive options.
func (c *Conn) SetKeepAliveOptions(opt *KeepAliveOptions) error {
	s, err := c.sysfd()
	if err != nil {
		return err
	}
	if opt.IdleInterval > 0 {
		if err := setKeepAliveIdleInterval(s, opt.IdleInterval); err != nil {
			return err
		}
	}
	if opt.ProbeInterval > 0 {
		if err := setKeepAliveProbeInterval(s, opt.ProbeInterval); err != nil {
			return err
		}
	}
	if opt.ProbeCount > 0 {
		if err := setKeepAliveProbeCount(s, opt.ProbeCount); err != nil {
			return err
		}
	}
	return nil
}

// BufferOptions represents buffer options.
type BufferOptions struct {
	// The runtime-integrated network poller doesn't report that
	// the connection is writable while the amount of unsent TCP
	// data size is greater than NotsentLowWatermark.
	//
	// For now only Darwin and Linux support this option.
	// See TCP_NOTSENT_LOWAT for further information.
	NotsentLowWatermark int
}

// SetBufferOptions sets buffer options.
func (c *Conn) SetBufferOptions(opt *BufferOptions) error {
	s, err := c.sysfd()
	if err != nil {
		return err
	}
	if opt.NotsentLowWatermark > 0 {
		if err := setInt(s, &sockOpts[ssoNotsentLowWatermark], opt.NotsentLowWatermark); err != nil {
			return err
		}
	}
	return nil
}

// Buffered returns the number of bytes that can be read from the
// underlying socket read buffer.
// It returns -1 when the platform doesn't support this feature.
func (c *Conn) Buffered() int {
	s, err := c.sysfd()
	if err != nil {
		return -1
	}
	return buffered(s)
}

// Available returns how many bytes are unused in the underlying
// socket write buffer.
// It returns -1 when the platform doesn't support this feature.
func (c *Conn) Available() int {
	s, err := c.sysfd()
	if err != nil {
		return -1
	}
	return available(s)
}

// Cork enables TCP_CORK option on Linux, TCP_NOPUSH option on Darwin,
// DragonFlyBSD, FreeBSD and OpenBSD.
func (c *Conn) Cork() error {
	s, err := c.sysfd()
	if err != nil {
		return err
	}
	return setCork(s, true)
}

// Uncork disables TCP_CORK option on Linux, TCP_NOPUSH option on
// Darwin, DragonFly BSD, FreeBSD and OpenBSD.
func (c *Conn) Uncork() error {
	s, err := c.sysfd()
	if err != nil {
		return err
	}
	return setCork(s, false)
}

// Info returns information of current connection.
// For now this option is supported on Darwin, FreeBSD and Linux.
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
