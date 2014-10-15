// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd

package tcp

import (
	"runtime"
	"syscall"
)

func (c *Conn) buffered() int {
	fd, err := c.sysfd()
	if err != nil {
		return 0
	}
	n, err := getIntByIoctl(fd, &sockOpts[ssoBuffered])
	if err != nil {
		return 0
	}
	return n
}

func (c *Conn) available() int {
	fd, err := c.sysfd()
	if err != nil {
		return 0
	}
	n, err := getIntByIoctl(fd, &sockOpts[ssoAvailable])
	if err != nil {
		return 0
	}
	if runtime.GOOS == "linux" {
		l, err := syscall.GetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_RCVBUF)
		if err != nil {
			return 0
		}
		return l - n
	}
	return n
}

func getIntByIoctl(fd int, opt *sockOpt) (int, error) {
	if opt.name < 1 || opt.typ != ssoTypeInt {
		return 0, errOpNoSupport
	}
	v, err := getsockoptIntByIoctl(fd, opt.name)
	if err != nil {
		return 0, err
	}
	return v, nil
}
