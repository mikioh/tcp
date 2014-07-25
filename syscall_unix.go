// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build freebsd linux,amd64 linux,arm

package tcp

import (
	"syscall"
	"unsafe"
)

type sysSockoptLen uint32

func getsockopt(fd int, level, name int, v unsafe.Pointer, l *sysSockoptLen) error {
	if _, _, errno := syscall.Syscall6(sysGETSOCKOPT, uintptr(fd), uintptr(level), uintptr(name), uintptr(v), uintptr(unsafe.Pointer(l)), 0); errno != 0 {
		return error(errno)
	}
	return nil
}
