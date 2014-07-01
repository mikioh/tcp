// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build freebsd dragonfly

package tcp

const (
	sysGetsockopt = 118

	sysSockoptTCPKeepAliveCount = 0x400
	sysSockoptTCPNoPush         = 0x4
)
