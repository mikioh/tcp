// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris

package tcp

import "time"

func (c *Conn) setKeepAliveIdleInterval(d time.Duration) error {
	return errOpNoSupport
}

func (c *Conn) setKeepAliveProbeInterval(d time.Duration) error {
	return errOpNoSupport
}

func (c *Conn) setKeepAliveProbes(max int) error {
	return errOpNoSupport
}

func (c *Conn) setCork(on bool) error {
	return errOpNoSupport
}

func (c *Conn) info() (*Info, error) {
	return nil, errOpNoSupport
}
