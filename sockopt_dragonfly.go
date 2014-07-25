// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.3

package tcp

import (
	"os"
	"syscall"
)

func (opt *opt) setMaxKeepAliveProbes(max int) error {
	fd, err := opt.sysfd()
	if err != nil {
		return err
	}
	return os.NewSyscallError("setsockopt", syscall.SetsockoptInt(fd, ianaProtocolTCP, sysTCP_KEEPCNT, max))
}

func (opt *opt) setCork(on bool) error {
	fd, err := opt.sysfd()
	if err != nil {
		return err
	}
	return os.NewSyscallError("setsockopt", syscall.SetsockoptInt(fd, ianaProtocolTCP, sysTCP_NOPUSH, boolint(on)))
}

func (opt *opt) info() (*Info, error) {
	return nil, errOpNoSupport
}
