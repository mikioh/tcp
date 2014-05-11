// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.2

package tcp

import (
	"syscall"
	"time"
)

// SetKeepAlivePeriod implemets the SetKeepAlivePeriod method of
// net.TCPConn.
func (c *Conn) SetKeepAlivePeriod(d time.Duration) error {
	if !c.opt.ok() {
		return syscall.EINVAL
	}
	return c.TCPConn.SetKeepAlivePeriod(d)
}
