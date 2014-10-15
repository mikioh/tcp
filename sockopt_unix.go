// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package tcp

import (
	"os"
	"syscall"
	"time"
)

func (c *Conn) setKeepAliveIdleInterval(d time.Duration) error {
	fd, err := c.sysfd()
	if err != nil {
		return err
	}
	d += (sockOpts[ssoKeepAliveIdleInterval].tmu - time.Nanosecond)
	v := int(d / sockOpts[ssoKeepAliveIdleInterval].tmu)
	if err := setInt(fd, &sockOpts[ssoKeepAliveIdleInterval], v); err != nil {
		return err
	}
	return nil
}

func (c *Conn) setKeepAliveProbeInterval(d time.Duration) error {
	fd, err := c.sysfd()
	if err != nil {
		return err
	}
	d += (sockOpts[ssoKeepAliveIdleInterval].tmu - time.Nanosecond)
	v := int(d / sockOpts[ssoKeepAliveIdleInterval].tmu)
	if err := setInt(fd, &sockOpts[ssoKeepAliveProbeInterval], v); err != nil {
		return err
	}
	return nil
}

func (c *Conn) setKeepAliveProbes(probes int) error {
	if probes < 1 {
		return errInvalidArgument
	}
	fd, err := c.sysfd()
	if err != nil {
		return err
	}
	if err := setInt(fd, &sockOpts[ssoKeepAliveProbes], probes); err != nil {
		return err
	}
	return nil
}

func (c *Conn) setCork(on bool) error {
	fd, err := c.sysfd()
	if err != nil {
		return err
	}
	if err := setInt(fd, &sockOpts[ssoCork], boolint(on)); err != nil {
		return err
	}
	return nil
}

func getInt(fd int, opt *sockOpt) (int, error) {
	if opt.name < 1 || opt.typ != ssoTypeInt {
		return 0, errOpNoSupport
	}
	v, err := syscall.GetsockoptInt(fd, ianaProtocolTCP, opt.name)
	if err != nil {
		return 0, os.NewSyscallError("getsockopt", err)
	}
	return v, nil
}

func setInt(fd int, opt *sockOpt, v int) error {
	if opt.name < 1 || opt.typ != ssoTypeInt {
		return errOpNoSupport
	}
	return os.NewSyscallError("setsockopt", syscall.SetsockoptInt(fd, ianaProtocolTCP, opt.name, v))
}
