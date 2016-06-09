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

func setCork(s uintptr, on bool) error {
	return setInt(s, &sockOpts[ssoCork], boolint(on))
}

func setKeepAliveIdleInterval(s uintptr, d time.Duration) error {
	d += (sockOpts[ssoKeepAliveIdleInterval].tmu - time.Nanosecond)
	v := int(d / sockOpts[ssoKeepAliveIdleInterval].tmu)
	return setInt(s, &sockOpts[ssoKeepAliveIdleInterval], v)
}

func setKeepAliveProbeInterval(s uintptr, d time.Duration) error {
	d += (sockOpts[ssoKeepAliveIdleInterval].tmu - time.Nanosecond)
	v := int(d / sockOpts[ssoKeepAliveIdleInterval].tmu)
	return setInt(s, &sockOpts[ssoKeepAliveProbeInterval], v)
}

func setKeepAliveProbeCount(s uintptr, probes int) error {
	return setInt(s, &sockOpts[ssoKeepAliveProbeCount], probes)
}

func getInt(s uintptr, opt *sockOpt) (int, error) {
	if opt.name < 1 || opt.typ != ssoTypeInt {
		return 0, errOpNoSupport
	}
	v, err := syscall.GetsockoptInt(int(s), ianaProtocolTCP, opt.name)
	if err != nil {
		return 0, os.NewSyscallError("getsockopt", err)
	}
	return v, nil
}

func setInt(s uintptr, opt *sockOpt, v int) error {
	if opt.name < 1 || opt.typ != ssoTypeInt {
		return errOpNoSupport
	}
	return os.NewSyscallError("setsockopt", syscall.SetsockoptInt(int(s), ianaProtocolTCP, opt.name, v))
}
