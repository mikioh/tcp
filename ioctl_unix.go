// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd

package tcp

import (
	"runtime"
	"syscall"
)

func (c *Conn) readBufferLen() int {
	fd, err := c.sysfd()
	if err != nil {
		return -1
	}
	n, err := getIntByIoctl(fd, &sockOpts[ssoReadBufferLen])
	if err != nil {
		return -1
	}
	return n
}

func (c *Conn) writeBufferSpace() int {
	fd, err := c.sysfd()
	if err != nil {
		return -1
	}
	n, err := getIntByIoctl(fd, &sockOpts[ssoWriteBufferSpace])
	if err != nil {
		return -1
	}
	if runtime.GOOS == "linux" {
		l, err := syscall.GetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_SNDBUF)
		if err != nil {
			return -1
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
