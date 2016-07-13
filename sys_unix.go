// Copyright 2014 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd

package tcp

import (
	"runtime"

	"github.com/mikioh/tcpopt"
)

func buffered(s uintptr) int {
	var b [4]byte
	if err := ioctl(s, options[soBuffered].name, b[:]); err != nil {
		return -1
	}
	return int(nativeEndian.Uint32(b[:]))
}

func available(s uintptr) int {
	var b [4]byte
	if err := ioctl(s, options[soAvailable].name, b[:]); err != nil {
		return -1
	}
	n := int(nativeEndian.Uint32(b[:]))
	if runtime.GOOS == "linux" {
		var o tcpopt.SendBuffer
		if err := getsockopt(s, o.Level(), o.Name(), b[:]); err != nil {
			return -1
		}
		return int(nativeEndian.Uint32(b[:])) - n
	}
	return n
}
