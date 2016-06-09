// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd

package tcp

import (
	"runtime"
	"syscall"
)

func buffered(s uintptr) int {
	n, err := ioctlGetInt(s, &sockOpts[ssoBuffered])
	if err != nil {
		return -1
	}
	return n
}

func available(s uintptr) int {
	n, err := ioctlGetInt(s, &sockOpts[ssoAvailable])
	if err != nil {
		return -1
	}
	if runtime.GOOS == "linux" {
		l, err := syscall.GetsockoptInt(int(s), syscall.SOL_SOCKET, syscall.SO_SNDBUF)
		if err != nil {
			return -1
		}
		return l - n
	}
	return n
}

func ioctlGetInt(s uintptr, opt *sockOpt) (int, error) {
	if opt.name < 1 || opt.typ != ssoTypeInt {
		return 0, errOpNoSupport
	}
	v, err := ioctl(s, opt.name)
	if err != nil {
		return 0, err
	}
	return v, nil
}
