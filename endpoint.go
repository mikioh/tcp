// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"net"
	"time"

	"github.com/mikioh/netreflect"
)

var _ net.Conn = &Conn{}

// A Conn represents a network endpoint that uses TCP connection.
// It allows to set non-portable, platform-dependent TCP-level socket
// options.
type Conn struct {
	net.Conn
	s int // socket descriptor for avoding data race
}

// A KeepAliveOptions represents keepalive options.
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
	if opt.IdleInterval >= 0 { // BSD variants accept 0, Linux doesn't
		if err := setKeepAliveIdleInterval(c.s, opt.IdleInterval); err != nil {
			return err
		}
	}
	if opt.ProbeInterval >= 0 { // BSD variants accept 0, Linux doesn't
		if err := setKeepAliveProbeInterval(c.s, opt.ProbeInterval); err != nil {
			return err
		}
	}
	if opt.ProbeCount >= 0 { // BSD variants accept 0, Linux doesn't
		if err := setKeepAliveProbeCount(c.s, opt.ProbeCount); err != nil {
			return err
		}
	}
	return nil
}

// A BufferOptions represents buffer options.
type BufferOptions struct {
	// The runtime-integrated network poller doesn't report that
	// the connection is writable while the amount of unsent data
	// size is greater than UnsentThreshold.
	//
	// For now only Darwin and Linux support this option.
	// See TCP_NOTSENT_LOWAT for further information.
	UnsentThreshold int
}

// SetBufferOptions sets buffer options.
func (c *Conn) SetBufferOptions(opt *BufferOptions) error {
	if opt.UnsentThreshold >= 0 {
		if err := setInt(c.s, &sockOpts[ssoUnsentThreshold], opt.UnsentThreshold); err != nil {
			return err
		}
	}
	return nil
}

// Buffered returns the number of bytes that can be read from the
// underlying socket read buffer.
// It returns -1 when the platform doesn't support this feature.
func (c *Conn) Buffered() int {
	return buffered(c.s)
}

// Available returns how many bytes are unused in the underlying
// socket write buffer.
// It returns -1 when the platform doesn't support this feature.
func (c *Conn) Available() int {
	return available(c.s)
}

// Cork enables TCP_CORK option on Linux, TCP_NOPUSH option on Darwin,
// DragonFlyBSD, FreeBSD and OpenBSD.
func (c *Conn) Cork() error {
	return setCork(c.s, true)
}

// Uncork disables TCP_CORK option on Linux, TCP_NOPUSH option on
// Darwin, DragonFly BSD, FreeBSD and OpenBSD.
func (c *Conn) Uncork() error {
	return setCork(c.s, false)
}

// Info returns information of current connection.
// For now this option is supported on Darwin, FreeBSD and Linux.
func (c *Conn) Info() (*Info, error) {
	return info(c.s)
}

// NewConn returns a new Conn.
func NewConn(c net.Conn) (*Conn, error) {
	s, err := netreflect.SocketOf(c)
	if err != nil {
		return nil, err
	}
	tc := &Conn{Conn: c, s: int(s)}
	return tc, nil
}
