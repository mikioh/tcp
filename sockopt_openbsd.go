// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"os"
	"syscall"
	"time"
)

func (c *Conn) setKeepAliveIdlePeriod(d time.Duration) error {
	return errOpNoSupport
}

func (c *Conn) setKeepAliveProbeInterval(d time.Duration) error {
	return errOpNoSupport
}

func (c *Conn) setMaxKeepAliveProbes(max int) error {
	return errOpNoSupport
}

func (c *Conn) setCork(on bool) error {
	fd, err := c.sysfd()
	if err != nil {
		return err
	}
	return os.NewSyscallError("setsockopt", syscall.SetsockoptInt(fd, ianaProtocolTCP, sysTCP_NOPUSH, boolint(on)))
}

func (c *Conn) info() (*Info, error) {
	return nil, errOpNoSupport
}
