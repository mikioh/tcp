// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"syscall"
	"unsafe"
)

const sysGETSOCKOPT = 0xf

func socketcall(call int, a0, a1, a2, a3, a4, a5 uintptr) (int, syscall.Errno)

func getsockopt(s uintptr, level, name int, v unsafe.Pointer, l *uint32) error {
	if _, errno := socketcall(sysGETSOCKOPT, s, uintptr(level), uintptr(name), uintptr(v), uintptr(unsafe.Pointer(l)), 0); errno != 0 {
		return error(errno)
	}
	return nil
}

func ioctl(s uintptr, ioc int) (int, error) {
	var i int
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, s, uintptr(ioc), uintptr(unsafe.Pointer(&i))); errno != 0 {
		return 0, error(errno)
	}
	return i, nil
}
